FROM golang:1.17-buster
RUN mkdir /app
COPY go.mod /app
COPY go.sum /app
WORKDIR /app
RUN go mod download
WORKDIR /
RUN rm -rf /app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/eheim-exporter main.go

FROM alpine:latest
COPY --from=0 /app/bin/eheim-exporter /eheim-exporter
CMD ["/eheim-exporter"]
