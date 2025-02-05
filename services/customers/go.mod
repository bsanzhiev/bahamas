module github.com/bsanzhiev/bahamas/services/customers

go 1.23.4

require (
	github.com/bsanzhiev/bahamas v0.0.0
	github.com/segmentio/kafka-go v0.4.47
	google.golang.org/grpc v1.70.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/protobuf v1.36.4 // indirect
)

// For local development use the following replace directives
replace github.com/bsanzhiev/bahamas => ../..
replace github.com/bsanzhiev/bahamas/libs => ../../libs
