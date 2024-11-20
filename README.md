# Hotel-reserv
Hotel reservation API using GO, fiber, and mongodb

> [!IMPORTANT]
> API docs with all the endpoints and signatures:
> [Docs/App_routes.md](https://github.com/Jucaza1/hotel-reserv/blob/main/docs/App_routes.md)

## Cloning repo
```bash
git clone https://github.com/jucaza1/hotel-reserv && cd hotel-reserv
```

---

## A) Build image and compose up both mongodb and hotel API
```bash
make docker-compose
```
By default on [http://localhost:4000](http://localhost:4000)
To change it edit docker-compose.yaml 4000:4000 -> XXXX:4000

---
## B.1a) Set up api with docker
```bash
make docker-api
```
## B.1b) Locally, needs to set up mongodb (see B.2)
```bash
make build

make run

make test
```
## B.2) Set up mongo with docker
### From makefile
```bash
make docker-mongo
```
### From bash
General config
```bash
docker run --name mongodb -p 27017:27017 -d mongo:latest
```
Or using this blueprint command
```bash
docker run -d --name YOUR_CONTAINER_NAME_HERE -p YOUR_LOCALHOST_PORT_HERE:27017 -e MONGO_INITDB_ROOT_USERNAME=YOUR_USERNAME_HERE -e MONGO_INITDB_ROOT_PASSWORD=YOUR_PASSWORD_HERE mongo
```
