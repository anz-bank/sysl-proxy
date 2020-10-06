package main

import (
	sysl_proxy "github.com/joshcarp/sysl-proxy"
	"log"
	"net"
	"net/http"

)

func main() {
	http.HandleFunc("/", sysl_proxy.ServeHTTP)
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.Serve(lis, nil))
}
