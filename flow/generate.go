package flow

//go:generate protoc -I./proto -I. -I/home/vtolstov/devel/projects/unistack/micro/micro-proto --go_out=paths=source_relative:./proto proto/test.proto

//go:generate protoc -I./proto -I. -I/home/vtolstov/devel/projects/unistack/micro/micro-proto --go-micro_out=components=micro|http,debug=true,tag_path=./proto,paths=source_relative:./proto proto/test.proto
