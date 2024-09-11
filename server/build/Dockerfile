FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /builder

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /builder/supplier main.go

FROM node:20-alpine3.19 AS final

WORKDIR /root

COPY --from=builder /builder/supplier .
COPY --from=builder /builder/internal/config/*.toml ./internal/config/
COPY --from=builder /builder/internal/third-party/eslint/eslint.config.mjs ./internal/third-party/eslint/

RUN npm i eslint --prefix ./internal/third-party/eslint

EXPOSE 8080

CMD ["./supplier"]