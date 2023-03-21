package main

import (
	polygon "github.com/polygon-io/client-go/rest"
)

func main() {
	polygonApi := polygon.New("FVcbDR6ZtUfTl2URJWZfPVRFNkL2kvnJ")
	storage, err := NewPostgresStore()

	if err != nil {
		return
	}
	server := NewApiServer(":8080", storage)
	server.Run()
}
