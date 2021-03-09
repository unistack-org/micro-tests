package http

//go:generate protoc -I./proto -I. -I/home/vtolstov/.cache/go-path/pkg/mod/github.com/unistack-org/micro-proto@v0.0.1 --go_out=paths=source_relative:./proto --micro_out=components=micro|http,debug=true,paths=source_relative:./proto proto/test.proto

///go:generate protoc -I./proto -I. -I/home/vtolstov/devel/projects/unistack/micro/micro-proto --go_out=paths=source_relative:./proto --micro_out=components=micro|http,debug=true,paths=source_relative:./proto proto/test.proto
