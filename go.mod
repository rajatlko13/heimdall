module github.com/maticnetwork/heimdall

go 1.15

require (
	github.com/cbergoon/merkletree v0.2.0
	github.com/cosmos/cosmos-sdk v0.40.0-rc0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.15.2
	github.com/maticnetwork/bor v0.2.1
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/regen-network/cosmos-proto v0.3.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/tendermint/tendermint v0.34.0-rc5
	github.com/tendermint/tm-db v0.6.2
	github.com/ugorji/go v1.1.4 // indirect
	google.golang.org/grpc v1.33.1
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

