package handlers

import "dev-book-api/services"

type Handler struct {
	UserService services.UserService
	PostService services.PostService
}
