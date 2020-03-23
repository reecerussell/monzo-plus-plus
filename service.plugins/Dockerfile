FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh ca-certificates openssl

COPY service.plugins/ /go/src/github.com/reecerussell/monzo-plus-plus/service.plugins/
COPY libraries/ /go/src/github.com/reecerussell/monzo-plus-plus/libraries/

WORKDIR /go/src/github.com/reecerussell/monzo-plus-plus/service.plugins

## Install Dependencies
RUN go get github.com/google/uuid
RUN go get github.com/go-sql-driver/mysql
RUN go get google.golang.org/grpc
RUN go get github.com/gorilla/mux

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main cmd/main.go

FROM scratch

COPY --from=builder /go/src/github.com/reecerussell/monzo-plus-plus/service.plugins/main  /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

EXPOSE 80

CMD ["./main"]