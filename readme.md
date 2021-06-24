# gokit ddd demo

## api gateway
### bulid
```
cd api_gateway
go build
```

* http server listen on :1323

## order service
### build
```
cd order_svc/app/ordersvc
go build
```

* grpc server listen on :8092
* http server listen on :8091

## user service
### build
```
cd user_svc/app/usersvc
go build
```

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
