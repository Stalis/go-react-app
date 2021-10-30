############################
# STEP 1 build react app
############################
FROM node:alpine as frontend-builder

WORKDIR /app
COPY ./web/app .

RUN yarn
RUN yarn build

############################
# STEP 2 build golang server
############################
FROM golang:alpine as server-builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Add certificates and timezones
RUN apk add ca-certificates tzdata

WORKDIR /app
COPY ./cmd ./internal ./migrations ./go.mod ./

# Disable cgo for standalone binary
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go mod download -x
RUN go mod tidy
RUN go build -o /app/bin/server ./cmd/server/main.go
RUN cp -r /app/migrations /app/bin/migrations

############################
# STEP 3 build a small image
############################
FROM scratch

WORKDIR /app/frontend
COPY --from=frontend-builder /app/build .

WORKDIR /app
COPY --from=server-builder /app/bin .
COPY --from=server-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=server-builder /usr/share/zoneinfo /usr/share/zoneinfo/

ENV FRONTEND_PATH=/app/frontend
ENV FRONTEND_INDEX=index.html

EXPOSE 80

ENTRYPOINT [ "/app/server" ]