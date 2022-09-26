package clientserver

import (
	"fmt"
	"net/http"
)

func StartServer(port string) {

	http.HandleFunc("/placements/request", handleRequest)

	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
