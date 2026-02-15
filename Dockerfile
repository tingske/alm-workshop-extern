#############
### BUILD ###
#############
FROM golang:1.23 AS build

# Everything will be executed in this working directory in the container
WORKDIR /app

# Copy all the Go source code
COPY workshop-service/. .

# Download dependencies
RUN go mod download && go mod tidy

# Run unit tests
RUN go test -v ./...

# Build the Go application and output the binary to path
RUN CGO_ENABLED=0 GOOS=linux go build -v -o bin/workshop

###########
### RUN ###
###########
FROM alpine:3.19 AS run

# Everything will be executed in this working directory in the container
WORKDIR /app

# Copy from the "build" container the binary that we created in there
# We don't need all the source code only the compiled binary is enough, this can drastically reduce the "run" container size
COPY --from=build /app/bin/workshop /app/bin/workshop

# This makes sure that the application will run correctly on OpenShift
RUN chgrp -R 0 /app/bin/workshop && \
    chmod +x /app/bin/workshop && \
    chmod -R g=u /app/bin/workshop

# This a documentation line, it does not actually open up port 3000
EXPOSE 3000

# Command that is executed when running the container
CMD ["/app/bin/workshop"]
