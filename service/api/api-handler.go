package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/session", rt.doLogin)
	rt.router.PATCH("/me/username", rt.setMyUserName)
	rt.router.PATCH("/me/photo", rt.setMyPhoto)
	rt.router.GET("/me/chats", rt.getMyConversations)
	rt.router.GET("/me/chats/:chat_id/messages", rt.getConversation)
	rt.router.POST("/me/chats/:chat_id/messages", rt.sendMessage)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
