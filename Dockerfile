FROM golang:1.17-alpine
RUN apk add build-base

WORKDIR /IranStocksCrawler

COPY go.mod ./
COPY go.sum ./

COPY . .

CMD ["/bin/sh" ,"-c" ,"go mod download"]

RUN go build -o ./stockscrawler

#EXPOSE 1212 

#CMD tail -f /dev/null

ENTRYPOINT [ "./stockscrawler" ]