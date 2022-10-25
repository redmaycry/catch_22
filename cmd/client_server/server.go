package clientserver

import (
	"fmt"
	"log"
	"net/http"
	customtypes "sample-choose-ad/cmd/custom_types"
	"time"
)

const MAX_TIME_PER_REQUEST = time.Duration(250 * time.Millisecond)

func StartServer(port string, partners []customtypes.PartnersAddress) {
	h := http.TimeoutHandler(handleRequest(partners), MAX_TIME_PER_REQUEST, "")
	http.Handle("/placements/request", h)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
