package main

import (
	"fmt"
	"net/http"
	"time"
)

func main()  {
	go func() {
		for {

			fmt.Println("yay")
			time.Sleep(1 * time.Second)
		}
	}()
	startHttpServer()
}

func startHttpServer(){
	http.Handle("/", http.FileServer(http.Dir("public")))

	http.ListenAndServe(":3000", nil)

}