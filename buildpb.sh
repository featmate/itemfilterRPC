protoc -I /Users/mac/WORKSPACE/GITHUB/featmate/itemfilterRPC  --go_out=. --go-grpc_out=. --grpc-gateway_out=.  --openapiv2_out=. pbschema/itemfilterRPC/*.proto