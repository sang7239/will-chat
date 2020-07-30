# will-chat
<img src="imgs/Home.png" align="left" height="400" width="400" >

<br/><br/>
## Overview
Will-chat is a messenger app that resembles popular communication services such as slack/microsoft teams/discord etc. WillChat users can create public/private channels, in which public channels are accessible to all WillChat users, whereas only the specified users can be added to private channels by the channel creator. The service implements authentication, session management and notification features. 

### Front-End
React.js

### Back-End
GOLANG for serverside (go channels for notifications)

MongoDB for users, channels, and messages collections 

Redis for session tokens for authenticating users. 


## Authentication

<img src="imgs/Signup.png" align="left" height="400" width="250" >
<br/><br/>
<img src="imgs/Login.png" align="left" height="200" width="400" >
<br/><br/>

### Sessions Library

sessionid.go: generates a cryptographically-random, digitally-signed session ID using bcrypt package in the standard GO Library.
redisstore.go: a session store backed by a redis database. This implements the abstract Store interface that is defined in store.go.
session.go: a set of package-level functions for beginning sessions, getting session IDs and state from an HTTP request, and ending sessions.

### Sign-Up, Sign-In, Sign-Out




## Channels/Messages
![Channels page](imgs/Channels.png) <!-- .element height="30%" width="30%" -->
![Chat page](imgs/Chat.png) <!-- .element height="30%" width="30%" -->



