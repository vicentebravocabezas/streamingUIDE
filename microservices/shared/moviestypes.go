package shared

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type movie struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaType   string `json:"mediaType"`
	Source      string `json:"source"`
}

func GetMovie(id string) (Media, error) {
	query := ConstructQuery(`SELECT media.media_id, media.title, media.description, media.media_type_id, media.source, media_types.media_type 
FROM media
LEFT JOIN media_types ON media.media_type_id = media_types.media_type_id
WHERE media.media_id = ? AND media_types.media_type = 'Movie'`, id).JSONReader()

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

	return &movie{
		Id:          int(decoded[0]["media_id"].(float64)),
		Title:       decoded[0]["title"].(string),
		Description: decoded[0]["description"].(string),
		MediaType:   decoded[0]["media_type"].(string),
		Source:      decoded[0]["source"].(string),
	}, nil
}

// lista de multimedia registrada en la base de datos
func MovieList() ([]Media, error) {
	// ejecutar consulta SQL
	query := ConstructQuery(
		`SELECT media.media_id, media.title, media.description, media.media_type_id, media.source, media_types.media_type 
FROM media 
LEFT JOIN media_types ON media.media_type_id = media_types.media_type_id
WHERE media_types.media_type = 'Movie'`,
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
		mediaList = append(mediaList, &movie{
			Id:          int(v["media_id"].(float64)),
			Title:       v["title"].(string),
			Description: v["description"].(string),
			MediaType:   v["media_type"].(string),
			Source:      v["source"].(string),
		})
	}

	return mediaList, nil
}

func (m *movie) GetId() int {
	return m.Id
}

func (m *movie) GetTitle() string {
	return m.Title
}

func (m *movie) GetDescription() string {
	return m.Description
}

func (m *movie) GetMediaType() string {
	return m.MediaType
}

// obtener url del objeto
func (m *movie) GetSource() string {
	return m.Source
}

func (m *movie) GetArtist() string {
	return ""
}

func (m *movie) GetAlbum() string {
	return ""
}
