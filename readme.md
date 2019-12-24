# gokit ddd demo

## api gateway
* http server listen on :1323

## order service
* grpc server listen on :8092
* http server listen on :8091

## user service
* grpc server listen on :8082
* http server listen on :8081

## request flow

launch request:
```shell script
curl http://127.0.0.1:1323/api/users?with_orders=true
```

1. `curl` launches http request to `api gateway`
2. `api gateway` lauches grpc request to `user service`
3. `api gateway` lauches grpc request to `order service` if with_orders is true
