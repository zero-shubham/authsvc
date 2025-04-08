FROM golang:bookworm

WORKDIR /authsvc

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

CMD [ "go", "run", "./cmd/server/main.go" ]