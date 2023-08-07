FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/anagram_finder/main.go

FROM alpine:latest  

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
