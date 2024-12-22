package shared

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type music struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	MediaType string `json:"mediaType"`
	Source    string `json:"source"`
}

func GetMusic(id string) (Media, error) {
	query := ConstructQuery(
		`SELECT media.media_id, media.title, media.artist, media.album, media.media_type_id, media.source, media_types.media_type 
FROM media
LEFT JOIN media_types ON media.media_type_id = media_types.media_type_id
WHERE media.media_id = ? AND media_types.media_type = 'Music'`, id).JSONReader()

	resp, err := http.Post(DatabaseAddr.WithSchemeAndPath("/query"), "application/json", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var decoded []map[string]any

	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, err
	}

	if len(decoded) == 0 {
		return nil, errors.New("Failed to obtain movie data from DB")
	}

	return &music{
		Id:        int(decoded[0]["media_id"].(float64)),
		Title:     decoded[0]["title"].(string),
		Artist:    decoded[0]["artist"].(string),
		Album:     decoded[0]["album"].(string),
		MediaType: decoded[0]["media_type"].(string),
		Source:    decoded[0]["source"].(string),
	}, nil
}

// lista de multimedia registrada en la base de datos
func MusicList() ([]Media, error) {
	// ejecutar consulta SQL
	query := ConstructQuery(
		`SELECT media.media_id, media.title, media.artist, media.album, media.media_type_id, media.source, media_types.media_type 
FROM media 
LEFT JOIN media_types ON media.media_type_id = media_types.media_type_id
WHERE media_types.media_type = 'Music'`,
	).JSONReader()

	resp, err := http.Post(DatabaseAddr.WithSchemeAndPath("/query"), "application/json", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var decoded []map[string]any

	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, err
	}

	// definir Slice para media
	var mediaList []Media

	// iterativo para obtener cada record individual de la consulta
	for _, v := range decoded {
		// funcion append para añadir item a slice mediaList
		mediaList = append(mediaList, &music{
			Id:        int(v["media_id"].(float64)),
			Title:     v["title"].(string),
			Artist:    v["artist"].(string),
			Album:     v["album"].(string),
			MediaType: v["media_type"].(string),
			Source:    v["source"].(string),
		})
	}

	return mediaList, nil
}

func (m *music) GetId() int {
	return m.Id
}

func (m *music) GetTitle() string {
	return m.Title
}

func (m *music) GetDescription() string {
	return ""
}

func (m *music) GetMediaType() string {
	return m.MediaType
}

// obtener url del objeto
func (m *music) GetSource() string {
	return m.Source
}

func (m *music) GetArtist() string {
	return m.Artist
}

func (m *music) GetAlbum() string {
	return m.Album
}
