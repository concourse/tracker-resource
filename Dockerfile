FROM ubuntu:14.04

RUN apt-get update && apt-get -y install wget git

ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /tmp/go

ADD . /tmp/go/src/github.com/concourse/tracker-resource

RUN wget -qO- https://storage.googleapis.com/golang/go1.3.linux-amd64.tar.gz | tar -C /usr/local -xzf - && \
      go build -o /opt/resource/out github.com/concourse/tracker-resource/out/cmd/out && \
      rm -rf /tmp/go /usr/local/go
