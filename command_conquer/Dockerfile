## Base image
FROM node:9

## Create base dir
RUN mkdir /src

## Working dir
WORKDIR /src

## Install packages
COPY ./package.json /src/package.json
COPY ./package-lock.json /src/package-lock.json
RUN npm install --silent

## Add application code
COPY ./app /src/app

## Set env to dev
ENV NODE_ENV development

EXPOSE 3000

CMD ["node", "/src/app/cmd_injection.js"]