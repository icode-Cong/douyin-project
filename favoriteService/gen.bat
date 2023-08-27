cd services/favoriteService
protoc --micro_out=./ --go_out=./ favoriteService.proto
protoc-go-inject-tag --input=./favoriteService.pb.go
cd .. && cd ..