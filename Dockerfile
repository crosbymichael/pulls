FROM ubuntu:12.04

RUN apt-get update && apt-get install -yq curl git

RUN curl -s https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz | tar -v -C /usr/local -xz
ENV PATH /usr/local/go/bin:/go/bin:$PATH
ENV GOPATH /go

ENTRYPOINT ["bash"]

WORKDIR /go/src
ADD . /go/src/github.com/crosbymichael/pulls
RUN cd /go/src/github.com/crosbymichael/pulls/pulls && go get && go install
