FROM golang:alpine

RUN apk update && apk add git

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /url-shortener

EXPOSE 8080

CMD [ "/url-shortener" ]
