# RTMP Stream Switcher

A lightweight Go program that acts as an RTMP server and forwards an incoming stream to an external RTMP destination. It supports primary and backup stream sources, automatically switching to the backup if the primary disconnects, and back to primary when it reconnects.

## Features

- Accepts RTMP streams via custom stream keys  
- Forwards incoming stream to a configured output RTMP URL  
- Automatically switches between a primary and backup stream source  
- Uses the [`joy4`](https://github.com/nareix/joy4) library for RTMP handling  

## How It Works

- The server listens on port `:1935` for incoming RTMP connections.  
- It checks the stream key against `FirstStreamKey` and `SecondStreamKey` from the config.  
- Only one stream (primary or backup) can be active at a time:  
  - If the primary connects, the backup is dropped.  
  - If the backup connects and the primary is inactive, it is forwarded.  
- Streams are pushed to the `OutputStreamURL`.  

## Configuration

The program uses a `config.json` file to load settings. Here's a sample format:

```json
{
  "FirstStreamKey": "/primary",
  "SecondStreamKey": "/backup",
  "OutputStreamURL": "rtmp://example.com/live/stream"
}
