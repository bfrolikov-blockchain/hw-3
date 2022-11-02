FROM golang:1.19.1-alpine3.16
RUN apk add --no-cache gcc musl-dev
COPY . /hw-3
WORKDIR /hw-3
RUN go mod download
RUN go build
ENTRYPOINT ["./hw-3"]