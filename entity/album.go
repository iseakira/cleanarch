package entity

import (
	"go-api-arch-clean-template/pkg"
	"time"
)

type Album struct {
	ID int
	Title string
	ReleaseDate time.Time
	CategoryID int
	Category Category
}


func (a *Album) Anniversary(clock pkg.Clock) int {
	now := clock.Now()
	years := now.Year() - a.ReleaseDate.Year()
	releaseDay := pkg.GetAdjustedReleaseDay(a.ReleaseDate,now)
	if now.YearDay() < releaseDay {
		years -= 1
	}
	return years
}