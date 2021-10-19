############################
# STEP 1 build react app
############################
FROM node:alpine as frontend-builder

WORKDIR /app

COPY ./frontend .

#RUN npm install -g yarn
RUN yarn
RUN yarn build

############################
# STEP 2 build golang server
############################
FROM golang:alpine as server-builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY ./server .

RUN go get -d -v
RUN CGO_ENABLED=0 go build -o /app/bin/server

############################
# STEP 3 build a small image
############################
FROM scratch

WORKDIR /app/frontend
COPY --from=frontend-builder /app/build .

WORKDIR /app
COPY --from=server-builder /app/bin .

ENV SERVER_PORT=80
ENV FRONTEND_PATH=/app/frontend
ENV FRONTEND_INDEX=index.html

EXPOSE 80

ENTRYPOINT [ "/app/server" ]