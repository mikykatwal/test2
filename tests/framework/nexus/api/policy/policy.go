package policy

//go:generate protoc --proto_path=../:../../../../../cnquery:. --go_out=. --go_opt=paths=source_relative --rangerrpc_out=. --iam-actions_out=. policy.proto