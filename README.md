# Pipe

Lightweight queryable windowed aggregator.

## Specification

Basic functionality include:

- In-memory event queue (a.k.a. pipe) for each received key
- Events in the pipe are sorted by their arrival times
- Old events (events that fall beyond the window) are evicted from the queue
- Metrics are collected (count)
- Queue is queryable

## Developing

To generate schemas, run:

```bash
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
protoc proto/pipe.proto --go_out=plugins=grpc:schema --proto_path=proto
```
