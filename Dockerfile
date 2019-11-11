FROM golang:alpine as builder
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk update && apk add --no-cache git
WORKDIR /go/src
COPY . /go/src/udpp/
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io
RUN go build -o /go/bin/udpp udpp/cmd/udpp/main.go

FROM alpine
COPY --from=builder /go/bin/udpp /go/bin/udpp
WORKDIR /go/bin
COPY config.toml /go/bin

ENTRYPOINT ["/go/bin/udpp"]
