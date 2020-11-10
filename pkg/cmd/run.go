package cmd

import (
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

const (
	msgEnds = "\n\n"
)

func Run(cmd string, args []string, stdin <-chan string, ctrl chan os.Signal) (<-chan string, <-chan string, <-chan error) {
	// handler for the active command we'll be launching
	active := exec.Command(cmd, args...)

	// channels we use to give stdout, stderr and go errors to the user
	retStdout := make(chan string)
	retStderr := make(chan string)
	retErr := make(chan error)

	// stdout & stderr pipes coming from the active process
	cmdStdout, _ := active.StdoutPipe()
	cmdStderr, _ := active.StderrPipe()
	// stdin pipe to the active process
	cmdStdIn, _ := active.StdinPipe()

	// an error chan & helper func to kick off read pumps for stderr & stdout
	pumpErrs := make(chan error)
	pump(cmdStdout, retStdout, pumpErrs, msgEnds) // messages divided by \n\n
	pump(cmdStderr, retStderr, pumpErrs, "\n")    // return any error lines

	go func() {
		defer close(retStdout)
		defer close(retStderr)
		defer close(retErr)
		defer close(pumpErrs)
		defer cmdStdout.Close()
		defer cmdStderr.Close()
		defer cmdStdIn.Close()

		pumpsFinished := 0
		for {
			// the read pumps launched above handle reading from the active commands
			// stdout/err and returning messages to the caller, we're going to:
			// - check our read pumps haven't exited & exit if they both finish
			// - write to the process stdin if the caller sends input
			// - signal the process to exit if the caller sends any signal
			// If we encounter any errors they'll be shipped over our errs return chan
			select {
			case input := <-stdin:
				if input == "" {
					continue
				}

				_, err := cmdStdIn.Write([]byte(input))
				if err != nil {
					retErr <- err
				}
			case <-ctrl:
				retErr <- active.Process.Kill()
				return
			case err := <-pumpErrs:
				pumpsFinished += 1
				if err != nil {
					retErr <- err
				}
				if pumpsFinished >= 2 {
					ctrl <- syscall.SIGINT
				}
			}
		}
	}()

	err := active.Start()
	if err != nil {
		retErr <- err
		ctrl <- syscall.SIGINT
	}

	return retStdout, retStderr, retErr
}

func determineMsgs(in, sep string) ([]string, string) {
	bits := strings.Split(in, sep)

	if len(bits) == 1 {
		return []string{}, in
	}

	return bits[0 : len(bits)-1], bits[len(bits)-1]
}

func pump(src io.Reader, drain chan<- string, errs chan<- error, sep string) {
	go func() {
		soFar := ""
		for {
			buf := make([]byte, 1024)
			_, err := src.Read(buf)

			if err != nil {
				if err == io.EOF {
					errs <- nil
				} else {
					errs <- err
				}
			}

			soFar += string(buf)

			msgs, remaining := determineMsgs(soFar, sep)
			for _, msg := range msgs {
				drain <- msg
				soFar = remaining
			}
		}
	}()
}
