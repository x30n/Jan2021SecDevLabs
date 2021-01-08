docker rm --force bb && docker build --tag bulletinboard:1.0 . && docker run --publish 8000:8080 --detach --name bb bulletinboard:1.0
