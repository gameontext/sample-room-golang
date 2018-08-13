#
# Build from the go source directory:
#   docker build -t goroom:1.5 .

FROM golang:1.10

# We need gorilla, so go get it.
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN cd $GOPATH

# Copy our go source files and build/install the
# code.  If successful, the executable file,
# $GOPATH/bin/sample-room-golang, will be built.
RUN mkdir -p $GOPATH/src/sample-room-golang
# COPY . $GOPATH/src/sample-room-golang/
COPY ./*.go $GOPATH/src/sample-room-golang/
COPY ./routers/ $GOPATH/src/sample-room-golang/routers/
COPY ./plugins/ $GOPATH/src/sample-room-golang/plugins/
COPY ./Gopkg.toml $GOPATH/src/sample-room-golang/
COPY ./Gopkg.lock $GOPATH/src/sample-room-golang/
COPY ./container-startup.sh /usr/bin/container-startup.sh
RUN cd $GOPATH/src/sample-room-golang && dep ensure
RUN cd $GOPATH/src/sample-room-golang && go install

# Our room should always listen on port 3000 (-lp)
# although the mapped host callback port (-cp) may
# be different.
EXPOSE 3000
WORKDIR $GOPATH/src/sample-room-golang

# The real work of running the game is done in the startup script
# which reads environment variables to drive its choices for startup.
ENTRYPOINT ["/bin/bash", "/usr/bin/container-startup.sh"]