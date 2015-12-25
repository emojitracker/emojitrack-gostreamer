FROM scratch
COPY build/linux-amd64/emojitrack-gostreamer /emojitrack-gostreamer
EXPOSE 80
ENTRYPOINT ["/emojitrack-gostreamer"]
