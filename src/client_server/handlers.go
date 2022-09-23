package clientserver

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	req_types "sample-choose-ad/src/requests_types"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var inpReqBody req_types.IncomingRequest

	err := json.Unmarshal(body, &inpReqBody)

	if err != nil {
		throwHTTPError("WRONG_SCHEMA", 400, &w)
		return
	}

	// Check if Id is empty
	if inpReqBody.Id == nil {
		throwHTTPError("EMPTY_FIELD", 400, &w)
		return
	}

	// Check if tiles is empty
	if len(inpReqBody.Tiles) == 0 {
		throwHTTPError("EMPTY_TILES", 400, &w)
		return
	}

	// ipv4 validation
	if wrongIPAddresFormat(inpReqBody.Context.Ip) {
		throwHTTPError("WRONG_SCHEMA", 400, &w)
		return
	}

	// TODO: fix url
	// TODO: может что то получится сделать с p_body?
	p_body := constructPartnersRequestBody(&inpReqBody)
	sendRequest("localhost:5059", &p_body)

}

func constructPartnersRequestBody(ir *req_types.IncomingRequest) io.Reader {
	var outReqBody req_types.OutgoingRequest

	var imps []req_types.Imp

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
