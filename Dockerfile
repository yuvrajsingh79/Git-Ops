from golang:latest
LABEL maintainer="Yuvraj Singh <singhyuvraj79@gmail.com>"

RUN mkdir $GOPATH/src/almabase
RUN go get github.com/gorilla/mux
RUN go get github.com/google/go-github/github
RUN go get golang.org/x/oauth2

ADD . $GOPATH/src/almabase/Git-Ops
WORKDIR $GOPATH/src/almabase/Git-Ops

RUN go build -o /app/main .

CMD ["/app/main"]
EXPOSE 8080