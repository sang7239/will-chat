# will-chat
![home_page](imgs/Signup.png) <!-- .element height="20%" width="20%" -->

## Overview
Will-chat is a messenger app that resembles popular communication services such as slack/microsoft teams/discord etc. WillChat users can create public/private channels, in which public channels are accessible to all WillChat users, whereas only the specified users can be added to private channels by the channel creator. The service implements authentication, session management and notification features. 

### Front-End
React.js

### Back-End
GOLANG for serverside (go channels for notifications)

MongoDB for users, channels, and messages collections 

Redis for session tokens for authenticating users. 


## Authentication

![Signup_page](imgs/Signup.png) <!-- .element height="40%" width="40%" -->
![Login_page] (imgs/Login.png) <!-- .element height="40%" width="40%" -->

### Sessions Library

sessionid.go: a cryptographically-random, digitally-signed session ID. Use the various crypto packages in the Go standard library, as directed in the comments.
redisstore.go: a session store backed by a redis database. This implements the abstract Store interface that is defined in store.go. For reference, look at the implementation in the memstore.go file, which is already provided for you.
session.go: a set of package-level functions for beginning sessions, getting session IDs and state from an HTTP request, and ending sessions.

### Sign-Up, Sign-In, Sign-Out



## Channels/Messages
![channels page](imgs/channels.png) <!-- .element height="30%" width="30%" -->
![channels page](imgs/chat.png) <!-- .element height="30%" width="30%" -->



