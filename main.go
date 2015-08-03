// Evolution project main.go
package main

import (
	"net/http"
	"log"
)

func main() {
	server := NewServer()
	go server.Listen()
	http.Handle("/", http.FileServer(http.Dir("webroot")))
	log.Fatal(http.ListenAndServe(":80", nil))
	//NewGame(&ConsoleChoiceMaker{"One"}, &ConsoleChoiceMaker{"Two"})
}
