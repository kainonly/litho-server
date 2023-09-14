FROM alpine:edge

WORKDIR /app

RUN apk add tzdata & mkdir /app/model

ADD server /app/
ADD model/*.json /app/model/
