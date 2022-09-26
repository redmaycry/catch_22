/*
Usage:

   sample-choose-ad [flags]

The flags are:
   -p
       Listening port
   -d
       Adversment partners list in format ip_p1:port,ip_p2:port2...ip_p10:port



*/
package main

import (
	"flag"
	"log"
	clientserver "sample-choose-ad/src/client_server"
	customtypes "sample-choose-ad/src/custom_types"
	"strings"
)

var partners []customtypes.PartnersAddress

func main() {
	log.Println("Starting server")

	port := flag.String("p", "", "-p 5050")
	addressesList := flag.String("d", "", "-d '10.10.10.10:5050,10.10.10.20:5050'")

	flag.Parse()

	if *port == "" {
		log.Fatalln("Error: Port number is require!")
	}

	if *addressesList == "" {
		log.Fatalln("Error: Partners list is require!")
	}

	// Parse first 10 ip:port pairs into `[]partners` slise

	for i, p := range strings.Split(*addressesList, ",") {

		if i == 10 {
			log.Println("Warning: Partners count must be less or equal 10!")
			return
		}

		ip, port, err := clientserver.ParsePartnersAddress(p)

		if err != nil {
			log.Println(err)
			continue
		}

		partners = append(partners, customtypes.PartnersAddress{
			Ip:   ip,
			Port: port})
	}

	clientserver.StartServer(*port, partners)
}
