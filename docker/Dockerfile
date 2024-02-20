FROM golang:1.21

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . ./

RUN go build -o ./app

EXPOSE 8001
EXPOSE 8002

CMD ["/src/app"]