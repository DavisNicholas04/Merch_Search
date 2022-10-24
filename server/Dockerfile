FROM golang:1.19.2-alpine3.16 as builder
WORKDIR /usr/merchSearch/server/src
COPY /src .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app .
#RUN apk --no-cache add ca-certificates

FROM alpine
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY --from=builder /etc/ssl/certs/ca-certificates.key /etc/ssl/certs/
COPY /34.207.90.86.crt /etc/ssl/certs/
COPY /34.207.90.86.key /etc/ssl/certs/
WORKDIR /root/
COPY --from=0 /usr/merchSearch/server/src/app .
CMD ["./app"]
