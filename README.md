
### Sessions Library

sessionid.go: generates a cryptographically-random, digitally-signed session ID using bcrypt package in the standard GO Library.
redisstore.go: a session store backed by a redis database. This implements the abstract Store interface that is defined in store.go.
session.go: a set of package-level functions for beginning sessions, getting session IDs and state from an HTTP request, and ending sessions.


### Sign-Up, Sign-In, Sign-Out

Sign-Up requests for a new user to be stored in the User Collection

Sign-In authenticates users based on email/password and assigns a new session id for the authenticated user, then redirects to main page

Sign-Out removes session id from the current user, then redirects to sign up page of the app


## Channels/Messages
<img src="imgs/Channels.png" align="left" height="400" width="200" >

### Supported Functionalities:

Get all channels a given user is allowed to see (i.e., chanells the user is a member of, as well as all public channels)
Insert a new channel

Get the most recent N messages posted to a particular channel

Update a channel's Name and Description

Delete a channel, as well as all messages posted to that channel

Insert a new message

### Added functionalities for private channels:

Add a user to a channel's Members list

Remove a user from a channel's Members list

<img src="imgs/Chat.png" align="left" height="400" width="400" >

<!-- ABOUT THE PROJECT -->
## About The Project

![Home](imgs/Home.png)

Will-chat is a messenger app that resembles popular communication services such as [slack](https://slack.com/),[microsoft teams](https://www.microsoft.com/en-us/microsoft-365/microsoft-teams/group-chat-software),[discord](https://discord.com/new) etc. WillChat users can create public/private channels, in which public channels are accessible to all WillChat users, whereas only the specified users can be added to private channels by the channel creator. The service implements authentication, session management and notification features. 

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

This is an example of how to list things you need to use the software and how to install them.
*
```go
sessionid.go: generates a cryptographically-random, digitally-signed session ID using bcrypt package in the standard GO Library.
redisstore.go: a session store backed by a redis database. This implements the abstract Store interface that is defined in store.go.
session.go: a set of package-level functions for beginning sessions, getting session IDs and state from an HTTP request, and ending sessions
```

### Sign-Up, Sign-In, Sign-Out

1. Get a free API Key at [https://example.com](https://example.com)
2. Clone the repo
```sh
git clone https://github.com/your_username_/Project-Name.git
```
3. Install NPM packages
```sh
npm install
```
4. Enter your API in `config.js`
```JS
const API_KEY = 'ENTER YOUR API';
```


<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request





