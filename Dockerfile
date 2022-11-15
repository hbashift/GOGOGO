FROM golang:alpine as builder

RUN go version
ENV GOPATH=/app
WORKDIR /src/app
COPY . /src/app

RUN go mod download
RUN go build -o server ./cmd/
CMD ["./server"]

