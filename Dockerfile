FROM scratch
MAINTAINER Matthew Rothenberg <mroth@mroth.info>

COPY build/linux-amd64/emojitrack-gostreamer /emojitrack-gostreamer
ENTRYPOINT ["/emojitrack-gostreamer"]
