FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o litho-api

FROM alpine:edge

WORKDIR /app

RUN apk --no-cache add tzdata

COPY --from=builder /src/litho-api /app/

EXPOSE 3000

CMD ["./litho-api"]
