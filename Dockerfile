FROM alpine:edge

WORKDIR /app

RUN apk add tzdata & mkdir /app/model

COPY dist /app
ADD model/*.json /app/model/
