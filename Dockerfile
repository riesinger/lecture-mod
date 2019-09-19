FROM golang:1-alpine AS builder
RUN mkdir /build
WORKDIR /build
ENV GOOS=linux
ENV CGO_ENABLED=0
ENV GOPATH=/go
COPY . /build
RUN apk add --update --no-cache git && \
		go get -v && \
		go vet -v && \
		go build -v -a -installsuffix cgo -o lecture-mod .

# We are using alpine for CA Certificates
FROM alpine:latest AS runner

RUN mkdir /app

COPY --from=builder /build/lecture-mod /app
COPY --from=builder /build/client /app/client

RUN apk update && apk add ca-certificates tzdata

WORKDIR /app

EXPOSE 10944
ENTRYPOINT /app/lecture-mod
