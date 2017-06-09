'use strict';

const express = require('express');
const fileUpload = require('express-fileupload');
const multer = require('multer');
const bodyParser = require('body-parser');
const path = require('path');
const fs = require('fs-extra');
// const crypto = require('crypto');
const amqp = require('amqplib/callback_api');

const config = require('./config');

let imageChannel = null;

// # Server
const api = express();

// # Multer, file upload handler
const storage = multer.diskStorage({
  destination: config.upload_dir
//   ,filename: function (req, file, cb) {
//     crypto.pseudoRandomBytes(16, function (err, raw) {
//       if (err) return cb(err)

//       cb(null, raw.toString('hex') + path.extname(file.originalname))
//     })
//   }
});

// # Handler
api.post('/api/submit', function (req, res) {
    var uploads = multer({
        storage: storage
    }).array('image');

    uploads(req, res, function (err) {
        if (err) {
            console.log(err.message);
            return res.send(400, err);
        }

        console.log('All the files are uploaded, push messages to the queue');

        for(let i = 0; i < req.files.length; i++) {
            const file = req.files[i];
            imageChannel.sendToQueue(config.imageQueue, new Buffer.from(JSON.stringify({
                path: config.upload_dir,
                image: file.filename,
                mimetype: file.mimetype
            }), 'utf-8'));
        }

        res.send(200);
    });
});

// Setup connection with que & Startup api server
amqp.connect(`amqp://${config.rabbit.host}`, function(err, conn) {
    conn.createChannel(function(err, ch) {
        ch.assertQueue(config.imageQueue, {durable: false});
        imageChannel = ch;
        // Start API server
        api.listen(config.port, function () {
            console.log('Image-Svc-API listening at port %s', config.port);
        });
    });
});
