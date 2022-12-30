# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

# Install ffmpeg
RUN apk add ffmpeg
# Install frei0r-plugins
# RUN apk add frei0r-plugins
# Install imagemagick
RUN apk add imagemagick

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . .

RUN go build -o ./mebender

# For using a custom port
ARG PORT
ENV PORT $PORT
EXPOSE ${PORT?8080}

CMD [ "./mebender" ]

