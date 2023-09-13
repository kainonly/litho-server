FROM alpine:edge

WORKDIR /app

RUN apk add tzdata & mkdir /app/model

COPY dist /app
COPY model/*.json /app/model

EXPOSE 3000

CMD [ "./main" ]
