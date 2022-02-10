FROM golang:latest as builder

ENV CGO_ENABLED=0

WORKDIR /app
COPY . .
RUN go build -o ./bin/ddns-updater ./cmd/ddns-updater

FROM scratch
COPY --from=builder /app/bin/ddns-updater /go/bin/ddns-updater
ENTRYPOINT ["/go/bin/ddns-updater"]
