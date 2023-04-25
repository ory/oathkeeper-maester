# Build the manager binary
FROM golang:1.17 as builder
WORKDIR /go/src/app
COPY . .
RUN apt update &&\
    apt install ca-certificates &&\
    make manager

FROM gcr.io/distroless/static:latest
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/app/manager .
USER 1000

ENTRYPOINT ["/manager"]
