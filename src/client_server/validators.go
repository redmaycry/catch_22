package clientserver

import (
	"log"
	"net/http"
	"regexp"
)

// Returns false if ipv4 `correct`.
func wrongIPAddresFormat(ipv4 string) bool {
	re, err := regexp.Compile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	if err != nil {
		log.Println(err)
	}
	return !re.Match([]byte(ipv4))
}

func throwHTTPError(err_text string, code int, w *http.ResponseWriter) {
	http.Error(*w, err_text, code)
	log.Printf("Error %d %v\n", code, err_text)
}
