FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=LINUX GOARCH=amd64 go build -a -installsuffix cgo -o /app/bin/main .

FROM alpine:latest

RUN apk add --no-cache ca-certificates \
	tzdata

ENV TZ Asia/Jakarta

COPY --from=builder /app/bin/main .

EXPOSE 80

CMD ["./main", "rest", "-p", "80"]
