FROM alpine
RUN apk add --no-cache ca-certificates
COPY apiserver /apiserver
EXPOSE 80
ENTRYPOINT ["/apiserver"]