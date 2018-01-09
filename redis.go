package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/garyburd/redigo/redis"
)

type redisHandler struct{}

func (h redisHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//get the current views count from redis storage
	v := getViews()
	v++
	//set the new views count to redis storage
	updateViews(v)
	//output the view count to browser
	fmt.Fprintf(w, "Viewed count from redis : "+strconv.Itoa(v))
}

//connects to redis service
func redisConnect() redis.Conn {
	c, err := redis.Dial("tcp", "redis:6379")
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func updateViews(cnt int) {

	c := redisConnect()
	defer c.Close()

	// set the value on redis for the key viewedcount
	reply, err := c.Do("SET", "viewedcount", cnt)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("GET ", reply)
}

func getViews() int {

	c := redisConnect()
	defer c.Close()
	// get the value from redis for the key viewedcount
	reply, err := c.Do("GET", "viewedcount")
	if err != nil {
		log.Fatal(err)
	}
	if reply != nil {
		s := string(reply.([]byte))
		log.Println("GET ", s)
		i, _ := strconv.Atoi(s)
		return i
	}

	return 0
}
