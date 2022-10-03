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

const MAX_PORT_NUM = 65535

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
	log.Printf("Error: %d %v\n", code, err_text)
}

// Wait string in format "10.10.10.10:8080", where `10.10.10.10` IPv4,
// and `8080` port. If ip or port has wrong format, returns error.
func ParsePartnersAddress(ipAndPort string) (string, int64, error) {
	var err error
	iap := strings.Split(ipAndPort, ":")

	ip := iap[0]
	if wrongIPAddresFormat(ip) {
		err = errors.New(fmt.Sprintf("Wrong ip address format in partner ip: %v", ip))
	}

	port, e := strconv.ParseInt(iap[1], 10, 32)
	if e != nil {
		err = errors.New(fmt.Sprintf("Wrong port format in partner ip: %v", e))
	}

	if port > MAX_PORT_NUM {
		err = errors.New(fmt.Sprintf("Wrong port in partner ip: grater than %v", MAX_PORT_NUM))
	}

	return ip, port, err
}
