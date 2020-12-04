FROM golang:1.15 as builder
ARG SOURCE_LOCATION=/api
WORKDIR ${SOURCE_LOCATION}

# Copy and download dependency using go mod
RUN go get github.com/gorilla/mux
RUN go get go.mongodb.org/mongo-driver/mongo 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
ARG SOURCE_LOCATION=/api
RUN apk --no-cache add curl
EXPOSE 8080
WORKDIR /root/
COPY --from=builder ${SOURCE_LOCATION} .
CMD ["./app"]  