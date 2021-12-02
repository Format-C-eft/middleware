# Builder

ARG GITHUB_PATH=github.com/Format-C-eft/middleware

FROM golang:1.17-alpine AS builder
RUN apk add --update make git curl
COPY . /home/${GITHUB_PATH}
WORKDIR /home/${GITHUB_PATH}
RUN make build-docker

# middleware

FROM alpine:latest as server
RUN apk add --update curl
LABEL org.opencontainers.image.source https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/middleware-docker .
COPY --from=builder /home/${GITHUB_PATH}/config.yml .

RUN chown root:root middleware-docker

EXPOSE 8000

CMD ["./middleware-docker"]
