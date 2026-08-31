FROM golang:1.26 AS builder

RUN mkdir server
WORKDIR /server

# go mod download
COPY . .
RUN go mod download
RUN echo "go mod download -- success --"


# go build
RUN go build -o server ./cmd
RUN echo "go build -- success --"

FROM debian:trixie-slim

RUN mkdir server
WORKDIR /server

COPY --from=builder /server/server ./server
RUN echo "copy .server -- success --"


CMD ./server