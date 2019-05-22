# We are using alpine for CA Certificates
FROM alpine:latest

RUN mkdir /app

COPY lecture-mod /app
COPY client /app/client

RUN apk update && apk add ca-certificates tzdata

WORKDIR /app

EXPOSE 10944
ENTRYPOINT /app/lecture-mod
