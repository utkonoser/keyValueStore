# keyValueStore !
<small>[is under development]</small>

Storage of key value pairs.
In this pet project, the following features are implemented:
- HTTP server using gorilla/mux
- REST service
- Transaction log
- Saving values in the database
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



