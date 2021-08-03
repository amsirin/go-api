FROM scratch

ENV PORT 8000
EXPOSE $PORT

COPY go-api /
CMD ["/go-api"]
