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

### Feedback and Mentorship

I’m actively seeking feedback on this project and welcome any suggestions or constructive criticism to help improve my software engineering skills.  
If you have experience in software development or related fields and would like to share insights or advice, please feel free to reach out.

Additionally, I am looking for mentors to support my learning journey in software engineering and infrastructure.  
If you are interested in mentoring, guiding, or collaborating with someone eager to grow and develop professionally, I would be grateful to connect.  
You can contact me through GitHub or via the project discussion channels.