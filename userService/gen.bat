cd services/protos
protoc --micro_out=./ --go_out=./ user.proto
@REM protoc --micro_out=../tokenproto --go_out=../tokenproto tokenproto/token_utils.proto
@REM protoc --micro_out=./ --go_out=./ from_relation/to_user.proto
protoc-go-inject-tag -input=../user.pb.go
cd .. && cd ..