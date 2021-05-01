FROM golang:1.16-alpine AS builder

RUN mkdir -p /go/src/github.com/emojitracker/emojitrack-gostreamer
WORKDIR /go/src/github.com/emojitracker/emojitrack-gostreamer

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -ldflags "-s" -a -installsuffix cgo -o emojitrack-gostreamer

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/emojitracker/emojitrack-gostreamer/emojitrack-gostreamer .
ENTRYPOINT ["./emojitrack-gostreamer"]
