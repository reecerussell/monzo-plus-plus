FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh ca-certificates openssl && \
    apk add --update npm

COPY service.mpp/ /go/src/github.com/reecerussell/monzo-plus-plus/service.mpp/
COPY libraries/ /go/src/github.com/reecerussell/monzo-plus-plus/libraries/
COPY ui/ /go/src/github.com/reecerussell/monzo-plus-plus/ui/

#
# BUILD GO
#

WORKDIR /go/src/github.com/reecerussell/monzo-plus-plus/service.mpp

## Install Dependencies
RUN go get github.com/google/uuid
RUN go get github.com/go-sql-driver/mysql
RUN go get google.golang.org/grpc
RUN go get github.com/gorilla/mux

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main cmd/main.go

# Generate SSL key and cert.
RUN openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 \
    -subj "/C=GB/ST=Buckinghamshire/L=Milton Keynes/O=IT/CN=auth.api" \
    -keyout ssl_key.pem  \
    -out ssl_cert.pem

#
# BUILD UI
#

WORKDIR /go/src/github.com/reecerussell/monzo-plus-plus/ui

RUN npm install
RUN npm audit fix
RUN npm run build

FROM scratch

COPY --from=builder /go/src/github.com/reecerussell/monzo-plus-plus/service.mpp/main  /app/
COPY --from=builder /go/src/github.com/reecerussell/monzo-plus-plus/service.mpp/ssl_key.pem  /app/ssl_key.pem
COPY --from=builder /go/src/github.com/reecerussell/monzo-plus-plus/service.mpp/ssl_cert.pem  /app/ssl_cert.pem
COPY --from=builder /go/src/github.com/reecerussell/monzo-plus-plus/ui/build /app/ui/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

ENV SSL_KEY_FILE=ssl_key.pem
ENV SSL_CERT_FILE=ssl_cert.pem

CMD ["./main"]