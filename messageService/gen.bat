cd services/messageService
protoc --micro_out=./ --go_out=./ messageService.proto
protoc-go-inject-tag -input=./messageService.pb.go
cd .. && cd ..