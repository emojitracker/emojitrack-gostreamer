.PHONY: container publish serve serve-container clean

app        := emojitrack-gostreamer
static-app := build/linux-amd64/$(app)
docker-tag := emojitracker/gostreamer

bin/$(app): *.go
	go build -o $@

$(static-app): *.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -ldflags "-s" -a -installsuffix cgo -o $(static-app)

container: $(static-app)
	docker build -t $(docker-tag) .

publish: container
	docker push $(docker-tag)

serve: bin/$(app)
	env PATH=$(PATH):./bin forego start web

serve-container:
	docker run -it --rm --env-file=.env -p 8001:8001 $(docker-tag)

clean:
	rm -rf bin build
