package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-api-arch-clean-template/adapter/controller/gin/presenter"
	"go-api-arch-clean-template/entity"
	"go-api-arch-clean-template/pkg/logger"
	"go-api-arch-clean-template/usecase"
)

type AlbumHandler struct {
	albumUseCase usecase.AlbumUseCase
}

func NewAlbumHandler(albumUseCase usecase.AlbumUseCase) *AlbumHandler {
	return &AlbumHandler{
		albumUseCase: albumUseCase,
	}
}

func albumToResponse(album *entity.Album) *presenter.AlbumResponse {
	return &presenter.AlbumResponse{
		Id:          album.ID,
		Title:       album.Title,
		ReleaseDate: presenter.ReleaseDate{Time: album.ReleaseDate},
		Category: presenter.Category{
			Id:   &album.Category.ID,
			Name: presenter.CategoryName(album.Category.Name),
		},
	}
}

func (a *AlbumHandler) CreateAlbum(c *gin.Context) {
	var requestBody presenter.CreateAlbumJSONRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
		return
	}
}