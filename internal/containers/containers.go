package containers

import "errors"

var (
	ErrNotFound      = errors.New("container not found")
	ErrAlreadyExists = errors.New("container already exists")
)

type Container struct {
	ID string `json:"containerid"`
}

type CreateContainerRequest struct {
	ImageURL string `json:"image_url" binding:"required"`
}
