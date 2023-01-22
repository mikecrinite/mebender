# Me, Bender

    Wait... there on the screen. It's that guy you are!

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
* [go-ffprobe](https://gopkg.in/vansante/go-ffprobe.v2) - Small library for executing an ffprobe process on a given file and getting an easy to use struct representing the returned ffprobe data.

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
    "StartTime": "5m55s",
    "EndTime": "6m7s",
    "VideoLocation": "Dragon Ball Z Kai - 56 - I Will Defeat Frieza! Another Super Saiyan!.mkv"
} 
```
* Example response body:
```
{
    "location": "/root/resources/output/1674365770809476600_DragonBallZKai-56-IW_clip.mkv",
    "success": true,
    "error": null,
    "duration": "69.12 s"
}
```
Nice. We also (a) prefix each file with a timestamp so each new file will be last in the folder and you don't have to search for it, and (b) shorten the filenames to 20 non-whitespace characters because long filenames are annoying and they should be unique because of timestamp anyway 

### Create .gif
* Supply the video location to extract the frames and create a gif
* Endpoint: `/gif`
* Method: `POST`
* Example request body:
```
{
    "StartTime": "19m48s",
    "EndTime": "19m57s",
    "VideoLocation": "Dragon Ball - S05E24.mkv",
    "OutputFilename": "he_thinks_he's_funny"
}
```
* Example response body:
```
{
    "location": "/root/resources/output/1674364283408007700_he_thinks_he's_funny.gif",
    "success": true,
    "error": null,
    "duration": "68.40 s"
}
```
Note: This seems like a long time to make a gif, although I don't have a frame of reference as this is the only way I've ever personally done it

### Extract Sound
* Supply the video location to extract the audio and create a wav file
* Endpoint: `/sound`
* Method: `POST`
* Example request body:
```
{
    "StartTime": "13m34s",
    "EndTime": "13m38s",
    "VideoLocation": "Dragon Ball - S03E06.mkv",
    "OutputFilename": "you_dont_scare_me"
}
```
* Example response body:
```
{
    "location": "/root/resources/output/1674108461352890900you_dont_scare_me.wav",
    "success": true,
    "error": null,
    "duration": "0.29 s"
}
```
Yes, audio is much quicker to extract. 

## TODO
- Optimize ffmpeg and imagemagick (the latter especially). Creating the gif takes multiple minutes sometimes which is probably excessive
- Pipe ffmpeg frames directly to imagemagick to reduce the need for saving them to the disk and deleting them afterwards
- Support other language audio extraction?
- Tests
- Two-stage Dockerfile 
- Validate requests
