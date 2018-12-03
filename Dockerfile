FROM golang:1.11 as builder

RUN go get -u github.com/dakalab/translate

FROM gcr.io/distroless/base
COPY --from=builder /go/bin/translate /

ENTRYPOINT ["/translate"]
