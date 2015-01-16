emojitrack-gostreamer
=====================

Implementation of an API compatible server for
the [emojitrack-streamer-spec][1] in Go.

Caveats
-------
![noidea](http://static.someecards.com/someecards/usercards/MjAxMi0wZmYwN2ZiYmIyNzVkNzIy.png)

This was my first major Go project, so there may potentially still be some dumb stuff left 
in the code here.Please feel encouraged to send pull requests or line comments even for 
idiomatic syntax issues when accompanied with an explanation that will help me learn.

Pipeline Flow
-------------
Goroutines and channels make it just too much fun to think about things as flow
diagrams.  This is the essential pipeline for information flow in the `main()`
function for the program.

![diagram](http://f.cl.ly/items/1S282j3A0H07382x3q2L/diagram.svg)

Package sseserver
-----------------
This program doesn't actually do too much beyond the pipeline flow above, as all
of the SSE/HTTP server logic has been extracted into a handy modular package
that should be usable for any generic SSE broadcasting.  For that stuff, delve into the
(freely licensed!) [sseserver](https://github.com/mroth/sseserver/) repo.

Is it any fast?
-----------
Yes.:racehorse:

TODO
----

 - [x] move scorepacker and connectionpool into their own packages
 - [ ] possibly add tests for them even!
 - [x] handle redis server reconnects
 - [x] parse standard single `REDIS_URL` env var
 - [x] dont emit empty msgs (but lets wait until done benchmarking)

[1]: http://github.com/mroth/emojitrack-streamer-spec
