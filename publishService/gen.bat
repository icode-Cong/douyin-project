@REM cd protos/publishService
@REM protoc --micro_out=./ --go_out=./ publishService.proto
@REM protoc-go-inject-tag -input=./publishService.pb.go
@REM cd .. && cd ..

cd protos/userPublish
protoc --micro_out=./ --go_out=./ userPublish.proto
protoc-go-inject-tag -input=./userPublish.pb.go
cd .. && cd ..