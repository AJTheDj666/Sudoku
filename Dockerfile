FROM golang:1.17.2-alpine3.14
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod download
RUN go build -o main .
CMD ["/app/main"]