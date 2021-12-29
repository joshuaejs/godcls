# gocls

`godcls` - a distributed commit log service written in Go.

## setup

- Ubuntu 21.10
  - install the protobuf compiler: `sudo apt-get update -qq && sudo apt-get install protobuf-compiler -y`
  - add to `go.mod`: `go get -d google.golang.org/protobuf/...@latest`
  - install `protoc-gen-go`: `go install fgoogle.golang.org/protobuf/...@latest`

```txt
$ sudo apt-get update -qq && sudo apt-get install protobuf-compiler -y
Reading package lists... Done
Building dependency tree... Done
Reading state information... Done
The following additional packages will be installed:
  libprotobuf-dev libprotobuf-lite23 libprotobuf23 libprotoc23
Suggested packages:
  protobuf-mode-el
The following NEW packages will be installed:
  libprotobuf-dev libprotobuf-lite23 libprotobuf23 libprotoc23 protobuf-compiler
0 upgraded, 5 newly installed, 0 to remove and 0 not upgraded.
Need to get 3127 kB of archives.
After this operation, 17.5 MB of additional disk space will be used.
Get:1 http://archive.ubuntu.com/ubuntu impish/main amd64 libprotobuf-lite23 amd64 3.12.4-1ubuntu3 [209 kB]
Get:2 http://archive.ubuntu.com/ubuntu impish/main amd64 libprotobuf23 amd64 3.12.4-1ubuntu3 [878 kB]
Get:3 http://archive.ubuntu.com/ubuntu impish/main amd64 libprotoc23 amd64 3.12.4-1ubuntu3 [664 kB]
Get:4 http://archive.ubuntu.com/ubuntu impish/main amd64 libprotobuf-dev amd64 3.12.4-1ubuntu3 [1347 kB]
Get:5 http://archive.ubuntu.com/ubuntu impish/universe amd64 protobuf-compiler amd64 3.12.4-1ubuntu3 [29.2 kB]
Fetched 3127 kB in 2s (1553 kB/s)
Selecting previously unselected package libprotobuf-lite23:amd64.
(Reading database ... 193676 files and directories currently installed.)
Preparing to unpack .../libprotobuf-lite23_3.12.4-1ubuntu3_amd64.deb ...
Unpacking libprotobuf-lite23:amd64 (3.12.4-1ubuntu3) ...
Selecting previously unselected package libprotobuf23:amd64.
Preparing to unpack .../libprotobuf23_3.12.4-1ubuntu3_amd64.deb ...
Unpacking libprotobuf23:amd64 (3.12.4-1ubuntu3) ...
Selecting previously unselected package libprotoc23:amd64.
Preparing to unpack .../libprotoc23_3.12.4-1ubuntu3_amd64.deb ...
Unpacking libprotoc23:amd64 (3.12.4-1ubuntu3) ...
Selecting previously unselected package libprotobuf-dev:amd64.
Preparing to unpack .../libprotobuf-dev_3.12.4-1ubuntu3_amd64.deb ...
Unpacking libprotobuf-dev:amd64 (3.12.4-1ubuntu3) ...
Selecting previously unselected package protobuf-compiler.
Preparing to unpack .../protobuf-compiler_3.12.4-1ubuntu3_amd64.deb ...
Unpacking protobuf-compiler (3.12.4-1ubuntu3) ...
Setting up libprotobuf23:amd64 (3.12.4-1ubuntu3) ...
Setting up libprotobuf-lite23:amd64 (3.12.4-1ubuntu3) ...
Setting up libprotoc23:amd64 (3.12.4-1ubuntu3) ...
Setting up protobuf-compiler (3.12.4-1ubuntu3) ...
Setting up libprotobuf-dev:amd64 (3.12.4-1ubuntu3) ...
Processing triggers for man-db (2.9.4-2) ...
Processing triggers for libc-bin (2.34-0ubuntu3) ...
/sbin/ldconfig.real: /usr/lib/wsl/lib/libcuda.so.1 is not a symbolic link
$ go get -d google.golang.org/protobuf
go get: added google.golang.org/protobuf v1.27.1
$ go install google.golang.org/protobuf/...@v1.27.1
```

## chapter 4

- updated `server_test.go` to fix deprecation warnings with grpc

```txt
$ diff

        "google.golang.org/grpc"
+       "google.golang.org/grpc/credentials/insecure"
+       "google.golang.org/grpc/status"


-       clientOptions := []grpc.DialOption{grpc.WithInsecure()}
+       clientOptions := []grpc.DialOption{
+               grpc.WithTransportCredentials(insecure.NewCredentials()),
+       }


-       got := grpc.Code(err)
-       want := grpc.Code(api.ErrOffsetOutOfRange{}.GRPCStatus().Err())
+       got := status.Code(err)
+       want := status.Code(api.ErrOffsetOutOfRange{}.GRPCStatus().Err())

```

## chapter 5

```txt
go install github.com/cloudflare/cfssl/cmd/cfssl@latest
go install github.com/cloudflare/cfssl/cmd/cfssljson@latest

mkdir test

mkdir -p internal/config

```

- an ACL is a table of rules
  - an ACL is easy to code, as it's *just a table*, can use a map or CSV for simple implementations; complex implementations would use a key-value store or relational db.

## chapter 6

- three types of telemetry data
  - observability
    - measure of how well understood the system's internals--its behavior and state--are from its external outputs: metrics, structured logs, and traces
  - metrics
    - measure numeric data over time.
    - this helps define:
      - service-level indicators (SLI), objectives (SLO), and agreements (SLA)
    - three kinds of metrics
      - counters
        - track the number of times an event happened
        - used to get a rate
      - histograms
        - show a distribution of data
        - used for percentiles
      - guages
        - track the current value of something
        - useful for saturation-type metrics, like a host's disk utilization
    - *Google's four golden signals* to measure:
      - latency
        - the time it takes the service to process a request
      - traffic
        - the amount of demand on the service
      - errors
        - the request failure rate of the service
        - internal server errors are especially important
      - saturation
        - a measure of the service's capacity
    - structured logs
      - a set of name and value ordered pairs encoded in a consistent schema and format
    - traces
      - capture request lifecycles and track request flows through the system
      - comprised of one or more *spans*
      - spans can have parent/child relationships, or be linked as siblings
      - go wide to begin
        - trace requests across all services end-to-end, with spans that being and end at the entry and exit points of the services
      - then go deep
        - trace important method calls

```txt
go get -d go.uber.org/zap@latest
go get -d go.opencensus.io@latest
```

```txt
$ cd internal/server/
$ go test -v -debug=true
=== RUN   TestServer
=== RUN   TestServer/produce/consume_stream_succeeds
    server_test.go:54: metrics log file: /tmp/metrics-518807290.log
    server_test.go:54: traces log file: /tmp/traces-951245228.log
2021-12-28T18:58:59.943-0800    INFO    server  zap/options.go:212      finished streaming call with code OK    {"grpc.start_time": "2021-12-28T18:58:59-08:00", "system": "grpc", "span.kind": "server", "grpc.service": "log.v1.Log", "grpc.method": "ConsumeStream", "peer.address": "127.0.0.1:38958", "grpc.code": "OK", "grpc.time_ns": 371700}
2021-12-28T18:58:59.943-0800    INFO    server  zap/options.go:212      finished streaming call with code Canceled      {"grpc.start_time": "2021-12-28T18:58:59-08:00", "system": "grpc", "span.kind": "server", "grpc.service": "log.v1.Log", "grpc.method": "ProduceStream", "peer.address": "127.0.0.1:38958", "error": "rpc error: code = Canceled desc = context canceled", "grpc.code": "Canceled", "grpc.time_ns": 1117400}
=== RUN   TestServer/consume_past_log_boundary_fail
    server_test.go:54: metrics log file: /tmp/metrics-601138573.log
    server_test.go:54: traces log file: /tmp/traces-3555482549.log
2021-12-28T18:59:01.449-0800    INFO    server  zap/options.go:212      finished unary call with code OK        {"grpc.start_time": "2021-12-28T18:59:01-08:00", "system": "grpc", "span.kind": "server", "grpc.service": "log.v1.Log", "grpc.method": "Produce", "peer.address": "127.0.0.1:59584", "grpc.code": "OK", "grpc.time_ns": 87300}
2021-12-28T18:59:01.450-0800    ERROR   server  zap/options.go:212      finished unary call with code Code(404) {"grpc.start_time": "2021-12-28T18:59:01-08:00", "system": "grpc", "span.kind": "server", "grpc.service": "log.v1.Log", "grpc.method": "Consume", "peer.address": "127.0.0.1:59584", "error": "rpc error: code = Code(404) desc = offset out of range: 1", "grpc.code": "Code(404)", "grpc.time_ns": 132300}
github.com/grpc-ecosystem/go-grpc-middleware/logging/zap.DefaultMessageProducer
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/logging/zap/options.go:212
github.com/grpc-ecosystem/go-grpc-middleware/logging/zap.UnaryServerInterceptor.func1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/logging/zap/server_interceptors.go:39
github.com/grpc-ecosystem/go-grpc-middleware.ChainUnaryServer.func1.1.1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/chain.go:25
github.com/grpc-ecosystem/go-grpc-middleware/tags.UnaryServerInterceptor.func1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/tags/interceptors.go:23
github.com/grpc-ecosystem/go-grpc-middleware.ChainUnaryServer.func1.1.1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/chain.go:25
github.com/grpc-ecosystem/go-grpc-middleware.ChainUnaryServer.func1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/chain.go:34
github.com/joshuaejs/godcls/api/v1._Log_Consume_Handler
        /home/jejs/repo/joshuaejs/godcls/api/v1/log_grpc.pb.go:193
google.golang.org/grpc.(*Server).processUnaryRPC
        /home/jejs/go/pkg/mod/google.golang.org/grpc@v1.43.0/server.go:1282
google.golang.org/grpc.(*Server).handleStream
        /home/jejs/go/pkg/mod/google.golang.org/grpc@v1.43.0/server.go:1616
google.golang.org/grpc.(*Server).serveStreams.func1.2
        /home/jejs/go/pkg/mod/google.golang.org/grpc@v1.43.0/server.go:921
=== RUN   TestServer/unauthorized_fails
    server_test.go:54: metrics log file: /tmp/metrics-3720234543.log
    server_test.go:54: traces log file: /tmp/traces-793760578.log
2021-12-28T18:59:02.960-0800    WARN    server  zap/options.go:212      finished unary call with code PermissionDenied  {"grpc.start_time": "2021-12-28T18:59:02-08:00", "system": "grpc", "span.kind": "server", "grpc.service": "log.v1.Log", "grpc.method": "Produce", "peer.address": "127.0.0.1:33854", "error": "rpc error: code = PermissionDenied desc = nobody not permitted to * to produce", "grpc.code": "PermissionDenied", "grpc.time_ns": 92400}
github.com/grpc-ecosystem/go-grpc-middleware/logging/zap.DefaultMessageProducer
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/logging/zap/options.go:212
github.com/grpc-ecosystem/go-grpc-middleware/logging/zap.UnaryServerInterceptor.func1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/logging/zap/server_interceptors.go:39
github.com/grpc-ecosystem/go-grpc-middleware.ChainUnaryServer.func1.1.1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/chain.go:25
github.com/grpc-ecosystem/go-grpc-middleware/tags.UnaryServerInterceptor.func1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/tags/interceptors.go:23
github.com/grpc-ecosystem/go-grpc-middleware.ChainUnaryServer.func1.1.1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/chain.go:25
github.com/grpc-ecosystem/go-grpc-middleware.ChainUnaryServer.func1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/chain.go:34
github.com/joshuaejs/godcls/api/v1._Log_Produce_Handler
        /home/jejs/repo/joshuaejs/godcls/api/v1/log_grpc.pb.go:175
google.golang.org/grpc.(*Server).processUnaryRPC
        /home/jejs/go/pkg/mod/google.golang.org/grpc@v1.43.0/server.go:1282
google.golang.org/grpc.(*Server).handleStream
        /home/jejs/go/pkg/mod/google.golang.org/grpc@v1.43.0/server.go:1616
google.golang.org/grpc.(*Server).serveStreams.func1.2
        /home/jejs/go/pkg/mod/google.golang.org/grpc@v1.43.0/server.go:921
2021-12-28T18:59:02.960-0800    WARN    server  zap/options.go:212      finished unary call with code PermissionDenied  {"grpc.start_time": "2021-12-28T18:59:02-08:00", "system": "grpc", "span.kind": "server", "grpc.service": "log.v1.Log", "grpc.method": "Consume", "peer.address": "127.0.0.1:33854", "error": "rpc error: code = PermissionDenied desc = nobody not permitted to * to produce", "grpc.code": "PermissionDenied", "grpc.time_ns": 62300}
github.com/grpc-ecosystem/go-grpc-middleware/logging/zap.DefaultMessageProducer
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/logging/zap/options.go:212
github.com/grpc-ecosystem/go-grpc-middleware/logging/zap.UnaryServerInterceptor.func1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/logging/zap/server_interceptors.go:39
github.com/grpc-ecosystem/go-grpc-middleware.ChainUnaryServer.func1.1.1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/chain.go:25
github.com/grpc-ecosystem/go-grpc-middleware/tags.UnaryServerInterceptor.func1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/tags/interceptors.go:23
github.com/grpc-ecosystem/go-grpc-middleware.ChainUnaryServer.func1.1.1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/chain.go:25
github.com/grpc-ecosystem/go-grpc-middleware.ChainUnaryServer.func1
        /home/jejs/go/pkg/mod/github.com/grpc-ecosystem/go-grpc-middleware@v1.3.0/chain.go:34
github.com/joshuaejs/godcls/api/v1._Log_Consume_Handler
        /home/jejs/repo/joshuaejs/godcls/api/v1/log_grpc.pb.go:193
google.golang.org/grpc.(*Server).processUnaryRPC
        /home/jejs/go/pkg/mod/google.golang.org/grpc@v1.43.0/server.go:1282
google.golang.org/grpc.(*Server).handleStream
        /home/jejs/go/pkg/mod/google.golang.org/grpc@v1.43.0/server.go:1616
google.golang.org/grpc.(*Server).serveStreams.func1.2
        /home/jejs/go/pkg/mod/google.golang.org/grpc@v1.43.0/server.go:921
=== RUN   TestServer/produce/consume_a_message_to/from_the_log_succeeds
    server_test.go:54: metrics log file: /tmp/metrics-2267114488.log
    server_test.go:54: traces log file: /tmp/traces-163317472.log
2021-12-28T18:59:04.469-0800    INFO    server  zap/options.go:212      finished unary call with code OK        {"grpc.start_time": "2021-12-28T18:59:04-08:00", "system": "grpc", "span.kind": "server", "grpc.service": "log.v1.Log", "grpc.method": "Produce", "peer.address": "127.0.0.1:47242", "grpc.code": "OK", "grpc.time_ns": 127000}
2021-12-28T18:59:04.470-0800    INFO    server  zap/options.go:212      finished unary call with code OK        {"grpc.start_time": "2021-12-28T18:59:04-08:00", "system": "grpc", "span.kind": "server", "grpc.service": "log.v1.Log", "grpc.method": "Consume", "peer.address": "127.0.0.1:47242", "grpc.code": "OK", "grpc.time_ns": 51700}
--- PASS: TestServer (6.04s)
    --- PASS: TestServer/produce/consume_stream_succeeds (1.51s)
    --- PASS: TestServer/consume_past_log_boundary_fail (1.51s)
    --- PASS: TestServer/unauthorized_fails (1.51s)
    --- PASS: TestServer/produce/consume_a_message_to/from_the_log_succeeds (1.51s)
PASS
ok      github.com/joshuaejs/godcls/internal/server     6.044s
```
