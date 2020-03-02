FROM golang:latest as goimage
ENV GOPATH=/usr/go
ENV GOBIN=$GOPATH/bin
RUN mkdir -p $GOPATH
WORKDIR $GOPATH
COPY . .
RUN go build -o bin/main.exe src/main.go

FROM alpine:3.6 as baseimagealp
RUN apk update && apk add bash && apk add --no-cache bash
ENV WORK_DIR=/docker/bin
WORKDIR $WORK_DIR
COPY --from=goimage /usr/go/bin/main.exe .
ENTRYPOINT /bin/bash
EXPOSE 9908