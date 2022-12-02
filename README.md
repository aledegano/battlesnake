# Battlesnake

This repository contains a simple implementation of the [battlesnake game](https://play.battlesnake.com/).

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

You will need the following software installed on your system to build and run the battlesnake game:

- [Go](https://golang.org/): The programming language used to implement the game.
- [Docker](https://www.docker.com/): The containerization platform used to package the game into a Docker image.

### Building the Game

To build the game, follow these steps:

1. Clone this repository onto your local machine:
    ```
    git clone https://github.com/aledegano/battlesnake.git
    ```
2. Build the game using the `go build` command:
    ```
    go build
    ```

### Running the Game

To run the game, follow these steps:

1. Build the game as described in the previous section.
2. Run the game using the `./battlesnake` command:
    ```
    ./battlesnake
    ```
3. The game will start listening for incoming requests on port 8080. You can use a tool like [Postman](https://www.postman.com/) to send requests to the game and see how it responds.

## Deployment

To deploy the game to a Kubernetes cluster, follow these steps:

1. Build the game as described in the previous section.
2. Build a Docker image for the game using the `docker build` command:
    ```
    docker build -t battlesnake:latest .
    ```
3. Push the Docker image to a Docker registry, such as [Docker Hub](https://hub.docker.com/):
    ```
    docker push battlesnake:latest
    ```
4. Create a Kubernetes deployment for the game using the `kubectl create deployment` command:
    ```
    kubectl create deployment battlesnake --image=battlesnake:latest
    ```
5. Expose the game deployment as a Kubernetes service using the `kubectl expose` command:
    ```
    kubectl expose deployment battlesnake --type=LoadBalancer --port=8080
    ```
6. The game should now be accessible at the service's external IP address. You can use the `kubectl get services` command to find the external IP address for the `battlesnake` service.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

