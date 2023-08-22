cd services/relationService
protoc --micro_out=./ --go_out=./ relationService.proto
protoc-go-inject-tag -input=./relationService.pb.go
cd .. && cd ..