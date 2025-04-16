package handler

// import (
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// )

// const (
// 	authorizationHeader = "Authorization"
// 	userCtx = "userId"
// )

// func (h *Handler) userIdentity(c *gin.Context){
// 	header := c.GetHeader(authorizationHeader)
// 	if header == ""{
// 		newErrorResponce(c, http.StatusUnauthorized, "empty auth header")
// 		return
// 	}

// 	headerParts:= strings.Split(header, " ")
// 	if len(headerParts) != 2{
// 		newErrorResponce(c, http.StatusUnauthorized, "invalid auth header")
// 		return
// 	}

// 	userId, err := h.services.Authorization.ParseToken(headerParts[1])
// 	if err!=nil{
// 		newErrorResponce(c, http.StatusUnauthorized, err.Error())
// 		return
// 	}
// 	c.Set(userCtx, userId)
// }