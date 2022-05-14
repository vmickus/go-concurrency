# GO concurrency

Experiment on how to handle concurrency in real projects.

Suppose that, for each request, we have to access N times an upstream that can take too long to answer (s3 for example).

We need:

* Each work upstream request should be done in its own goroutine
* We should have a timeout, in case an upstream takes too long to return
* Collect the result of each upstream

## Run

```shell
go run ./...
```

## Test

```shell
go test -v ./...
```
