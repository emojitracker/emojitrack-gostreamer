package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

//TODO: maybe get rid of this and just use redis.Message if there is generic interface
type RedisMsg struct {
	channel string
	data    []byte
}

func RedisGo() (<-chan RedisMsg, <-chan RedisMsg) {

	/* connec to the redis server */
	server := os.Getenv("REDIS_URL")
	pass := os.Getenv("REDIS_PASS")
	c, err := redis.Dial("tcp", server)
	if err != nil {
		panic(err)
	}
	_, err2 := c.Do("AUTH", pass)
	if err2 != nil {
		panic(err)
	}

	/* set up structures and channels to stream events out on */
	scoreUpdates := make(chan RedisMsg)
	detailUpdates := make(chan RedisMsg)

	/* subscribe to and handle streams */
	psc := redis.PubSubConn{c}
	psc.Subscribe("stream.score_updates")
	psc.PSubscribe("stream.tweet_updates.*")

	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				//fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
				scoreUpdates <- RedisMsg{v.Channel, v.Data} //string(v.Data)
			case redis.PMessage:
				//fmt.Printf("pattern: %s, channel: %s, data: %s\n", v.Pattern, v.Channel, v.Data)
				detailUpdates <- RedisMsg{v.Channel, v.Data}
			case error:
				fmt.Println("redis errored?@&*(#)akjd")
				panic(v)
			}
		}
	}()

	return scoreUpdates, detailUpdates
}
