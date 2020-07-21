FROM golang:1.14

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ENV APP_USER app
ENV APP_HOME /go/src/voluum-apis

# setting working directory
WORKDIR /go/src/app

COPY / /go/src/app/

# installing dependencies
RUN go get -u golang.org/x/oauth2
RUN go get -u golang.org/x/oauth2/google
RUN go get -u google.golang.org/api/sheets/v4

RUN go build -o voluum-apis

CMD ["./voluum-apis"]