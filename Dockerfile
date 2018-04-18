#FROM scratch
FROM alpine

WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY migrations ./migrations
COPY main ./main

ENTRYPOINT ["./main"]
