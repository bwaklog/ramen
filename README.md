## Ramen

A back-end for testing load balancers. The project consists of a compose file `compose.yml` that spins up a given amount of docker containers for the back-end (default of 3) as part of the `web` service. 

The `nginx` service defines a single Nginx container that sits infront of the 3 containers which send requests to a redis backend that acts as the storage

![service structure](https://i.imgur.com/n6YxvTd.png)

---

**Project Structure**

```text
├── Dockerfile
├── compose.yml
├── go.mod
├── go.sum
├── main.go 
├── nginx
│   └── nginx.conf 
└── pkg
    └── redis.go // all packages and methods 
                 // to interface with the redis backend
```

---

Requirements
1. [Docker Desktop](https://docs.docker.com/desktop/) for building the image and using compose
2. [GoLang](https://go.dev) version 1.23.2

---

Building

You need to use the compose file for testing, unless you decide to run all services locally. Here is a short guide to how to use compose with the given `compose.yml`

```sh
# use docker compose to start services
docker compose up

# to compose and de-attach the service logs from terminal
docker compose up -d 

# stop all services
docker compose down
```

