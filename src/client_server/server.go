package clientserver

import (
	"fmt"
	"net/http"
	customtypes "sample-choose-ad/src/custom_types"
)

func StartServer(port string, partners []customtypes.PartnersAddress) {

	http.HandleFunc("/placements/request", handleRequest(partners))
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)

}
