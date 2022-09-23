package clientserver

import (
	"fmt"
	"net/http"
	"sync"
)

func StartServer(port string, wg *sync.WaitGroup) {
	defer wg.Done()

	http.HandleFunc("/placements/request", handleRequest)

	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
