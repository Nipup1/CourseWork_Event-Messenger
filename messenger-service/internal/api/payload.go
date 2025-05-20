package api

import (
	"go/course_work_6/internal/models"

	ssov6 "github.com/Nipup1/SSO_Protos/gen/go/sso"
)

//-------------------------------------------------------

type HandleChatsRequest struct {
	Name      string `json:"name"`
	MemberIDs []int  `json:"member_ids"`
}

type HandleChatsResponse struct {
	ChatId int `json:"chat_id"`
}

//-------------------------------------------------------

type HandleChatsListResponse struct {
	Chats []models.Chat `json:"chats"`
}

//-------------------------------------------------------

type HandleMessagesResponse struct {
	Messages []models.Message `json:"messages"`
}

//-------------------------------------------------------

type HandleSendMessageRequest struct {
	Content string `json:"content"`
}

//-------------------------------------------------------

type HandleUserListResponse struct{
	UserIDs []int `json:"user_ids"`
}

//-------------------------------------------------------

type RegisterRequest struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	UserId int `json:"user_id"`
}

//-------------------------------------------------------

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

//-------------------------------------------------------

type UsersResponse struct {
	Users []*ssov6.User `json:"users"`
}

//-------------------------------------------------------

type UserResponse struct{
	UserID int64 `json:"user_id"`
	Email string `json:"email"`
	FullName string `json:"full_name"`
}
