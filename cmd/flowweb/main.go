package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Gixgine-budi/flow-free-bridge-solver/internal/web"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	log.Printf("Flow Free Bridges solver UI on http://localhost%s", *addr)
	if err := http.ListenAndServe(*addr, web.Handler()); err != nil {
		fmt.Println("error:", err)
	}
}
