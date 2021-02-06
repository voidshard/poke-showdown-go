FROM node

# as of node-15 or so we need to run npm in a non-root dir
WORKDIR /opt/app
RUN chown -R node:node /opt/app/

# we want pokemon-showdown to be available to the user
USER node
RUN npm install pokemon-showdown
COPY app .

# symlink pokemon-showdown to /usr/bin
USER root
RUN ln -s /opt/app/node_modules/pokemon-showdown/pokemon-showdown /usr/bin/
USER node

EXPOSE 8080
CMD ["/opt/app/app"]
