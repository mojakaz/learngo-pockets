//go:generate bash -c "protoc -I=proto/ --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto"
package api
