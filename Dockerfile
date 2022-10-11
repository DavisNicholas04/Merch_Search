FROM golang:1.19.2-alpine3.16 as builder
WORKDIR /usr/merchSearch/src
COPY /src .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app .
RUN apk --no-cache add ca-certificates

FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /root/
COPY --from=0 /usr/merchSearch/src/app .
CMD ["./app"]