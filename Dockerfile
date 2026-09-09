FROM golang:1.26-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOARCH=amd64
ENV GOOS=linux
ENV GOPROXY=https://goproxy.cn,direct

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.cloud.tencent.com/g' /etc/apk/repositories
RUN apk --no-cache add tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY ./config/values.build.yml ./config/values.yml

RUN go run ./builder
RUN go build -o ./dist/server ./

FROM alpine:edge

WORKDIR /app

RUN apk --no-cache add tzdata

COPY --from=builder /app/dist/server .

CMD ["./server"]
