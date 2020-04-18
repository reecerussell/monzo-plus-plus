FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh ca-certificates openssl

COPY service.auth/ /go/src/github.com/reecerussell/monzo-plus-plus/service.auth/
COPY libraries/ /go/src/github.com/reecerussell/monzo-plus-plus/libraries/

WORKDIR /go/src/github.com/reecerussell/monzo-plus-plus/service.auth

## Install Dependencies
RUN go get github.com/google/uuid
RUN go get github.com/go-sql-driver/mysql
RUN go get google.golang.org/grpc
RUN go get golang.org/x/crypto/pbkdf2

ARG JWT_PRIVATE_KEY_PASS=v3rYS3cR3tKEY
ARG JWT_PRIVATE_KEY_NAME=jwt-private-key.pem

RUN openssl genrsa \
    -passout pass:${JWT_PRIVATE_KEY_PASS} \
    -out ${JWT_PRIVATE_KEY_NAME} \
    2048

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main cmd/main.go

FROM scratch

ARG JWT_PRIVATE_KEY_PASS=v3rYS3cR3tKEY
ARG JWT_PRIVATE_KEY_NAME=jwt-private-key.pem

COPY --from=builder /go/src/github.com/reecerussell/monzo-plus-plus/service.auth/main  /app/
COPY --from=builder /go/src/github.com/reecerussell/monzo-plus-plus/service.auth/${PRIVATE_KEY_NAME}  /app/
COPY --from=builder /go/src/github.com/reecerussell/monzo-plus-plus/service.auth/jwt-config.json  /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

EXPOSE 80

ENV JWT_CONFIG=jwt-config.json
ENV JWT_PRIVATE_KEY=${JWT_PRIVATE_KEY_NAME}
ENV JWT_PRIVATE_KEY_PASS=${JWT_PRIVATE_KEY_PASS}

CMD ["./main"]