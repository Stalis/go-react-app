FROM golang:latest

WORKDIR /app/frontend
COPY ./frontend/build .

WORKDIR /app/src
COPY ./server .
RUN go build -o ../server

WORKDIR /app

ENV SERVER_PORT=80
ENV FRONTEND_PATH=/app/frontend
ENV FRONTEND_INDEX=index.html

EXPOSE 80

CMD [ "/app/server" ]

