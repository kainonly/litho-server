FROM golang:1.26-alpine AS builder

WORKDIR /src

ENV CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOPROXY=https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.cloud.tencent.com/g' /etc/apk/repositories

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./dist/server ./

FROM alpine:edge

WORKDIR /app

RUN apk --no-cache add tzdata

COPY --from=builder /src/dist/server .
COPY --from=builder /src/config/resources.yml ./config/

EXPOSE 3000

CMD ["./server"]
