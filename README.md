### Build and Run

1. Start the jaeger service

```sh
docker run -d -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest
```

2. Open the jaeger UI

http://localhost:16686

3. Run the server.go

```sh
go run server.go
```

4. Run the client.go (optional)

```sh
go run client.go
```
