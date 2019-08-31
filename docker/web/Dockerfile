FROM golang:latest as builder

LABEL maintainer="Vance Longwill <vancelongwill@gmail.com>"

WORKDIR /app
#
# Build Args
# ARG LOG_DIR=/app/logs

# Create Log Directory
# RUN mkdir -p ${LOG_DIR}

# Environment Variables
# ENV LOG_FILE_LOCATION=${LOG_DIR}/app.log

COPY go.mod go.sum ./

# Enable go modules
ENV GO111MODULE=on

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# VOLUME [${LOG_DIR}]

# Command to run the executable
CMD ["./main"] 
