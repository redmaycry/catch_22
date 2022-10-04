package clientserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	customtypes "sample-choose-ad/src/custom_types"
	req_types "sample-choose-ad/src/requests_types"
	"sort"
)

const PARTNER_ENDPOINT = "bid_request"

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

	var inpReqBody req_types.IncomingRequest
	var err error

	//check request method. Only POST valid.
	if r.Method == "GET" {
		return inpReqBody, errors.New("Wrong request method")
	}

	// Check if body in incoming request is empty
	body, _ := ioutil.ReadAll(r.Body)

	if json.Unmarshal(body, &inpReqBody) != nil {
		return inpReqBody, throwHTTPError("WRONG_SCHEMA", 400, &w)
	}

	// Check if Id is empty
	if inpReqBody.Id == nil {
		return inpReqBody, throwHTTPError("EMPTY_FIELD", 400, &w)
	}

	// Check if tiles is empty
	if len(inpReqBody.Tiles) == 0 {
		return inpReqBody, throwHTTPError("EMPTY_TILES", 400, &w)
	}

	// ipv4 validation
	if wrongIPAddresFormat(inpReqBody.Context.Ip) {
		return inpReqBody, throwHTTPError("WRONG_SCHEMA", 400, &w)
	}

	return inpReqBody, err
}

// Request handler with wrapper (make request for each partner in `[]partners`).
func handleRequest(partners []customtypes.PartnersAddress) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Parse incoming request and return an error, if it's empty
		// or contains wrong/empty fields
		incReq, err := parseAndCheckIncomingRequest(w, r)
		if err != nil {
			log.Println(err)
			return
		}

		p_body := constructPartnersRequestBody(&incReq)

		// Two data structures:
		// partnersRespones for getting price with O(1) complexity
		// []prices as slice of actual prices
		// var partnersRespones map[float64]req_types.RespImp
		partnersRespones := make(map[uint]map[float64]req_types.RespImp)
		prices := make(map[uint][]float64)

		for _, p := range partners {
			url := fmt.Sprintf("http://%v:%v/%v", p.Ip, p.Port, PARTNER_ENDPOINT)

			re, err := sendRequest(url, &p_body)

			if err != nil {
				log.Println(err)
				continue
			}
			// adding only successful responses
			for _, r := range re.Imp {
				if partnersRespones[r.Id] == nil {
					partnersRespones[r.Id] = make(map[float64]req_types.RespImp)
				}
				partnersRespones[r.Id][r.Price] = r
				prices[r.Id] = append(prices[r.Id], r.Price)
			}

		}

		if len(partnersRespones) == 0 {
			log.Println("Error: no responses from partners.")
			return
		}

		// Sorting prices, now biggest price at index len-1
		for _, p := range prices {
			sort.Float64s(p)
		}

		var bestOptions []req_types.RespImp

		// for each tile peak best price
		for _, tile := range incReq.Tiles {
			last := len(prices[tile.Id]) - 1
			biggestPrice := prices[tile.Id][last]
			_ = biggestPrice
			bestOptions = append(bestOptions, partnersRespones[tile.Id][biggestPrice])
		}

		response := req_types.SuccesResponse{
			Id:  *incReq.Id,
			Imp: bestOptions,
		}

		respJSON, err := json.Marshal(response)

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respJSON)
	}
}
