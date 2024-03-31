# Hotel-reserv
Hotel reservation API using GO

## MAKEFILE commands
```
make build

make run

make test

make docker
```
## Set up mongo with docker
### From makefile
```
make docker
```
### From bash
General config
```
docker run --name mongodb -p 27017:27017 -d mongo:latest
```
Or using this blueprint command
```
docker run -d --name YOUR_CONTAINER_NAME_HERE -p YOUR_LOCALHOST_PORT_HERE:27017 -e MONGO_INITDB_ROOT_USERNAME=YOUR_USERNAME_HERE -e MONGO_INITDB_ROOT_PASSWORD=YOUR_PASSWORD_HERE mongo
```