FROM golang:latest AS builder

ENV DEP_VERSION=v0.3.2

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/${DEP_VERSION}/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/musicstalker
WORKDIR /go/src/github.com/musicstalker

COPY Gopkg.toml Gopkg.lock ./
# copies the Gopkg.toml and Gopkg.lock to WORKDIR

RUN dep ensure -vendor-only

ADD . .
RUN go install .

ENTRYPOINT [ "/go/bin/musicstalker" ]