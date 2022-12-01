FROM golang:1.19.2-alpine3.16 as builder
WORKDIR /usr/merchSearch/polling_agent/src
COPY /src .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app .

FROM alpine
WORKDIR /root/
COPY --from=0 /usr/merchSearch/polling_agent/src/app .
CMD ["./app"]