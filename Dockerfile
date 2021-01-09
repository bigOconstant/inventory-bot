############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Install git.
RUN apk update && apk add --no-cache git

RUN apk add -U --no-cache ca-certificates && update-ca-certificates

WORKDIR $GOPATH/src/mypacgoinventory/

COPY . .

# Fetch dependencies.# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/goinventory
RUN ls /go/bin/

############################
# STEP 2 build a small image
############################
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=builder /go/bin/goinventory /app/goinventory
# Run the hello binary.

COPY settings.json .
COPY html/ /app/html/

ENTRYPOINT ["/app/goinventory"]