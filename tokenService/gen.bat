cd protos/token
protoc --micro_out=./ --go_out=./ token.proto
protoc-go-inject-tag -input=./token.pb.go
cd .. && cd ..