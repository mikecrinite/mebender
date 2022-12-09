# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

# Install ffmpeg
RUN apk add ffmpeg
# Install imagemagick
RUN apk add imagemagick

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./

RUN go mod download
COPY . .

RUN go build -o ./mebender

EXPOSE 8080

CMD [ "./mebender" ]

