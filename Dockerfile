FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod download && go mod verify

RUN make build

EXPOSE 5465

ENTRYPOINT ["/app/bin/proxy"]