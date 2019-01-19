FROM golang:1.11.3-alpine AS builder
LABEL maintainer="fabian.gehrlicher@outlook.de"

RUN apk add --no-cache ca-certificates git curl
RUN curl -fsSL -o /usr/local/bin/dep \
    https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
    chmod +x /usr/local/bin/dep

COPY . /go/src/gitlab.osram.info/osram/deployer
WORKDIR /go/src/gitlab.osram.info/osram/deployer

RUN dep ensure
RUN CGO_ENABLED=0 go build -o /deployer .

FROM alpine:3.6
LABEL maintainer="fabian.gehrlicher@outlook.de"

RUN apk --update add openssh-client git
RUN apk --no-cache add ca-certificates && update-ca-certificates
COPY --from=builder /deployer /deployer

ENV CLONE_METHOD="ssh"

ENTRYPOINT ["/deployer"]
