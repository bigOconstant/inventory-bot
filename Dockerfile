############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Install git.
RUN apk update && apk add --no-cache git
RUN apk add --no-cache --upgrade bash
RUN apk add --no-cache make 


RUN apk add -U --no-cache ca-certificates && update-ca-certificates

WORKDIR $GOPATH/src/mypacgoinventory/
RUN pwd
COPY . .


WORKDIR $GOPATH/src/mypacgoinventory/cmd/app

RUN go get -d -v
WORKDIR $GOPATH/src/mypacgoinventory/

# Build the binary.
RUN make


############################
# STEP 2 build a small image
############################
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=builder /go/src/mypacgoinventory/goinventory /app/goinventory
# Run the hello binary.

COPY settings.json .


ENTRYPOINT ["/app/goinventory"]