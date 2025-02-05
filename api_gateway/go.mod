module github.com/bsanzhiev/bahamas/api_gateway

go 1.23.4

require (
	github.com/bsanzhiev/bahamas/services/customers v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.70.0
)

require (
	github.com/bsanzhiev/bahamas v0.0.0 // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/protobuf v1.36.4 // indirect
)

replace (
	github.com/bsanzhiev/bahamas => ../
	github.com/bsanzhiev/bahamas/libs/pb/customers => ../libs/pb/customers
	github.com/bsanzhiev/bahamas/services/customers => ../services/customers
)
