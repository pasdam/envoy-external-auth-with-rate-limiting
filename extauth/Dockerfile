# Build image
########################
FROM golang:1.20.5-alpine3.18 as builder

WORKDIR /var/tmp/app

# Download the dependencies
ADD ./go.mod ./go.sum ./
RUN go mod download

# copy artifacts into the container
ADD ./*.go ./
RUN go build -o .build/app .

# Final image
########################
FROM alpine:3.18.2

WORKDIR /opt/app

COPY --from=builder /var/tmp/app/.build/app .

CMD [ "./app" ]
