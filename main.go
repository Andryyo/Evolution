// Evolution project main.go
package main

import (
	"net/http"
	"log"
	"github.com/Andryyo/Evolution/EvolutionServer"
)

func main() {
	server := EvolutionServer.NewServer()
	go server.Listen()
	http.Handle("/", http.FileServer(http.Dir("EvolutionServer/webroot")))
	log.Fatal(http.ListenAndServe(":8080", nil))
	//NewGame(&ConsoleChoiceMaker{"One"}, &ConsoleChoiceMaker{"Two"})
}
