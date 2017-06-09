FROM golang:1.8.3-alpine

RUN apk add --no-cache git

RUN mkdir -p /go/src/image-worker
COPY . /go
WORKDIR /go/src/image-worker

RUN ls

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]