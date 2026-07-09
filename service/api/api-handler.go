package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	//TAG: login
	rt.router.POST("/session", rt.doLogin)

	//TAG: user
	rt.router.PATCH("/me/username", rt.setMyUserName)
	rt.router.PUT("/me/photo", rt.setMyPhoto)

	//TAG: chat
	rt.router.GET("/chats", rt.getMyConversations)
	rt.router.POST("/chats", rt.createChat)

	//TAG: message
	rt.router.GET("/chats/:chatId/messages", rt.getConversation)
	rt.router.POST("/chats/:chatId/messages", rt.sendMessage)
	rt.router.DELETE("/chats/:chatId/messages/:messageId", rt.deleteMessage)
	rt.router.POST("/chats/:chatId/messages/:messageId/forward", rt.forwardMessage)

	//TAG: comment
	rt.router.POST("/chats/:chatId/messages/:messageId/comments", rt.commentMessage)
	rt.router.DELETE("/chats/:chatId/messages/:messageId/comments/:commentId", rt.deleteComment)

	//TAG: group
	rt.router.PATCH("/chats/:chatId/name", rt.setGroupName)
	rt.router.PUT("/chats/:chatId/photo", rt.setGroupPhoto)
	rt.router.POST("/chats/:chatId/members", rt.addToGroup)
	rt.router.DELETE("/chats/:chatId/members", rt.leaveGroup)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
