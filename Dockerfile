FROM golang:1.10 AS build

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o webstalker .


FROM scratch
WORKDIR /
COPY --from=build /go/src/app/webstalker /
ENTRYPOINT ["/webstalker"]
