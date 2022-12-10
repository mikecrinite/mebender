# Me, Bender

Make short video clips, turn videos to gifs, extract audio from videos, and more, by making API calls to me, Bender!

This service is a WIP. I'm developing it for personal use for reasons including but not limited to:
* Creating gifs from source material to spam my friends with (gif captioning pending).
* Extracting MF-DOOM-esque audio samples from my favorite shows in a simple, repeatable way, to make lots of sci-fi-themed electronic music that nobody but me will ever listen to.

Feel free to submit any thoughts, questions, concerns, etc. via GitHub 

## Table of Contents
* [Table of Contents](#table-of-contents)
* [Requirements](#requirements)
* [Libraries](#libraries)
* [Run](#run)
* [Test](#test)
* [Example Calls and Responses](#example-calls-and-responses)
* [TODO](#todo)

## Requirements
* [golang](https://go.dev/) (obviously) - An open-source programming language supported by Google
* [ffmpeg](https://ffmpeg.org/) - A complete, cross-platform solution to record, convert and stream audio and video. 
* [imagemagick](https://imagemagick.org/index.php) - Convert, Edit, or Compose Digital Images
* [make](https://www.gnu.org/software/make/manual/make.html) - A Unix utility that can run recipes designed in a makefile

`golang`, `ffmpeg`, and `imagemagick` are actually not required if you choose to use the supplied Dockerfile, as it will download them for you as part of the container build process

`make` is of course also not required. The `makefile` contains all the rules, which are basically just shell scripts. Running the individual commands in a terminal will accomplish the same thing

## Libraries
* [go-ffprobe](gopkg.in/vansante/go-ffprobe.v2) - Small library for executing an ffprobe process on a given file and getting an easy to use struct representing the returned ffprobe data.

## Run
* The simplest way to run the service for development purposes is of course:

        go run main.go
* Run via docker/docker-compose for a pre-packaged single-command run option:

        make run-docker

## Test
* It's a work in progress, and unfortunately that means there are no tests yet
* See [TODO](#todo)

## Example Calls and Responses
### Cut Video
* Supply a start time, end time, and the video's location within the input folder in order to trim the video to a shortened clip
* Endpoint: `/cut`
* Method: `POST`
* Example request body:
```
{
    "StartTime": "19m46s",
    "EndTime": "20m7s",
    "VideoLocation": "sample.mkv"
}   
```
* Example response body:
```
{
    "location": "/root/resources/output/sample_clip_1670621753960165200.mkv",
    "success": true,
    "error": null,
    "duration": "15.50 s"
}
```

### Create .gif
* Supply the video location to extract the frames and create a gif
* Endpoint: `/gif`
* Method: `POST`
* Example request body:
```
{
    "VideoLocation": "sample_clip_1670533061723588900.mkv"
}
```
* Example response body:
```
{
    "location": "/root/resources/output/1670649210303948800.gif",
    "success": true,
    "error": null,
    "duration": "225.55 s"
}
```
Note: 225.55s is kind of a ridiculous amount of time to create a gif, but in somewhat of a defense, it's a 20s gif of 212 PNGs. No but really, this needs to be optimized

### Create .gif
* Supply the video location to extract the audio and create a wav file
* Endpoint: `/sound`
* Method: `POST`
* Example request body:
```
{
    "VideoLocation": "sample_clip_1670533061723588900.mkv"
}
```
* Example response body:
```
{
    "location": "",
    "success": true,
    "error": null,
    "duration": "1.07 s"
}
```
Yes, audio is much quicker to extract. Still, 1+ s response time is not ideal

## TODO
- Optimize ffmpeg and imagemagick (the latter especially). Creating the gif takes multiple minutes which is probably excessive
- Pipe ffmpeg frames directly to imagemagick to reduce the need for saving them to the disk and deleting them afterwards
- Remove frames directory when gif is successfully created (should be unnecessary if the above todo is completed first)
- Support other language audio extraction?
- Tests
- Two-stage Dockerfile 
- Validate requests
- Allow a full length video to be supplied to /cut and /gif as well as start and end times to trim them to a shorter clip