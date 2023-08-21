cd services/feedService
protoc --micro_out=./ --go_out=./ feedService.proto
protoc-go-inject-tag -input=../feedService.pb.go
cd .. && cd ..