FROM alpine:3.14
RUN apk add --no-cache ca-certificates
COPY shortener .
EXPOSE 8080
ENTRYPOINT [ "./shortener" ]