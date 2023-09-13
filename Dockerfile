FROM alpine:edge

RUN apk add tzdata

COPY dist /app
COPY model/*.json /app/model

WORKDIR /app

EXPOSE 3000

CMD [ "./main" ]
