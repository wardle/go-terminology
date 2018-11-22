FROM alpine:3.4

RUN apk -U add ca-certificates

EXPOSE 8080 8081

ADD gts-linux-amd64 /

ENTRYPOINT ["./gts-linux-amd64"]
CMD ["-db", "${DATABASE}", "-server"]