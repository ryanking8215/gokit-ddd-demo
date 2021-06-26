FROM golang:1.16 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.io,direct

WORKDIR /go/src/gokit-ddd-demo/
COPY . .
RUN go mod tidy

# build
RUN cd api_gateway && go build -a -installsuffix cgo -o api_gateway .
RUN cd user_svc/app/usersvc && go build -a -installsuffix cgo -o usersvc .
RUN cd order_svc/app/ordersvc && go build -a -installsuffix cgo -o ordersvc .

FROM alpine:latest AS api_gateway
RUN apk --no-cache add ca-certificates
WORKDIR /opt/gokit-ddd-demo/api_gateway
COPY --from=builder /go/src/gokit-ddd-demo/api_gateway/api_gateway .
EXPOSE 1323
CMD ["./api_gateway"]

FROM alpine:latest AS user_svc
RUN apk --no-cache add ca-certificates
WORKDIR /opt/gokit-ddd-demo/user_svc
COPY --from=builder /go/src/gokit-ddd-demo/user_svc/app/usersvc/usersvc .
EXPOSE 8081
EXPOSE 8082
CMD ["./usersvc"]

FROM alpine:latest AS order_svc
RUN apk --no-cache add ca-certificates
WORKDIR /opt/gokit-ddd-demo/order_svc
COPY --from=builder /go/src/gokit-ddd-demo/order_svc/app/ordersvc/ordersvc .
EXPOSE 8091
EXPOSE 8092
CMD ["./ordersvc"]