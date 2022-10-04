package clientserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	req_types "sample-choose-ad/src/requests_types"
	"time"
)

func sendRequest(url string, body *io.Reader) (req_types.SuccesResponse, error) {
	var pResp req_types.SuccesResponse

	c := &http.Client{
		Timeout: 200 * time.Millisecond,
	}

	resp, err := c.Post(url, "application/json", *body)

	if err != nil {
		log.Println(err)
		eText := fmt.Sprintf("%v\n not responding", url)
		return pResp, errors.New(eText)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return pResp, errors.New("No content")
	}

	b, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(b, &pResp)
	if err != nil {
		log.Println(err)
	}
	return pResp, nil
}

/*
   key string

   map[price]{Imp}

*/

func sendRequest2(url string, body *io.Reader) ([]req_types.RespImp, error) {
	var pResp req_types.SuccesResponse

	c := &http.Client{
		Timeout: 200 * time.Millisecond,
	}

	resp, err := c.Post(url, "application/json", *body)

	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode == 204 {
		return pResp.Imp, errors.New("No content")
	}

	b, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(b, &pResp)
	if err != nil {
		log.Println(err)
	}

	return pResp.Imp, nil
}
