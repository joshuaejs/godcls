# gocls

`gocls` - a distributed commit log service written in Go.

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
