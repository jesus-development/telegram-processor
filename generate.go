package main

// in terminal: protoc --proto_path . --go_out . --go_opt paths=source_relative --go-grpc_out . --go-grpc_opt paths=source_relative --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out . --openapiv2_opt use_go_templates=true ./pkg/api/proto/*.proto
//go:generate find pkg/api/proto/. -type f -name "*.proto" -exec protoc -I=. --go_out . --go_opt paths=source_relative --go-grpc_out . --go-grpc_opt paths=source_relative --grpc-gateway_out . --grpc-gateway_opt paths=source_relative --openapiv2_out . --openapiv2_opt use_go_templates=true {} ;
//go:generate find pkg/models/proto/. -type f -name "*.proto" -exec protoc -I=. --go_out . --go_opt paths=source_relative {} ;
//go:generate find pkg/models/json/. -type f -name "*.go" -not -name "*easyjson.go" -exec easyjson -all {} ;
//go:generate mockgen -source=internal/services/processor/processor.go -destination=internal/services/processor/mock/processor.go -package=mock
