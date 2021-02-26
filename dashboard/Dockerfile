# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
WORKDIR /go/src/app


ADD . .

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)



RUN go get -u github.com/DronRathore/goexpress
RUN go get -u github.com/DronRathore/goexpress/request
RUN go get -u github.com/DronRathore/goexpress/response

RUN go get -u github.com/go-redis/redis
RUN go get -u github.com/googollee/go-socket.io
RUN go get -u gopkg.in/olivere/elastic.v5

RUN go get -u github.com/kidoman/embd


RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

CMD ["go-wrapper", "run"] # ["app"]


# Document that the service listens on port 8080.
EXPOSE 8080
