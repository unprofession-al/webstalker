FROM golang:1.11rc2 AS build
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o webstalker .

FROM alpine
WORKDIR /
COPY --from=build /go/src/app/webstalker /
ENTRYPOINT ["/webstalker"]
