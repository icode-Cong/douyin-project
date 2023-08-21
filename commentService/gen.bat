cd services/commentService
protoc --micro_out=./ --go_out=./ commentService.proto
protoc-go-inject-tag -input=./commentService.pb.go
cd .. && cd ..
