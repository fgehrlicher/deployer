FROM golang:1.14-alpine AS builder
LABEL maintainer="fabian.gehrlicher@outlook.de"

RUN apk add --no-cache ca-certificates git curl
WORKDIR /deployer

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build .

FROM alpine:3.6
LABEL maintainer="fabian.gehrlicher@outlook.de"

RUN apk --update add openssh-client git
RUN apk --no-cache add ca-certificates && update-ca-certificates

COPY default_known_hosts /root/.ssh/known_hosts
ENV SSH_KNOWN_HOSTS="/root/.ssh/known_hosts"
ENV KEY_FILE_PATH="/root/.ssh/id_rsa"
ENV CLONE_METHOD="ssh"

COPY --from=builder /deployer/deployer /deployer

ENTRYPOINT ["/deployer"]
