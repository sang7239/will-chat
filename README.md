
<!-- ABOUT THE PROJECT -->
## Overview

![Home](imgs/Home.png)

Will-chat is a messenger app that resembles popular communication services such as [slack](https://slack.com/), [microsoft teams](https://www.microsoft.com/en-us/microsoft-365/microsoft-teams/group-chat-software), [discord](https://discord.com/new) etc. WillChat users can create public/private channels, in which public channels are accessible to all WillChat users, whereas only the specified users can be added to private channels by the channel creator. The service implements authentication, session management and notification features. 

### Built With
* [GO](https://golang.org/)
* [React](https://reactjs.org)
* [Bootstrap](https://getbootstrap.com/)
* [MongoDB](https://www.mongodb.com/)
* [Redis](https://redis.io/)




<!-- GETTING STARTED -->
## Authentication

This is an example of how you may give instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps.

### Sessions Library

The apiserver/sessions directory contains files for a reusable sessions library that provides digitally-signed session IDs, as well as two different session state stores: one backed by an in-memory cache, and one backed by a redis server

`sessionid.go`: generates a cryptographically-random, digitally-signed session ID using bcrypt package in the standard GO Library.

`redisstore.go`: a session store backed by a redis database. This implements the abstract Store interface that is defined in store.go.

`session.go`: a set of package-level functions for beginning sessions, getting session IDs and state from an HTTP request, and ending sessions

### Sign-Up, Sign-In, Sign-Out

`UsersHandler()`: allows new users to sign-up (POST) or returns all users

If the request method is "POST":
+ Decode the request body into a models.NewUser struct
+ Validate the models.NewUser
+ Ensure there isn't already a user in the UserStore with the same email address
+ Ensure there isn't already a user in the UserStore with the same user name
+ Insert the new user into the UserStore
+ Begin a new session
+ Respond to the client with the models.User struct encoded as a JSON object
If the request method is "GET":
+ Get all users from the UserStore and write them to the response as a JSON-encoded array

`SessionsHandler()`: allows existing users to sign-in

The request method must be "POST"
+ Decode the request body into a models.Credentials struct
+ Get the user with the provided email from the UserStore; if not found, respond with an http.StatusUnauthorized
+ Authenticate the user using the provided password; if that fails, respond with an http.StatusUnauthorized
+ Begin a new session
+ Respond to the client with the models.User struct encoded as a JSON object

`SessionsMineHandler()`: allows authenticated users to sign-out
The request method must be "DELETE"
+ End the session
+ Respond to the client with a simple message saying that the user has been signed out

`UsersMeHanlder()`: Get the session state
+ Respond to the client with the session state's User field, encoded as a JSON object


## Channels

`ChannelsHandler()`: This will handle all requests made to the /v1/channels path.

+ `GET`: get the channels the current user can see and write the returned slice to the response as a JSON-encoded array.
+ `POST`: add the current user to the new channel's Members list, insert the new channel, and write the newly-inserted Channel object to the response

`SpecificChannelHandler()`: This will handle all requests made to the /v1/channels/<channel-id> path. Get the specific channel ID from the last part of the request's URL path 
  
+ `GET`: if the user is allowed to read messages from this channel (user is a member or the channel is public), get the most recent 500 messages from the specified channel and write those to the response.

+ `PATCH`: if the current user is the channel creator, update the specified channel's Name/Description and write the updated Channel object to the response

+ `DELETE`: if the current user is the channel creator, delete the specified channel

+ `LINK`: if the specified channel is private, add the specified user to the Members list of the current channel.

+ `UNLINK`: if the specified channel is private, remove the current user to the Members list of the specified channel.

## Messages

`MessagesHandler()`: This will handle all requests made to the /v1/messages path. What you do will depend on the request method.

+ `POST`: insert the new message, and respond by writing the newly-inserted Message to the response.

### Added functionalities for private channels:

Add a user to a channel's Members list

Remove a user from a channel's Members list

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request





