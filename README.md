
# Image Processing Service
> Proof Of Concept

## Goal
A large image (in widht & height) is uploaded and will be processed so we have 3 different sizes coming out. 

### Steps
1. The image which is uploaded to the nodejs API will be uploaded to a shared volume
2. When the image is uploaded a message will be pushed on the queue
3. The workers which will be written in Golang will pickup the message, take the image and resize/scale them so we have three different sizes. (small, medium, large)
