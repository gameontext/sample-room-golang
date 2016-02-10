#
# Build from the go source directory:
#   docker build -t goroom:1.5 .

FROM golang:1.5

# We need gorilla, so go get it.
RUN cd $GOPATH && go get github.com/gorilla/websocket

# Copy our go source files and build/install the
# code.  If successful, the executable file,
# $GOPATH/bin/gameon-room-go, will be built.
RUN mkdir -p $GOPATH/src/gameon-room-go
COPY ./*.go $GOPATH/src/gameon-room-go/
COPY ./container-startup.sh /usr/bin/container-startup.sh
RUN cd $GOPATH/src/gameon-room-go && go install

# Our room should always listen on port 3000 (-lp)
# although the mapped host callback port (-cp) may
# be different.
EXPOSE 3000
WORKDIR $GOPATH/gameon-room-go/src

# The real work of running the game is done in the startup script
# which reads environment variables to drive its choices for startup.
ENTRYPOINT ["/bin/bash", "/usr/bin/container-startup.sh"]
