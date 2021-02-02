FROM golang:1.15-alpine as builder

WORKDIR /app

COPY . .

RUN env GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o svc -mod=vendor -a -installsuffix cgo -ldflags '-w'

FROM alpine
RUN apk --update --no-cache add ca-certificates && \
	addgroup -S svc && adduser -S -g svc svc
USER svc
COPY --from=builder /app/svc / 
COPY --chown=svc:svc ./index.html /index.html
CMD ["/svc"]
EXPOSE 8000