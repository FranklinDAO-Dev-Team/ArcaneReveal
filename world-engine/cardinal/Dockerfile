# Base Builder Stage
FROM golang:1.22 AS builder

WORKDIR /usr/src/app
COPY . .

RUN export GOPATH="$HOME/go"
RUN PATH="$GOPATH/bin:$PATH"

RUN make setup-vendor-docker
RUN go build -v -o /usr/local/bin/cardinal

CMD ["cardinal"]
