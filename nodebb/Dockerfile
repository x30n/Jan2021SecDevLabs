FROM node:current-slim

WORKDIR /usr/src/app
COPY package.json .
RUN npm install

COPY people.js /data/

EXPOSE 8080
CMD [ "npm", "start" ]

COPY . .
