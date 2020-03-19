
#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN apk add --no-cache git
RUN go get -d -v ./...
RUN go install -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/go-docker-httpc /app
COPY tmpl ./tmpl/
ENTRYPOINT ./app
LABEL Name=go-docker-httpc Version=0.0.1
EXPOSE 8080
