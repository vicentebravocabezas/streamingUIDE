package media

import (
	"log"

	"github.com/vicentebravocabezas/streamingUIDE/database"
)

type Media struct {
	Id          int
	Title       string
	Description string
	MediaType   string
	Source      string
}

// lista de multimedia registrada en la base de datos
func MediaList() ([]Media, error) {
	// ejecutar consulta SQL
	rows, err := database.OpenDB().Query("SELECT media.*, media_types.media_type FROM media LEFT JOIN media_types ON media.media_type_id = media_types.media_type_id")
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
		rows.Scan(&id, &title, &description, &media_type_id, &source, &media_type)

		// funcion append para añadir item a slice mediaList
		mediaList = append(mediaList, Media{
			Id:          id,
			Title:       title,
			Description: description,
			MediaType:   media_type,
			Source:      source,
		})
	}

	return mediaList, nil
}

// obtener url del objeto
func (m *Media) GetSource() string {
	return m.Source
}

// almacenar nuevo registro de media en la base de datos
func (m *Media) Store() {
	var mediaId int

	if err := database.OpenDB().QueryRow("SELECT media_type_id FROM media_types WHERE media_type = ?", m.MediaType).Scan(&mediaId); err != nil {
		log.Println(err)
	}

	database.OpenDB().Exec("INSERT INTO media VALUES (?, ?, ?, ?)", m.Title, m.Description, mediaId, m.Source)
}
