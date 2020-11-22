package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDetermineMsgs(t *testing.T) {
	cases := []struct {
		Given        string
		Sep          string
		ExpectParsed []string
		ExpectRemain string
	}{
		{"a\n\nb\n\ncd\n", "\n\n", []string{"a", "b"}, "cd\n"},
		{"a\nb\n\nc", "\n", []string{"a", "b", ""}, "c"},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			messages, remaining := determineMsgs(tt.Given, tt.Sep)

			assert.Equal(t, tt.ExpectParsed, messages)
			assert.Equal(t, tt.ExpectRemain, remaining)
		})
	}

}
