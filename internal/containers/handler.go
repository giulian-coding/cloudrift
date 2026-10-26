package containers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(router *gin.RouterGroup) {
	router.POST("/create", h.createContainer)
	router.GET("/:id", h.getContainer)
	router.GET("/:id/shell", h.getContainerShell)
}

func (h *Handler) createContainer(ctx *gin.Context) {
	var req CreateContainerRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Wenn das JSON ungültig ist oder image_url fehlt, mit 400 Bad Request abbrechen
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.svc.createContainer(ctx, req.ImageURL)

	ctx.JSON(http.StatusOK, gin.H{"message": "Container wird erstellt", "image": req.ImageURL})
}

func (h *Handler) getContainer(ctx *gin.Context) {

}

func (h *Handler) getContainerShell(ctx *gin.Context) {

}
