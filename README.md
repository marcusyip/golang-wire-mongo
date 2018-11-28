# golang-wire-mongo

Attempt to use wire for dependency injection for app running and test container

* Wire (Dependency injection) (https://github.com/google/go-cloud/tree/master/wire)
* Gin
* Mongo (mongo-go-driver 0.0.18) (https://github.com/mongodb/mongo-go-driver)
* ginkgo (e2e)

## Generate wire_gen.go

```
make generate
```

## Running e2e test

```
docker-compose up
make e2e
```
