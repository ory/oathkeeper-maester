# Build the manager binary
FROM alpine



RUN apk add -U --no-cache ca-certificates

FROM scratch

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY oathkeeper-k8s-controller /usr/bin/oathkeeper-k8s-controller
USER 1000

ENTRYPOINT ["oathkeeper-k8s-controller"]
