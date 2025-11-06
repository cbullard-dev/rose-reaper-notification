# Rose Reaper Notification Service

The Rose Reaper Notification Service is a bridge between a Twitch stream and a Discord community.  
It listens for Twitch channel events in real time and automatically notifies the Discord server when the streamer goes live or ends their stream.

### Overview

This service uses Twitch’s WebSocket EventSub connection to subscribe to events such as:

- Stream start (when the streamer goes live)
- Stream end (when the stream finishes)

When these events occur, the service sends a corresponding message to a designated Discord announcements channel.

### Goals

- Provide a reliable notification system for stream activity.  
- Maintain a stable WebSocket connection to Twitch’s EventSub system.  
- Integrate seamlessly with Discord through webhooks or bot interactions.  

### Planned Features

- Basic Twitch stream start/end notifications (initial goal)  
- Rich embed messages in Discord announcements  
- Support for customized notification templates  
- Optional persistence and monitoring for uptime  

### Status

This project is currently in early development and focuses on establishing the initial Twitch WebSocket connection and Discord message delivery system.  
