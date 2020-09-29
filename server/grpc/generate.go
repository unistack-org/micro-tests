package grpc

//go:generate protoc -I./internal/errors -I. --go-grpc_out=paths=source_relative:./internal/errors --go_out=paths=source_relative:./internal/errors --micro_out=paths=source_relative:./internal/errors internal/errors/errors.proto
