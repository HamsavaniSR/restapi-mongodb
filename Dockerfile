FROM golang:1.15 AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
EXPOSE 8080
RUN go build -o main ./api
WORKDIR /dist
RUN cp /build/main .
FROM scratch
COPY --from=builder /dist/main /
# Command to run the executable
ENTRYPOINT ["/main"]
