package clientserver

import (
	"io"
	"log"
	"net/http"
	"time"
)

func sendRequest(url string, body *io.Reader) {
	c := &http.Client{
		Timeout: 200 * time.Millisecond,
	}

	resp, err := c.Post(url, "application/json", *body)

	if err != nil {
		log.Println(err)
	}
	log.Println(resp)
}
