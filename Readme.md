# avcapture

avcapture allows you to run a that captures the content and pipes the audio/video for any URL to FFmpeg for encoding including generating a live playlist.

## Build and Run

- **Build**: `docker build -f ./Dockerfile -t <imagename>:<tag> .`
- **Run**: `docker run -it --net test -v <directory to store output>:<mapped directory inside docker container> --name avcapture -p 8080:8080 <imagename>:<tag>`

## API

### start-recording

- POST: <http://IP:8080/start_recording>
- "ffmpeg:options" and "chrome:options" are optional. If it is specified, the original arguments to these applications will be replaced completely with the provided one. User has to take care on the arguments passed for proper functionality.

```json
{
  "ffmpeg": {
    "params": [
      [" -hls_time", "6"],
      ["-hls_playlist_type", "event"],
      ["-hls_segment_filename", "/work/out%04d.ts"],
      ["/work/play.m3u8"]
    ],
    "options": [
      ["-y", ""],
      ["-v", "info"],
      ["-f", "x11grab"],
      ["-draw_mouse", "0"],
      ["-r", "24"],
      ["-s", "1280x720"],
      ["-thread_queue_size", "4096"],
      ["-i", ":99.0+0,0"],
      ["-f", "pulse"],
      ["-thread_queue_size", "4096"],
      ["-i", "default"],
      ["-acodec", "aac"],
      ["-strict", "-2"],
      ["-ar", "48000"],
      ["-c:v", "libx264"],
      ["-x264opts", "no-scenecut"],
      ["-preset", "veryfast"],
      ["-profile:v", "main"],
      ["-level", "3.1"],
      ["-pix_fmt", "yuv420p"],
      ["-r", "24"],
      ["-crf", "25"],
      ["-g", "48"],
      ["-keyint_min", "48"],
      ["-force_key_frames", "\"expr:gte(t,n_forced*2)\""],
      ["-tune", "zerolatency"],
      ["-b:v", "3600k"],
      ["-maxrate", "4000k"],
      ["-bufsize", "5200k"]
    ]
  },
  "chrome": {
    "url": "<https://www.youtube.com/watch?v=Bey4XXJAqS8>",
    "options": [
      "--enable-logging=stderr",
      "--autoplay-policy=no-user-gesture-required",
      "--no-sandbox",
      "--start-maximized",
      "--window-position=100,300",
      "--window-size=1280,720"
    ]
  }
}
```

### stop-recording

- POST: <http://IP:8080/stop_recording>
- No parameter is passed to this call.

## Configuration

By default, the application will run server on port **8080**. If the application has to run on a different port, set the environment variable `PORT` with the new port inside the Dockerfile before the `ENTRYPOINT`.
For example:

```dockerfile
ENV PORT=":8080"
```

## Output

User is supposed to map a directory from host system to the docker image. Along with this, user has to provide the output path (as part of `/start_recording` api) which will direct the output generated to the corresponding directory.

## Architecture

The docker image contains google chrome, ffmpeg and wrapper application.
As part of startup, the wrapper application will configure the system to run chrome browser on the given display id and to capture audio using pulseaudio.
Once the `/start_recording` is received, chrome will be started to render the `url` provided. An instance of ffmpeg will be started to capture the display.

## Version Info

| Component         | Version       | Details                                                                   |
| ----------------- | ------------- | ------------------------------------------------------------------------- |
| Base docker image | ubuntu:16.04  |                                                                           |
| ffmpeg            | 4.0.2         | <https://ffmpeg.org/releases/ffmpeg-4.0.2.tar.bz2>                        |
| libx264           | latest        | git://git.videolan.org/x264.git                                           |
| nasm              | 2.13.03       | <https://www.nasm.us/pub/nasm/releasebuilds/2.13.03/nasm-2.13.03.tar.bz2> |
| Google chrome     | Latest stable | deb <http://dl.google.com/linux/chrome/deb/> stable main                  |
| pulseaudio        | 8.0           |

## Known limitations

- The solution is validated on ubuntu 16.04 (with pulseaudio v8.0). Ubuntu 18.04 contains pulseaudio 11.1 which breaks a few features. With pulseaudio 11.1,
- To run pulseaudio as daemon on root, `--system` has to be mentioned.
- User will not be able to add modules to pulseaudio.
