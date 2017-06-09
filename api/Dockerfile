# Pulling nodejs base image
FROM node:8.0.0-alpine

RUN apk add --no-cache bash netcat-openbsd

# Skip npm install if package.json didn't changed
ADD package.json /tmp/package.json

RUN cd /tmp && \
    npm install --production && \
    rm -f ~/.npmrc

RUN mkdir -p /opt/app && cp -a /tmp/node_modules /opt/app

# Loading application code
WORKDIR /opt/app
ADD . /opt/app

# Expose API port
EXPOSE 9080

CMD ["node", "."]