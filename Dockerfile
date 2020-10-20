# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.12
ARG GO_VERSION=1.13

# First stage: Build the binary
FROM golang:${GO_VERSION}-alpine as golang

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Copy Go Module config
COPY go.mod .
COPY go.sum .

# Download Go Modules
RUN go mod download

# Import the code from the context.
COPY . .

# Static build required so that we can safely copy the binary over.
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"'

# Build app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app .

# Set permissions on app
RUN chmod +x ./app

# Second stage: Setup environment
FROM alpine:latest as alpine

# Install dependencies
RUN apk --no-cache add tzdata zip ca-certificates

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Set workdir to zoneinfo
WORKDIR /usr/share/zoneinfo

# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -q -r -0 /zoneinfo.zip .

# Final stage: the running container.
FROM scratch

# Set timezone data
ENV ZONEINFO /zoneinfo.zip

# Import zone information
COPY --from=alpine /zoneinfo.zip /

# Import the user and group files from the first stage.
COPY --from=alpine /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --chown=nobody:nobody --from=golang /src/app /app

# Perform any further action as an unprivileged user.
USER nobody:nobody

ENTRYPOINT ["/app"]