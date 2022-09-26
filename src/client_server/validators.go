package clientserver

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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

func ParsePartnersAddress(ipAndPort string) (string, int64, error) {
	var err error
	iap := strings.Split(ipAndPort, ":")

	ip := iap[0]
	if wrongIPAddresFormat(ip) {
		err = errors.New(fmt.Sprintf("Wrong ip address format in partner ip: %v", ip))
	}

	port, _ := strconv.ParseInt(iap[1], 10, 32)
	return ip, port, err
}
