FROM golang:1.11

# RUN apk add --no-cache gcc musl-dev

RUN apt-get -y update && apt-get install -y gcc 

COPY ./ /go/src/github.com/user/myProject/app

RUN go get github.com/pradeep-pyro/triangle
RUN go get -d -v ./...
RUN go install -v ./...

