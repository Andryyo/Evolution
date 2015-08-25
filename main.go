// Evolution project main.go
package main

import (
	"github.com/Andryyo/Evolution/EvolutionServer"
	"log"
	"net/http"
)

func main() {
	server := EvolutionServer.NewServer()
	go server.Listen()
	http.Handle("/", http.FileServer(http.Dir("EvolutionServer/webroot")))
	//log.Fatal(http.ListenAndServe(":8081", nil))
	log.Fatal(http.ListenAndServe(":80", nil))
	//NewGame(&ConsoleChoiceMaker{"One"}, &ConsoleChoiceMaker{"Two"})
}
