# WebSocket Chat Server

This project is a WebSocket-based chat server where users can register and communicate in real time via WebSocket connections. It allows users to register with a name and automatically generates a unique 4-digit user ID (UID) for each registered user. Once registered, users can connect to a WebSocket endpoint using their UID and send messages to all other connected clients.

## Features

- **User Registration**: Users can register by providing their name, and the server generates a unique 4-digit UID for each user.
- **WebSocket Communication**: Registered users can connect to the server using WebSocket, send messages, and receive broadcasted messages in real-time.
- **Broadcasting**: Messages from any connected user are broadcast to all other connected users.
- **Thread-Safe**: The server handles concurrent WebSocket connections using mutexes to ensure thread-safe access to user data.

## Endpoints

### `POST /register`
- **Description**: Registers a new user.
- **Request Body**:
    ```json
    {
      "name": "John Doe"
    }
    ```
- **Response**:
    ```json
    {
      "uid": "1234"
    }
    ```
  The `uid` field contains the unique 4-digit ID generated for the user.

### `GET /ws?uid=<UID>`
- **Description**: Establishes a WebSocket connection for the user with the specified `UID`.
- **WebSocket Communication**:
    - Once connected, users can send messages to the server, which will broadcast the message to all connected clients.
    - Messages are expected in plain text (UTF-8 encoded).
    - Example message from a user:
        ```
        Hello, world!
        ```

### `GET /dummy`
- **Description**: A simple endpoint that returns "Hello, world!" for testing purposes.

## How to Run with Docker

### Prerequisites
- Docker and Docker Compose installed.

### Steps

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Mahaveer86619/chat-ws.git
   cd chat-ws

2. **Build and Start the Docker Container**:
   ```bash
   docker-compose up -d
   ```