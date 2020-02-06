FROM golang:1.13.7-alpine as builder

WORKDIR /app
COPY ./ /app
RUN go mod download
RUN go build -o ./builds/linux/explorer_genesis_uploader ./cmd/explorer_genesis_uploader.go

FROM alpine:3.7

COPY --from=builder /app/builds/linux/explorer_genesis_uploader /usr/bin/explorer_genesis_uploader
RUN addgroup minteruser && adduser -D -h /minter -G minteruser minteruser
USER minteruser
WORKDIR /minter
ENTRYPOINT ["/usr/bin/explorer_genesis_uploader"]
CMD ["explorer_genesis_uploader"]
