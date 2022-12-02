FROM golang:1.15-alpine

# Create a working directory and copy the source files into it
WORKDIR /app
COPY . .

# Build the battlesnake game
RUN go build -o battlesnake

# Run the battlesnake game when the container starts
ENTRYPOINT ["/app/battlesnake"]
