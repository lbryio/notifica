package main

import "github.com/lbryio/notifica/cmd"

//go:generate protoc --proto_path=app --go_out=app/ --go_opt=paths=source_relative types/notification.proto

func main() {
	cmd.Execute()
}
