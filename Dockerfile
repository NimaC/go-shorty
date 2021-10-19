FROM golang:1.17.2-stretch

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /goshorty

CMD [ "/goshorty" ]