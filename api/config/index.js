'use strict';

// Define a schema
let config = {
    name: {
        doc: 'Application name',
        default: 'image-svc-api'
    },
    port: {
        doc: "The port to bind the http server",
        format: "port",
        default: 9080,
        env: "HTTP_PORT"
    },
    upload_dir: {
        doc: "Image upload directory",
        format: String,
        default: null,
        env: "IMAGE_UPLOAD_PATH"
    },
    rabbit: {
        host: {
            doc: "Host of rabbitmq",
            format: String,
            default: null,
            env: "RABBIT_HOST"
        },
        port: {
            doc: "Port of rabbitmq",
            format: "port",
            default: "5672",
            env: "RABBIT_PORT"
        }
    },
    imageQueue: {
        doc: "Que name of the processing worker QUEUE",
        format: String,
        default: null,
        env: "IMAGE_SERVICE_QUEUE"
    }
};

var convict = require('convict');
config = convict(config);
config.validate();
module.exports = config.get();