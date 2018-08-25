FROM golang:1.11rc2 AS build
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o webstalker .

FROM alpine
WORKDIR /
COPY --from=build /go/src/app/webstalker /
ENTRYPOINT ["/webstalker"]
