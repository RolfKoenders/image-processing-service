
# Image Processing Service
> Proof Of Concept

## Goal
A large image (in widht & height) is uploaded and will be processed so we have different sizes coming out. 

### Plan / Approach
1. The image which is uploaded to the nodejs API will be uploaded to a shared volume
2. When the image is uploaded a message will be pushed on the queue
3. The workers which will be written in Golang will pickup the message, take the image and resize/scale them so we have three different sizes. (small, medium, large)

## Starting up
Before starting everything up with docker compose, make sure to launch the `setup.sh` script. This script will create the missing folder where the uploaded images will be stored and install all the dependencies for the nodejs api.

```bash
# In the project dir
$ ./setup.sh
```

Then we can bring everything up with docker compose
```bash
$ docker-compose up --build
```

This will build and spin up the following container
- RabbitMQ
- Nginx
- API
- Worker

### RabbitMQ
[RabbitMQ](https://www.rabbitmq.com/) is the most widely deployed open source message broker. The management portal is listening on port `15672`

### Nginx
Proxy layer for the API and used to serve static files. Open a browser window and go to http://localhost:8080 to see the upload form.

### API
Simple API layer written in [node.js](https://nodejs.org/en/). This is the API where the submitted form will end up.

### Worker
The worker which is going to process the images and takes care of resizing them. Since there is no reconnect mechanism in place there is a chance that the worker container is in exitted state. With the command `docker-compose restart worker` its possible to bring him to life. 
