'use strict';

const express = require('express');
const fileUpload = require('express-fileupload');
const multer = require('multer');
const bodyParser = require('body-parser');
const path = require('path');
const fs = require('fs-extra');
const crypto = require('crypto');
const amqp = require('amqplib/callback_api');

const config = require('./config');

let imageChannel = null;

// # Server
const api = express();

// # Multer, file upload handler
const storage = multer.diskStorage({
  destination: config.upload_dir,
  filename: function (req, file, cb) {
    crypto.pseudoRandomBytes(16, function (err, raw) {
      if (err) return cb(err)

      cb(null, raw.toString('hex') + path.extname(file.originalname))
    })
  }
});

// # Handler
api.post('/api/submit', function (req, res) {
    var uploadProfileImgs = multer({
        storage: storage
    }).single('image');

    uploadProfileImgs(req, res, function (err) {
        if (err) {
            console.log(err.message);
            return res.send(400, err);
        }
        console.log('Everything went fine, image is uploaded.');
        
        // Send message to que
        imageChannel.sendToQueue(config.imageQue, new Buffer.from(JSON.stringify({
            path: config.upload_dir,
            image: '1f8abb90e377b32c26ae6190f714e609.jpeg'
        }), 'utf-8'));

        return res.send(200);
    });
});

// Setup connection with que & Startup api server
amqp.connect('amqp://localhost', function(err, conn) {
    conn.createChannel(function(err, ch) {
        const q = config.imageQue;
        ch.assertQueue(q, {durable: false});
        imageChannel = ch;

        // Start API server
        api.listen(config.port, function () {
            console.log('Image-Svc-API listening at port %s', config.port);
        });
    });
});
