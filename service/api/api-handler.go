package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	//rt.wrap aggiunge un parametro in più, ovvero ctx reqcontext.RequestContext

	//TAG: login
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	//TAG: user
	rt.router.GET("/me", rt.wrap(rt.getMyInformation))
	rt.router.PATCH("/me/username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/me/photo", rt.wrap(rt.setMyPhoto))
	rt.router.GET("/users", rt.wrap(rt.searchUser))

	//TAG: chat
	rt.router.GET("/chats", rt.wrap(rt.getMyConversations))
	rt.router.POST("/chats", rt.wrap(rt.createChat))

	//TAG: message
	rt.router.GET("/chats/:chatId/messages", rt.wrap(rt.getConversation))
	rt.router.POST("/chats/:chatId/messages", rt.wrap(rt.sendMessage))
	rt.router.DELETE("/chats/:chatId/messages/:messageId", rt.wrap(rt.deleteMessage))
	rt.router.POST("/chats/:chatId/messages/:messageId/forward", rt.wrap(rt.forwardMessage))

	//TAG: comment
	rt.router.POST("/chats/:chatId/messages/:messageId/comments", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/chats/:chatId/messages/:messageId/comments/:commentId", rt.wrap(rt.uncommentMessage))

	//TAG: group
	rt.router.PATCH("/chats/:chatId/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/chats/:chatId/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.POST("/chats/:chatId/members", rt.wrap(rt.addToGroup))
	rt.router.DELETE("/chats/:chatId/members", rt.wrap(rt.leaveGroup))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
