# ![RealWorld Example App](https://user-images.githubusercontent.com/25560203/128340208-d07d731e-883c-46df-bde3-e236fb326d24.png)  
![workflow](https://github.com/zacscoding/echo-gorm-realworld-app/actions/workflows/check.yaml/badge.svg)

> ### Go/Echo/GORM codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.

### [Demo](https://github.com/gothinkster/realworld)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)

This codebase was created to demonstrate a fully fledged fullstack application built with **[YOUR_FRAMEWORK]** including
CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the **[YOUR_FRAMEWORK]** community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to
the [RealWorld](https://github.com/gothinkster/realworld) repo.

# How it works

![Simple Architecture](https://user-images.githubusercontent.com/25560203/128342030-bfeafe65-cf90-4856-90ef-65e345645d39.png)

**Note**: cache layer is not implemented yet.

API Server technology stack is

- Server code: **golang**
- Database: **MySQL**
- Migrate: **[golang-migrate/migrate](github.com/golang-migrate/migrate)**
- ORM: **[go-gorm/gorm v2](https://github.com/go-gorm/gorm)**
- Logging: [uber-go/zap](https://github.com/uber-go/zap)
- Unit Testing: **go test**, **[stretchr/testify](https://github.com/stretchr/testify)**
- Integration Testing: **[newman](https://github.com/postmanlabs/newman)**
- Configuration management: **[knadh/koanf](github.com/knadh/koanf)**

# Getting started

## Using docker-compose

```shell
// build app-server and start mysql, app-server. 
$ make compose.up
./scripts/compose.sh up
Creating network "echo-gorm-realworld-app_default" with the default driver
Creating db ... done
Creating app-server ... done
...

$ docker ps -a
CONTAINER ID   IMAGE                                COMMAND                  CREATED          STATUS          PORTS                                                    NAMES
01c48db0a1ff   zacscoding/echo-gorm-realworld-app   "app-server --config…"   18 minutes ago   Up 17 minutes   0.0.0.0:8080->8080/tcp, :::8080->8080/tcp                app-server
79db7c93b358   mysql:8.0.17                         "docker-entrypoint.s…"   18 minutes ago   Up 18 minutes   33060/tcp, 0.0.0.0:43306->3306/tcp, :::43306->3306/tcp   db
```

This server will serve on localhost:8080. Schema also migrated
from [golang-migrate/migrate](github.com/golang-migrate/migrate) and [migrations](./migrations).

- check api specs at http://localhost:8080/docs in ur browser.
- see [docker-compose.yaml](./docker-compose.yaml) for more info.
- see [config-docker.yaml](./config-docker.yaml) for configuration.
- see [migrations](./migrations) for schema.

# More commands

## Tests and checks lint, build

```shell
// This command includes clean tests cache, unit test, datarace, build, lint.
$ make tests

// also can be run each.

// run unit test
$ make test

// run datarace
$ make test.datarace

// check build
$ make test.build

// check lint
$ make lint
```

## Integration tests

After run servers(e.g: `make compose.up`), u can run integration tests.

```shell
$ make it.postman

+++ dirname integration/postman/run-api-tests.sh
++ cd integration/postman

...

┌─────────────────────────┬───────────────────┬──────────────────┐
│                         │          executed │           failed │
├─────────────────────────┼───────────────────┼──────────────────┤
│              iterations │                 1 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│                requests │                32 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│            test-scripts │                48 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│      prerequest-scripts │                18 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│              assertions │               263 │                0 │
├─────────────────────────┴───────────────────┴──────────────────┤
│ total run duration: 17.8s                                      │
├────────────────────────────────────────────────────────────────┤
│ total data received: 5.46KB (approx)                           │
├────────────────────────────────────────────────────────────────┤
│ average response time: 27ms [min: 7ms, max: 120ms, s.d.: 30ms] │
└────────────────────────────────────────────────────────────────┘
```  

---  

# TODO

- [ ] implement cache([go-redis/redis](https://github.com/go-redis/redis))
- [ ] unit tests of article handler
- [ ] integration tests with golang and [gavv/httpexpect](https://github.com/gavv/httpexpect)
