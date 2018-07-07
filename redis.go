package main

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

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

func myRedisSubscriptions() (<-chan redis.Message, <-chan redis.Message) {
	// set up structures and channels to stream events out on
	scoreUpdates := make(chan redis.Message)
	detailUpdates := make(chan redis.Message)

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
		const scoreKey = "stream.score_updates"
		const detailKey = "stream.tweet_updates.*"
		psc := redis.PubSubConn{Conn: conn}
		psc.Subscribe(scoreKey)
		psc.PSubscribe(detailKey)

		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				switch {
				case v.Channel == scoreKey:
					scoreUpdates <- v
				case v.Pattern == detailKey:
					detailUpdates <- v
				default:
					log.Println("Received a message on an unexpected channel ", v.Channel)
				}
			case error:
				log.Println("redis subscribe connection errored?@&*(#)akjd")
				// probable cause is connection was EOF
				// reminder: in this context, "Close" means just return to pool
				// pool will detect if connection is errored via testOnBorrow
				conn.Close()

				log.Println("attempting to get a new one in 5 seconds...")
				time.Sleep(5 * time.Second)
				conn = redisPool.Get()
				psc = redis.PubSubConn{Conn: conn}
				psc.Subscribe("stream.score_updates")
				psc.PSubscribe("stream.tweet_updates.*")
			}
		}
	}()

	return scoreUpdates, detailUpdates
}
