package main

import (
	"fmt"

	"net/http"
)

func main() {
	//http requests comes to the root will be handled by the indexHandler
	http.Handle("/", indexHandler{})
	//http requests comes to the /redis/* will be handled by the redisHandler from the redis.go
	http.Handle("/redis/", redisHandler{})
	//http requests comes to the /db/* will be handled by the s3Handler from the s3.go
	http.Handle("/s3/", s3Handler{})
	//http requests comes to the /db/* will be handled by the dbHandler from the db.go
	http.Handle("/db/", dbHandler{})

	//listen to the post :8080
	http.ListenAndServe(":8080", nil)

	fuga()
}

func fuga() {
	// go vet unreachable
	// return
	fmt.Println("fuga")
}

type indexHandler struct{}

func (h indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// out put a text to the browser.
	fmt.Fprintf(w, "hello 2, you've hit %s\n", r.URL.Path)
}

// function to demonstrate a simple test in main_test.go
func hello() string {
	return "Hello"
}
