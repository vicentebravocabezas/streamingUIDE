package media

import (
	"github.com/vicentebravocabezas/streamingUIDE/database"
)

type movie struct {
	id          int
	title       string
	description string
	mediaType   string
	source      string
}

func GetMovie(id string) (Media, error) {
	row := database.DB().QueryRow(
		`SELECT media.media_id, media.title, media.description, media.media_type_id, media.source, media_types.media_type 
FROM media
LEFT JOIN media_types ON media.media_type_id = media_types.media_type_id
WHERE media.media_id = ?`, id,
	)

	var resultid int
	var title string
	var description string
	var media_type_id int
	var source string
	var media_type string
	if err := row.Scan(&resultid, &title, &description, &media_type_id, &source, &media_type); err != nil {
		return nil, err
	}

	return &movie{
		id:          resultid,
		title:       title,
		description: description,
		mediaType:   media_type,
		source:      source,
	}, nil
}

// lista de multimedia registrada en la base de datos
func MovieList() ([]Media, error) {
	// ejecutar consulta SQL
	rows, err := database.DB().Query(
		`SELECT media.media_id, media.title, media.description, media.media_type_id, media.source, media_types.media_type 
FROM media 
LEFT JOIN media_types ON media.media_type_id = media_types.media_type_id
WHERE media_types.media_type = 'Movie'`,
	)
	if err != nil {
		return nil, err
	}

	// definir Slice para media
	var mediaList []Media

	// iterativo para obtener cada record individual de la consulta
	for rows.Next() {
		var id int
		var title string
		var description string
		var media_type_id int
		var source string
		var media_type string
		if err := rows.Scan(&id, &title, &description, &media_type_id, &source, &media_type); err != nil {
			return nil, err
		}

		// funcion append para añadir item a slice mediaList
		mediaList = append(mediaList, &movie{
			id:          id,
			title:       title,
			description: description,
			mediaType:   media_type,
			source:      source,
		})
	}

	return mediaList, nil
}

func (m *movie) Id() int {
	return m.id
}

func (m *movie) Title() string {
	return m.title
}

func (m *movie) Description() string {
	return m.description
}

func (m *movie) MediaType() string {
	return m.mediaType
}

// obtener url del objeto
func (m *movie) Source() string {
	return m.source
}

func (m *movie) Artist() string {
	return ""
}

func (m *movie) Album() string {
	return ""
}
