# builder image
FROM golang:latest

# Install curl
RUN apt-get install -y curl


WORKDIR /stg

COPY . .
RUN go mod download
RUN go build -o main .


# Expose the Redis + app port
EXPOSE 50003


# Set the command to run your Go application
#CMD ["sleep", "300"]
WORKDIR /stg

#CMD ["sleep", "3000"]
CMD ["sh", "-c", "go run main.go"]




