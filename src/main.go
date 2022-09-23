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
	"sync"
	// req_types "sample-choose-ad/src/requests_types"
)

// ============

func main() {
	var wg sync.WaitGroup
	log.Println("Starting server")

	/*
		file := flag.String("f", "", "path to file")
			port := flag.String("p", "5050", "listening port")

			flag.Parse()

			if *file == "" {
				fmt.Println("Please specify the path to the file!")
				return
			}
	*/
	port := flag.String("p", "", "-p 5050")

	flag.Parse()

	if *port == "" {
		log.Println("Port number is require!")
		return
	}

	wg.Add(1)
	go clientserver.StartServer(*port, &wg)
	wg.Wait()
}
