FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN GOARCH=amd64 GOOS=linux go build -o app .

FROM amd64/alpine:latest

RUN apk add --no-cache libc6-compat

WORKDIR /root/

COPY --from=builder /app/app .
COPY --from=builder /app/src .

CMD ["./app"]
