FROM golang:1.22

WORKDIR /demo
COPY main.go .
RUN go build main.go
CMD [ "/demo/main" ]