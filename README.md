# LL-HLS Streaming

This project setups OvenMediaEngine, NGINX and a mapper service to stream video and audio using the well know low latency protocol called LL-HLS.

The project has a few components:
* OvenMediaEngine: The media server that will stream the video and audio.
* NGINX: The proxy/web server that will deliver audio and video, caching everything that is possible.
* Mapper: A service that will map the requested streaming to the right origin.

## Run

To run the project and all the components, you can use the following command:

```
just run
```

## Ingest

To ingest a stream, you can use the following commands:

### SRT

```
just ingest-srt
```

You can also use any other client that supports SRT to ingest the stream.

## Play

Play using the following URL http://localhost:8080/ll/llhls/playlist.m3u8
