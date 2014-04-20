package main

import (
  "fmt"
  "time"
  "math/rand"
  "os"
  "github.com/garyburd/redigo/redis"
  _"github.com/joho/godotenv/autoload"
)

/* fake scoreUpdate stream just good enough to write code without wasting network connection */
func fakeRedisStream() <-chan string {
  c := make(chan string)
  keys := []string{"AAAA", "BBBB", "CCCC", "DDDD"}
  go func() {
    for {
      c <- keys[ rand.Intn(3) ]
      time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
    }
  }()
  return c
}


type RedisMsg struct {
  channel string
  data []byte
}

func RedisGo() (<-chan string, <-chan RedisMsg) {

  /* connec to the redis server */
  server := os.Getenv("REDIS_URL")
  pass   := os.Getenv("REDIS_PASS")
  c, err := redis.Dial("tcp", server)
  if err != nil {
    panic(err)
  }
  _, err2 := c.Do("AUTH", pass);
  if err2 != nil {
    panic(err)
  }

  /* set up structures and channels to stream events out on */
  scoreUpdates := make(chan string)
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
            scoreUpdates <- string(v.Data)
        case redis.PMessage:
            //fmt.Printf("pattern: %s, channel: %s, data: %s\n", v.Pattern, v.Channel, v.Data)
            //detailUpdates <- RedisMsg{v.Channel, v.Data}
        case error:
            fmt.Println("error")
      }
    }
  }()

  return scoreUpdates, detailUpdates
}
