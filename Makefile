.PHONY: container publish serve serve-container clean

app        := emojitrack-gostreamer
docker-tag := emojitracker/gostreamer

bin/$(app): *.go
	go build -o $@

container:
	docker build -t $(docker-tag) .

publish: container
	docker push $(docker-tag)

serve: bin/$(app)
	env PATH=$(PATH):./bin forego start web

serve-container:
	docker run -it --rm --env-file=.env -p 8001:8001 $(docker-tag)

clean:
	rm -rf bin
