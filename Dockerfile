# Build the manager binary
FROM alpine

RUN apk add -U --no-cache ca-certificates

FROM scratch

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY bin/manager .
USER 1000

ENTRYPOINT ["/manager"]
