FROM golang:1.12

LABEL maintainer="rictorlome@gmail.com"

WORKDIR $GOPATH/src/github.com/rictorlome/pupil

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 8080

CMD ["pupil"]