package main

import (
	"os"

	serv "github.com/featmate/itemfilterRPC/itemfilterRPC_serv"

	log "github.com/Golang-Tools/loggerhelper/v2"
	s "github.com/Golang-Tools/schema-entry-go/v2"
)

func main() {
	serv, err := s.NewEndPoint(&serv.Server{}, s.WithName("itemfilterRPC"), s.WithUsage("itemfilterRPC [options]"))
	if err != nil {
		log.Error("create serv node get error", log.Dict{"err": err.Error()})
		os.Exit(2)
	}

	serv.Parse(os.Args)
}
