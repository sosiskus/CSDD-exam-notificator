## Build
FROM golang:1.22 AS build

WORKDIR /app

# Download dependencies
COPY ./go.mod ./
COPY ./go.sum ./
# RUN go mod download
COPY *.go ./

RUN go build ./csdd.go

CMD ["./csdd"]

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /app/csdd /csdd

EXPOSE 80/tcp

USER root:root

ENTRYPOINT ["/csdd"]
