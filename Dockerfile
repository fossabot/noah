FROM golang:latest as noah-build
RUN echo $GOPATH
RUN mkdir -p $GOPATH/src/github.com/readystock/noah
COPY ./ $GOPATH/src/github.com/readystock/noah
WORKDIR $GOPATH/src/github.com/readystock/noah
RUN go get -u
RUN go build -o noah

FROM noah-build
RUN mkdir /ndb
COPY --from=noah-build $GOPATH/src/github.com/readystock/noah /ndb/noah
CMD ["/ndb/noah"]
