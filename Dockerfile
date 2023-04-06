FROM alpine:edge

RUN apk add tzdata

COPY dist /app
WORKDIR /app

EXPOSE 9000

CMD [ "./main" ]
