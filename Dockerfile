FROM golang:1.11

WORKDIR /go/src/github.com/JonathonGore/api-check
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-check .

FROM alpine:latest  

WORKDIR /root/
COPY --from=0 /go/src/github.com/JonathonGore/api-check/api-check .
CMD ["./api-check"] 
