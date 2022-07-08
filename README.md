# Zenon Twitter bot

Twitter bot that publishes updates to Zenon accelerator projects. For the data backend, see [zenon-az-service](https://github.com/dumeriz/zenon-az-service).

## Requirements

Requirements for running this repository:

- Golang or Docker installation
- Twitter API key and access token
- [Zenon-ZMQ-backend](https://github.com/dumeriz/zenon-az-service)

## Setup

Before starting the application, we first need to provide the required environment variables:

```bash
API_KEY=
API_KEY_SECRET=
ACCESS_TOKEN=
ACCESS_TOKEN_SECRET=
ZMQ_URL=tcp://localhost:6666
ZENON_URL=ws://127.0.0.1:35998
```

After that, you can start the application using Docker or Golang:

```bash
go run cmd/main.go
# or
docker build -t zenon-twitter .
docker run --name zenon-twitter -d zenon-twitter
```

## Supported events

List of currently supported events:

- Submission of a new A-Z proposal up for a vote
- Final vote of the A-Z proposal
