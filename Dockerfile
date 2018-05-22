# We are using alpine for CA Certificates
FROM alpine:latest

COPY lecture-mod /usr/bin
COPY client /usr/bin

RUN apk update && apk add ca-certificates

EXPOSE 9988
CMD ["lecture-mod"]
