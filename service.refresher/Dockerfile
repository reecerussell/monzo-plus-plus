FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh ca-certificates

COPY . /go/src/github.com/reecerussell/monzo-plus-plus/service.refresher/

WORKDIR /go/src/github.com/reecerussell/monzo-plus-plus/service.refresher

RUN go get github.com/go-sql-driver/mysql

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch

COPY --from=builder /go/src/github.com/reecerussell/monzo-plus-plus/service.refresher/main  /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

CMD ["./main"]