package shared

import (
	"errors"
	"slices"
)

type Media interface {
	GetId() int

	GetTitle() string

	GetDescription() string

	// si es una pelicula, "Movie", o una canción, "Music"
	GetMediaType() string

	//url del archivo
	GetSource() string

	GetArtist() string

	GetAlbum() string
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

type MediaGlobalType struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaType   string `json:"mediaType"`
	Source      string `json:"source"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
}

func (m *MediaGlobalType) GetId() int {
	return m.Id
}

func (m *MediaGlobalType) GetTitle() string {
	return m.Title
}

func (m *MediaGlobalType) GetDescription() string {
	return m.Description
}

func (m *MediaGlobalType) GetMediaType() string {
	return m.MediaType
}

func (m *MediaGlobalType) GetSource() string {
	return m.Source
}

func (m *MediaGlobalType) GetArtist() string {
	return m.Artist
}

func (m *MediaGlobalType) GetAlbum() string {
	return m.Album
}
