FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /builder

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /builder/supplier main.go

FROM alpine:latest AS final

WORKDIR /root

COPY --from=builder /builder/supplier .

EXPOSE 8080

CMD ["./supplier"]