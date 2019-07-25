FROM golang:latest

 RUN apt-get update \
  && apt-get install -y curl unzip

RUN mkdir -p /tmp/protoc && \
  curl -L https://github.com/google/protobuf/releases/download/v3.0.0/protoc-3.0.0-linux-x86_64.zip > /tmp/protoc/protoc.zip && \
  cd /tmp/protoc && \
  unzip protoc.zip && \
  cp /tmp/protoc/bin/protoc /usr/local/bin && \
  chmod go+rx /usr/local/bin/protoc && \
  cd /tmp && \
  rm -r /tmp/protoc

WORKDIR /go/src/todo

RUN go get -u google.golang.org/grpc && \
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway && \
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger && \
    go get -u firebase.google.com/go

ENV GOOGLE_APPLICATION_CREDENTIALS /go/src/todo/service-account-file.json
