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
    }
};

var convict = require('convict');
config = convict(config);
config.validate();
module.exports = config.get();