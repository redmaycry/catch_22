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

	// TODO: fix url TODO: может что то получится сделать с p_body?
	// Такие запросы надо разослать на список партнеров (-d). Из их
	// ответов потом нужно выбрат
	/*
		    Получив все ответы, ваш сервис должен для каждого элемента tiles из запроса
			плейсмента выбрать среди imp с таким же id тот, у которого максимальная цена.
			Одинаковых цен с одним id не будет
			Если imp с каким-то id не получено, такого id не должно быть в ответе. Порядок imp в
			ответе должен соответствовать порядку tiles в запросе плейсмента. Формат ответа:
	*/
	// Отпралять запросы как горутины!
	p_body := constructPartnersRequestBody(&inpReqBody)
	sendRequest("localhost:5059", &p_body)

}

func constructPartnersRequestBody(ir *req_types.IncomingRequest) io.Reader {
	var outReqBody req_types.OutgoingRequest

	var imps []req_types.Imp

	// WARN: не знаю как правильно перемножать uint и float
	for _, tile := range ir.Tiles {
		imps = append(imps, req_types.Imp{
			Id:        tile.Id,
			Minwidth:  tile.Width,
			Minheight: uint(math.Floor(float64(tile.Width * uint(tile.Ratio))))})
	}

	outReqBody.Id = *ir.Id
	outReqBody.Imp = imps
	outReqBody.Context = ir.Context
	log.Println(*ir)
	log.Println(outReqBody)
	t, _ := json.Marshal(outReqBody)
	return bytes.NewReader(t)
}
