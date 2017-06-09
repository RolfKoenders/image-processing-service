'use strict';

const restify = require('restify');
const plugins = require('restify-plugins');

const config = require('./config');

const server = restify.createServer({
    name: config.name,
    version: '1.0.0'
});

server.use(plugins.acceptParser(server.acceptable));
server.use(plugins.queryParser());
server.use(plugins.bodyParser());

// Startup server
server.listen(config.port, function () {
    console.log('%s listening at port %s', server.name, config.port);
});