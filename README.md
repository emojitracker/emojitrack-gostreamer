Attempting to implement an API compatible server from emojitrack-streamer-spec in Go.

My first Go project, so there will probably be some dumb stuff here.

Cribbing a lot for the HTTP stuff from:
http://gary.burd.info/go-websocket-chat

TODO:
 - move scorepacker and connectionpool into their own packages
 - possibly add tests for them even!
 - handle redis server reconnects
 - parse standard single `REDIS_URL` env var
 - dont emit empty msgs (but lets wait until done benchmarking)
