.PHONY: container publish serve serve-container clean

app        := emojitrack-gostreamer
linux-app  := build/linux-amd64/$(app)
docker-tag := emojitracker/gostreamer

$(app): *.go
	go build

$(linux-app):
	env GOOS=linux GOARCH=amd64 go build -o $(linux-app)

container: $(linux-app)
	docker build -t $(docker-tag) .

publish: container
	docker push $(docker-tag)

serve: $(app)
	env PATH=$(PATH):. forego start web

serve-container:
	docker run -it --rm --env-file=.env -p 8001:8001 $(docker-tag)

clean:
	go clean
	rm -rf build
