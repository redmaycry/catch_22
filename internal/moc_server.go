package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	vres := `{
    "id": "123",
    "imp": [{
        "id": 123,
        "width": 144,
        "height": 122,
        "title": "Title1",
        "url": "example.com",
        "price": 123.5
    },{
        "id": 123,
        "width": 155,
        "height": 133,
        "title": "Title2",
        "url": "upachka.com",
        "price": 143.5
    }]}`

	addr := flag.String("l", "", "-l 127.0.0.1:5059")
	flag.Parse()
	if *addr == "" {
		log.Fatalln("Error: listening address is required!")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		log.Println(string(b))

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(vres))
	})
	log.Fatal(http.ListenAndServe(*addr, nil))
}
