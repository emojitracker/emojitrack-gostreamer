package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

// TODO: maybe get rid of this and just use redis.Message if there is generic interface
// sigh golang...: https://github.com/garyburd/redigo/issues/51
type RedisMsg struct {
	channel string
	data    []byte
}

// direct copy from http://godoc.org/github.com/garyburd/redigo/redis#Pool
// why do we need to cut and paste code instead of having it be built-in
// to the package?  because golang!
func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// this should make a global var we can access wherever?
var (
	redisPool *redis.Pool
)

func initRedisPool() {
	host, pass := envRedis()
	redisPool = newPool(host, pass)
}

func myRedisSubscriptions() (<-chan RedisMsg, <-chan RedisMsg) {

	// set up structures and channels to stream events out on
	scoreUpdates := make(chan RedisMsg)
	detailUpdates := make(chan RedisMsg)

	// get a new redis connection from pool.
	// since this is the first time the app tries to do something with redis,
	// die if we can't get a valid connection, since something is probably
	// configured wrong.
	conn := redisPool.Get()
	_, err := conn.Do("PING")
	if err != nil {
		log.Fatal("Could not connect to Redis, check your configuration.")
	}

	// subscribe to and handle streams
	psc := redis.PubSubConn{conn}
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
				//TODO: at some point we might need to also match the pattern here for kiosk mode
				detailUpdates <- RedisMsg{v.Channel, v.Data}
			case error:
				fmt.Println("redis subscribe connection errored?@&*(#)akjd")
				// probable cause is connection was closed, but force close just in case
				conn.Close()

				fmt.Println("attempting to get a new one in 5 seconds...")
				time.Sleep(5 * time.Second)
				conn = redisPool.Get()
			}
		}
	}()

	return scoreUpdates, detailUpdates
}
