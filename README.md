# GASS - Gait Analysis from Single camera Setup
## Introduction 
This project is the backend server for GASS (Gait Analysis from Single camera Setup).
## Requirements
[Docker](https://www.docker.com/)
## Usage
First build the image with Docker
```
docker build --platform linux/amd64 -t gass-server:v1 -f Dockerfile . 
```
Then run it using docker run. You will also have to fill the environment variables with your own values. Also please note that the database provided needs to be Postgres.
```
docker run -d -p <port>:<port> \
-e PORT="PORT" \
-e DB_HOST="DB_HOST_IP" \
-e DB_USER="DB_USERNAME" \
-e DB_PASSWORD="DB_PASSWORD" \
-e DB_NAME="DB_NAME" \
-e DB_PORT="DB_PORT" \
-e TOKEN_TTL="JWT_TOKEN_TTL" \
-e JWT_PRIVATE_KEY="JWT_TOKEN_PRIVATE_KEY" \
-e INFERENCE_SERVER="http://INFERENCE_SERVER_IP" \
--name gass-server gass-server:v1
```
