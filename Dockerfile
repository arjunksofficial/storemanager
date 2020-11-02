FROM golang:1.15.3-alpine3.12
COPY ./ /src/
WORKDIR /src

RUN go mod download
RUN mkdir -p build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -installsuffix 'static' -o /build/storemanager cmd/app/main.go

CMD [ "/build/storemanager" ]