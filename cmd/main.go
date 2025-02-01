package main

import (
	"flag"

	"github.com/felipefrmelo/crm_backend_udacity"
)

func main() {
	server := flag.String("server", "fiber", "server to use (fiber, gorilla)")

  flag.Parse()

	repo := crm.NewRepo()
	app := crm.NewServer(repo, *server)

	app.Listen(":3000")

}
