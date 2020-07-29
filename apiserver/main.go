package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/will-slack/apiserver/events"
	"github.com/will-slack/apiserver/handlers"
	"github.com/will-slack/apiserver/middleware"
	"github.com/will-slack/apiserver/models/messages"
	"github.com/will-slack/apiserver/models/users"
	"github.com/will-slack/apiserver/sessions"
	"gopkg.in/mgo.v2"
	"gopkg.in/redis.v5"
)

const defaultPort = "443"

const (
	apiRoot                = "/v1/"
	apiSummary             = apiRoot + "summary"
	apiUsers               = apiRoot + "users"
	apiSessions            = apiRoot + "sessions"
	apiSessionsMine        = apiRoot + "sessions/mine"
	apiUserMe              = apiRoot + "users/me"
	apiSpecificUser        = apiRoot + "users/"
	apiChannels            = apiRoot + "channels"
	apiSpecificChannel     = apiRoot + "channels/"
	apiSpecificChannelInfo = apiRoot + "channelInfo/"
	apiMessages            = apiRoot + "messages"
	apispecificMessage     = apiRoot + "messages/"
	apiWebSocket           = apiRoot + "websocket"
)

func main() {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	if len(port) == 0 {
		port = defaultPort
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")

	sessionKey := os.Getenv("SESSIONKEY")
	if len(sessionKey) == 0 {
		log.Fatal("NO SESSIONKEY variable set in the env")
	}

	redisAddr := os.Getenv("REDISADDR")
	fmt.Printf("Connecting to redis server at %s...\n", redisAddr)

	redisOptions := redis.Options{
		Addr: redisAddr,
	}

	redisClient := redis.NewClient(&redisOptions)
	sessionStore := sessions.NewRedisStore(redisClient, -1)

	dbAddr := os.Getenv("DBADDR")
	fmt.Printf("Dialing User Database Server at %s...\n", dbAddr)
	dbSession, err := mgo.Dial(dbAddr)
	if err != nil {
		log.Fatalf("Error connecting to Database: %v", err)
	}

	userStore := users.NewMongoStore(dbSession, "production", "users")
	if userStore == nil {
		log.Fatal("Error creating user database")
	}
	messageStore := messages.NewMongoStore(dbSession, "production")
	if messageStore == nil {
		log.Fatal("Error creating messages database")
	}
	notifier := events.NewNotifier()
	ctx := &handlers.Context{
		SessionKey:   sessionKey,
		SessionStore: sessionStore,
		UserStore:    userStore,
		MessageStore: messageStore,
		Notifier:     notifier,
	}
	go ctx.Notifier.Start()

	mux := http.NewServeMux()
	mux.HandleFunc(apiSummary, handlers.SummaryHandler)
	mux.HandleFunc(apiUsers, ctx.UserHandler)
	mux.HandleFunc(apiSessions, ctx.SessionHandler)
	mux.HandleFunc(apiSessionsMine, ctx.SessionMineHandler)
	mux.HandleFunc(apiUserMe, ctx.UserMeHandler)
	mux.HandleFunc(apiSpecificUser, ctx.SpecificUserHandler)
	mux.HandleFunc(apiChannels, ctx.ChannelsHandler)
	mux.HandleFunc(apiSpecificChannel, ctx.SpecificChannelHandler)
	mux.HandleFunc(apiSpecificChannelInfo, ctx.SpecificChannelInfoHandler)
	mux.HandleFunc(apiMessages, ctx.MessagesHandler)
	mux.HandleFunc(apispecificMessage, ctx.SpecificMessageHandler)
	mux.HandleFunc(apiWebSocket, ctx.WebSocketUpgradeHandler)
	http.Handle(apiRoot, middleware.Adapt(mux, middleware.CORS("", "", "", "")))
	fmt.Printf("server is listening at %s...\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, nil))
}
