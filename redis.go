package main

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// it'd be nice to get rid of this and just use redis.Message generic interface
// ...but sigh Go: https://github.com/garyburd/redigo/issues/51
//
// let's just keep this and warp, even though there is some slight copy overhead
// probably, but it's better than having to have to have the rest of the code
// differentiate between structs based on where they originated.

// RedisMsg is our wrapper around redis.Message and redis.PMessage, since they
// really should be the same.
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

	go func() {
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
				log.Println("redis subscribe connection errored?@&*(#)akjd")
				// probable cause is connection was EOF
				// reminder: in this context, "Close" means just return to pool
				// pool will detect if connection is errored via testOnBorrow
				conn.Close()

				log.Println("attempting to get a new one in 5 seconds...")
				time.Sleep(5 * time.Second)
				conn = redisPool.Get()
				psc = redis.PubSubConn{conn}
				psc.Subscribe("stream.score_updates")
				psc.PSubscribe("stream.tweet_updates.*")
			}
		}
	}()

	return scoreUpdates, detailUpdates
}
