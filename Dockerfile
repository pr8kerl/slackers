FROM golang:1.10 as build

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && chmod +x /usr/local/bin/dep

WORKDIR /go/src/github.com/pr8kerl/slackers
COPY . .
RUN make clean && make

FROM scratch
COPY --from=build /go/src/github.com/pr8kerl/slackers/slackers /
COPY ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/slackers"]
CMD [ "--help" ]
