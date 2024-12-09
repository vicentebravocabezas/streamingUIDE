package media

import (
	"errors"
	"slices"
)

type Media interface {
	Id() int

	Title() string

	Description() string

	MediaType() string

	Source() string

	Artist() string

	Album() string
}

func MediaList() ([]Media, error) {
	movieList, err := MovieList()
	if err != nil {
		return nil, err
	}

	musicList, err := MusicList()
	if err != nil {
		return nil, err
	}

	return slices.Concat(movieList, musicList), nil
}

func GetMedia(id string, mediaType string) (Media, error) {
	if mediaType == "Movie" {
		return GetMovie(id)
	} else if mediaType == "Music" {
		return GetMusic(id)
	}

	return nil, errors.New("Invalid media type")
}
