package clientserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	customtypes "sample-choose-ad/cmd/custom_types"
	req_types "sample-choose-ad/cmd/requests_types"
	"strconv"
	"sync"
)

const PARTNER_ENDPOINT = "bid_request"

// Parsing and checking incoming request.
func parseAndCheckIncomingRequest(w http.ResponseWriter, r *http.Request) (req_types.IncomingRequest, error) {

	var inpReqBody req_types.IncomingRequest
	var err error

	//check request method. Only POST valid.
	if r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return inpReqBody, errors.New("Wrong request method")
	}

	// Check if body in incoming request is empty
	body, _ := ioutil.ReadAll(r.Body)

	if json.Unmarshal(body, &inpReqBody) != nil {
		log.Println("Unmarshaling problem", string(body))
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
	// UserAgent validation
	if len(inpReqBody.Context.UserAgent) == 0 {
		return inpReqBody, throwHTTPError("EMPTY_FIELD", 400, &w)
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

		wg := new(sync.WaitGroup)
		responsesCh := make(chan []req_types.RespImp, len(partners))
		var respImps []req_types.RespImp
		// Send requests to partners, collecting responses in `responses` channel
		for _, p := range partners {
			wg.Add(1)
			url := fmt.Sprintf("http://%v:%v/%v", p.Ip, p.Port, PARTNER_ENDPOINT)
			go makeRequest(url, &p_body, responsesCh, wg)
		}
		wg.Wait()
		close(responsesCh)

		for r := range responsesCh {
			respImps = append(respImps, r...)
		}
		//We have no identical pairs `id` and `price`
		partnersRespones := make(map[uint]req_types.RespImp)
		for _, resp := range respImps {
			if _, exist := partnersRespones[resp.Id]; !exist {
				partnersRespones[resp.Id] = resp
				continue
			}

			// Replase with new Imp, if last saved price smaller
			// Using type switch and addition checks, 'cause `Price` is interface
			var oldPrice float64
			switch partnersRespones[resp.Id].Price.(type) {
			case string:
				oldPrice, _ = strconv.ParseFloat(partnersRespones[resp.Id].Price.(string), 64)
			case float64:
				oldPrice = partnersRespones[resp.Id].Price.(float64)
			}

			var newPrice float64
			switch resp.Price.(type) {
			case string:
				newPrice, _ = strconv.ParseFloat(resp.Price.(string), 64)
			case float64:
				newPrice = resp.Price.(float64)
			}

			if oldPrice < newPrice {
				partnersRespones[resp.Id] = resp
			}
		}

		var bestOptions []req_types.RespImp

		// tile.Id == RespImp.Id
		// for each tile peak best price
		for _, tile := range incReq.Tiles {
			if val, exist := partnersRespones[tile.Id]; exist {
				bestOptions = append(bestOptions, val)
			}

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
