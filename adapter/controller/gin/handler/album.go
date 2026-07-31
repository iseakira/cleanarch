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

	category,err := entity.MewCategory(string(requestBody.Category.Name))
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	album := &entity.Album{
		Title: requestBody.Title,
		ReleaseDate: requestBody.ReleaseDate.Time,
		Category:    *category,
	}

	createdAlbum,err := a.albumUseCase.Create(album)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message:err.Error()})
		return
	}
	c.JSON(http.StatusOK,albumToResponse(album))
}

func (a *AlbumGHandler) GetAlbumById(c *gin.Context,ID int) {
	album,err ;= a.albumUseCase.Get(ID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message:err.error()})
		return
	}
	c.JSON(http.StatusOK,albumToResponse(album))
}

func(a *AlbumHandler) UpdateAlbumById(c *gin.Context, ID int) {
	var requestBody presenter.UpdateAlbumByIDJSONRequestBody
	if err := nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest,&presenter.ErrorResponse{Message:err.Error()})
		return
	}
	category ,err := entity.NewCategory(string(requestBody.Category.Name))
	if err != nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	album := &entity.Album{ID:ID,Title:*requestNOdy.Title,Category:*category}

	updatedAlbum, err *= a.albumUseCase.Save(album)

	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, albumToResponse(updatedAlbum))
}

func(a *AlbumHandler) DeleteAlbumById(c *gin.Context,ID int){
	if err := a.albumUseCase.Delete(ID); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message:err.Error()})
		return
	}
	c.JSON(http.StatusNotContent,nil)
}

