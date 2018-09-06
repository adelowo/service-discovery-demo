#### Service discovery and Loadalancing with consul and fabio

```sh
$ consul agent -dev
```

```sh
$ go run main.go
$ # In another terminal window, run go run cmd/main.go -http ":3500"
$ # In another, go run cmd/main.go -http ":3900"
```

```sh
$ fabio
```
