#this is only for development
FROM golang:1.9.2-alpine

ARG app_env
ENV APP_ENV $app_env

#install git
RUN apk add --no-cache git mercurial

# install dependency tool
RUN go get -u github.com/golang/dep/cmd/dep

# Fresh for rebuild on code change, no need for production
RUN go get -u github.com/pilu/fresh


COPY . /go/src/github.com/monstar-lab/fr-circle-api

WORKDIR /go/src/github.com/monstar-lab/fr-circle-api

# dep ensure to ensure the availability of required libraries used by the go source
# for development, pilu/fresh is used to automatically build the application everytime you save a Go or template file

CMD dep ensure && fresh

# for production, it just builds and runs the binary
# CMD dep ensure && go build && ./fr-circle-api

EXPOSE 8080