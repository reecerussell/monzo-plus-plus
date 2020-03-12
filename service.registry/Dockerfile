FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh ca-certificates

COPY service.registry/ /go/src/github.com/reecerussell/monzo-plus-plus/service.registry/

WORKDIR /go/src/github.com/reecerussell/monzo-plus-plus/service.registry

## Install Dependencies
RUN go get google.golang.org/grpc

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main cmd/main.go

FROM scratch

COPY --from=builder /go/src/github.com/reecerussell/monzo-plus-plus/service.registry/main  /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

CMD ["./main"]