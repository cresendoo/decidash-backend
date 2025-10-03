# decidash-backend


# deps
```bash
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.73.0
```

```bash
protoc \
--proto_path=proto \
--go_out=pkg \
--go_opt=paths=source_relative \
--go-grpc_out=pkg \
--go-grpc_opt=paths=source_relative \
    proto/aptos/indexer/v1/grpc.proto \
    proto/aptos/indexer/v1/raw_data.proto \
    proto/aptos/indexer/v1/filter.proto \
    proto/aptos/transaction/v1/transaction.proto \
    proto/aptos/util/timestamp/timestamp.proto \
    proto/aptos/remote_executor/v1/network_msg.proto \
    proto/aptos/internal/fullnode/v1/fullnode_data.proto
```