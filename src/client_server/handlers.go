package clientserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	customtypes "sample-choose-ad/src/custom_types"
	req_types "sample-choose-ad/src/requests_types"
)

// Create requset body based in incoming reqest `ir` and return
// `OutgoingRequest` as bytes.Reader from marshaled JSON
func constructPartnersRequestBody(ir *req_types.IncomingRequest) io.Reader {
	var outReqBody req_types.OutgoingRequest

	var imps []req_types.Imp

	// WARN: uint and float multiplication may cause problems
	for _, tile := range ir.Tiles {
		imps = append(imps, req_types.Imp{
			Id:        tile.Id,
			Minwidth:  tile.Width,
			Minheight: uint(math.Floor(float64(tile.Width * uint(tile.Ratio))))})
	}

	outReqBody.Id = *ir.Id
	outReqBody.Imp = imps
	outReqBody.Context = ir.Context
	t, _ := json.Marshal(outReqBody)
	return bytes.NewReader(t)
}

// Parsing and checking incoming request.
func parseAndCheckIncomingRequest(w http.ResponseWriter, r *http.Request) (req_types.IncomingRequest, error) {
	body, _ := ioutil.ReadAll(r.Body)

	var inpReqBody req_types.IncomingRequest
	var err error

	if json.Unmarshal(body, &inpReqBody) != nil {
		throwHTTPError("WRONG_SCHEMA", 400, &w)
		return inpReqBody, err
	}

	// Check if Id is empty
	if inpReqBody.Id == nil {
		throwHTTPError("EMPTY_FIELD", 400, &w)
		return inpReqBody, err
	}

	// Check if tiles is empty
	if len(inpReqBody.Tiles) == 0 {
		throwHTTPError("EMPTY_TILES", 400, &w)
		return inpReqBody, err
	}

	// ipv4 validation
	if wrongIPAddresFormat(inpReqBody.Context.Ip) {
		throwHTTPError("WRONG_SCHEMA", 400, &w)
		return inpReqBody, err
	}

	return inpReqBody, err
}

// Request handler with closure (make request for each partner in `[]partners`).
func handleRequest(partners []customtypes.PartnersAddress) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		incReq, err := parseAndCheckIncomingRequest(w, r)
		if err != nil {
			log.Println(err)
		}

		p_body := constructPartnersRequestBody(&incReq)

		var partnersRespones []req_types.SuccesResponse
		for _, p := range partners {
			url := fmt.Sprintf("http://%v:%v", p.Ip, p.Port)

			re, err := sendRequest(url, &p_body)

			if err != nil {
				log.Println(err)
				continue
			}
			// adding only successful responses
			partnersRespones = append(partnersRespones, re)
		}

		if len(partnersRespones) == 0 {
			log.Println("Error: no responses from partners.")
			return
		}

		// для каждого tile в incReq надо найти
		// maxPrice := 0

		// var bestOptionRespone req_types.SuccesResponse
		var bestOptions []req_types.Imp
		_ = bestOptions

		for _, tile := range incReq.Tiles {
			for _, partner := range partnersRespones {
			}
		}

	}
}
