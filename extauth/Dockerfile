# Build image
########################
FROM golang:1.15.7-alpine3.13 as builder

WORKDIR /var/tmp/app

RUN apk add git

# Download the dependencies
ADD ./go.mod ./go.mod
RUN go mod download

# copy artifacts into the container
ADD ./*.go ./
RUN go build -o .build/app .

# Final image
########################
FROM alpine:3.13.0

WORKDIR /opt/app

COPY --from=builder /var/tmp/app/.build/app .

CMD [ "./app" ]
