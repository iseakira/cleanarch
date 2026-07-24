package usecase

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go-api-arch-clean-template/entity"
	"go-api-arch-clean-template/pkg"
)

type mockAlbumRepository struct {
	mock.Mock
}

func NewMockAlbumRepository() *mockAlbumRepository {
	return &mockAlbumRepository{}
}

func (m *mockAlbumRepository) Create(album *entity.Album) (*entity.Album, error) {
	args := m.Called(album)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Album), args.Error(1)
}

func (m *mockAlbumRepository) Get(ID int) (*entity.Album, error) {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Album), args.Error(1)
}
func (m *mockAlbumRepository) Save(album *entity.Album) (*entity.Album, error) {
	args := m.Called(album)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Album), args.Error(1)
}

func (m *mockAlbumRepository) Delete(ID int) error {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return nil
}

type AlbumUseCaseSuite struct {
	suite.Suite
	albumUseCase *albumUseCase
}

func TestAlbumUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumUseCaseSuite))
}

func (suite *AlbumUseCaseSuite) SetupSuite() {
}

func (suite *AlbumUseCaseSuite) TestCreate() {
	now := pkg.Str2time("2023-01-01")
	title := "album"
	categoryName := "sports"
	mockAlbumRepository := NewMockAlbumRepository()
	suite.albumUseCase = NewAlbumUseCase(mockAlbumRepository)

	category := entity.Category{Name: entity.CategoryName(categoryName)}
	album := &entity.Album{
		Title:       title,
		ReleaseDate: now,
		Category:    category,
	}

	mockAlbumRepository.On("Create", album).Return(&entity.Album{
		ID:          1,
		Title:       title,
		ReleaseDate: now,
		CategoryID:  1,
		Category:    category,
	}, nil)

	album, err := suite.albumUseCase.Create(album)
	suite.Assert().Nil(err)
	suite.Assert().Equal(title, album.Title)
	suite.Assert().Equal(now, album.ReleaseDate)
	suite.Assert().Equal(category.ID, album.Category.ID)
	suite.Assert().Equal(string(categoryName), string(album.Category.Name))
}

func (suite *AlbumUseCaseSuite) TestGet() {
	now := pkg.Str2time("2023-01-01")
	title := "album"
	categoryName := "sports"
	mockAlbumRepository := NewMockAlbumRepository()
	suite.albumUseCase = NewAlbumUseCase(mockAlbumRepository)
	category := &entity.Category{
		ID:   1,
		Name: entity.CategoryName(categoryName),
	}
	mockAlbumRepository.On("Get", 1).Return(&entity.Album{
		ID:          1,
		Title:       title,
		ReleaseDate: now,
		CategoryID:  category.ID,
		Category:    *category,
	}, nil)

	album, err := suite.albumUseCase.Get(1)
	suite.Assert().Nil(err)
	suite.Assert().Equal(title, album.Title)
	suite.Assert().Equal(now, album.ReleaseDate)
	suite.Assert().Equal(category.ID, album.Category.ID)
	suite.Assert().Equal(string(categoryName), string(album.Category.Name))
}

func (suite *AlbumUseCaseSuite) TestUpdate() {
	now := pkg.Str2time("2023-01-01")
	title := "album"
	categoryName := "sports"
	category := entity.Category{Name: entity.CategoryName(categoryName)}
	album := &entity.Album{
		ID:          1,
		Title:       title,
		ReleaseDate: now,
		Category:    category,
	}

	mockAlbumRepository := NewMockAlbumRepository()
	suite.albumUseCase = NewAlbumUseCase(mockAlbumRepository)
	mockAlbumRepository.On("Save", album).Return(&entity.Album{
		ID:          1,
		Title:       title,
		ReleaseDate: now,
		CategoryID:  1,
		Category:    category,
	}, nil)

	album, err := suite.albumUseCase.Save(album)
	suite.Assert().Nil(err)
	suite.Assert().Equal(title, album.Title)
	suite.Assert().Equal(now, album.ReleaseDate)
	suite.Assert().Equal(category.ID, album.Category.ID)
	suite.Assert().Equal(string(categoryName), string(album.Category.Name))
}

func (suite *AlbumUseCaseSuite) TestDelete() {
	mockAlbumRepository := NewMockAlbumRepository()
	suite.albumUseCase = NewAlbumUseCase(mockAlbumRepository)
	mockAlbumRepository.On("Delete", 1).Return(nil, nil)
	err := suite.albumUseCase.Delete(1)
	suite.Assert().Nil(err)
}