FROM golang:onbuild
FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o noah .
CMD ["/app/noah"]