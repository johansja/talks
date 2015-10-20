//go:generate protoc -I. -I/usr/local/include -I$GOPATH/src/github.com/gengo/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. player2.proto
//go:generate protoc -I. -I/usr/local/include -I$GOPATH/src/github.com/gengo/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. player2.proto
package player
