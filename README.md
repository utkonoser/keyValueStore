# keyValueStore !


Storage of key value pairs.
In this pet project, the following features are implemented:
- HTTP server using gorilla/mux
- REST service
- Transaction logger
- Saving values in the database (PostgreSQL)
- Containerization with Docker

## Installation

clone this repo:
```shell
git clone git@github.com:utkonoser/keyValueStore.git
```
and run the command:
```shell
docker compose up
```

## Environment

Add the configuration to the `.env`file at the root of the project:

```shell
DB_NAME=postgres
DB_PASSWORD=1234qwer 
HOST=localhost       #for Docker use the name of the db container 
DB_USER=postgres
```

## Examples

```shell
$ curl -X PUT -d 'Hello, KVS!' -v http://0.0.0.0:8080/v1/key-a
*   Trying 0.0.0.0:8080...
* Connected to 0.0.0.0 (127.0.0.1) port 8080 (#0)
> PUT /v1/key-a HTTP/1.1
> Host: 0.0.0.0:8080
> User-Agent: curl/7.81.0
> Accept: */*
> Content-Length: 11
> Content-Type: application/x-www-form-urlencoded
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 201 Created
< Date: Wed, 07 Dec 2022 12:25:53 GMT
< Content-Length: 0

$ curl -X GET -v http://0.0.0.0:8080/v1/key-a
*   Trying 0.0.0.0:8080...
* Connected to 0.0.0.0 (127.0.0.1) port 8080 (#0)
> GET /v1/key-a HTTP/1.1
> Host: 0.0.0.0:8080
> User-Agent: curl/7.81.0
> Accept: */*

* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Wed, 07 Dec 2022 12:26:20 GMT
< Content-Length: 11
< Content-Type: text/plain; charset=utf-8
< 
* Connection #0 to host 0.0.0.0 left intact
Hello,KVS!

$ curl -X DELETE -v http://0.0.0.0:8080/v1/key-a
*   Trying 0.0.0.0:8080...
* Connected to 0.0.0.0 (127.0.0.1) port 8080 (#0)
> DELETE /v1/key-a HTTP/1.1
> Host: 0.0.0.0:8080
> User-Agent: curl/7.81.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Wed, 07 Dec 2022 12:26:34 GMT
< Content-Length: 0
```
PostgreSQL:
```postgresql
postgres=# SELECT * FROM transactions;
 sequence | event_type |  key  |      value      
----------+------------+-------+-----------------
        1 |          2 | key-a | Hello%2C+KVS%21
        2 |          1 | key-a | 
(2 rows)
```


