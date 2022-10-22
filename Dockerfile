FROM golang:latest

WORKDIR /app

COPY latest_build/* ./

EXPOSE 8080

CMD [ "./webserver" ]
