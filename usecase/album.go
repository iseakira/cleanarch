package usecase

import (
	"go-api-arch-clean-template/adapter/gateway"
	"go-api-arch-clean-template/entity"
)


type (
	AlbumUseCase interface {
		Create(album *entity.Album) (*entity.Album, error)
		Get(ID int) (*entity.Album, error)
		Save(*entity.Album) (*entity.Album, error)
		Delete(ID int) error
	}
)

type albumUseCase struct {
	albumRepository gateway.AlbumRepository
}

func NewAlbumUseCase(albumRepository gateway.AlbumRepository) AlbumUseCase {
	return &albumUseCase{
		albumRepository: albumRepository,
	}
}

func (a *albumUseCase) Create(album *entity.Album) (*entity.Album, error) {
	return a.albumRepository.Create(album)
}

func (a *albumUseCase) Get(ID int) (*entity.Album, error) {
	return a.albumRepository.Get(ID)
}

func (a *albumUseCase) Save(album *entity.Album) (*entity.Album, error) {
	return a.albumRepository.Save(album)
}

func (a *albumUseCase) Delete(ID int) error {
	return a.albumRepository.Delete(ID)
}