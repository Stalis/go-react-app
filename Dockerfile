FROM golang:latest

WORKDIR /app

COPY ./frontend/build ./frontend
COPY ./server ./src

RUN cd src
RUN go build -o ../server

RUN cd ..

ENV SERVER_HOST=localhost
ENV SERVER_PORT=80
ENV FRONTEND_PATH=frontend
ENV FRONTEND_INDEX=index.html

CMD [ "/app/server" ]

