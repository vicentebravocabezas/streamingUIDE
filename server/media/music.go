package media

import "github.com/vicentebravocabezas/streamingUIDE/database"

type music struct {
	id        int
	title     string
	artist    string
	album     string
	mediaType string
	source    string
}

func GetMusic(id string) (Media, error) {
	row := database.OpenDB().QueryRow(
		`SELECT media.media_id, media.title, media.artist, media.album, media.media_type_id, media.source, media_types.media_type 
FROM media
LEFT JOIN media_types ON media.media_type_id = media_types.media_type_id
WHERE media.media_id = ?`, id,
	)

	var resultid int
	var title string
	var artist string
	var album string
	var media_type_id int
	var source string
	var media_type string
	if err := row.Scan(&resultid, &title, &artist, &album, &media_type_id, &source, &media_type); err != nil {
		return nil, err
	}

	return &music{
		id:        resultid,
		title:     title,
		artist:    artist,
		album:     album,
		mediaType: media_type,
		source:    source,
	}, nil
}

// lista de multimedia registrada en la base de datos
func MusicList() ([]Media, error) {
	// ejecutar consulta SQL
	rows, err := database.OpenDB().Query(
		`SELECT media.media_id, media.title, media.artist, media.album, media.media_type_id, media.source, media_types.media_type 
FROM media 
LEFT JOIN media_types ON media.media_type_id = media_types.media_type_id
WHERE media_types.media_type = 'Music'`,
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
		var artist string
		var album string
		var media_type_id int
		var source string
		var media_type string
		if err := rows.Scan(&id, &title, &artist, &album, &media_type_id, &source, &media_type); err != nil {
			return nil, err
		}

		// funcion append para añadir item a slice mediaList
		mediaList = append(mediaList, &music{
			id:        id,
			title:     title,
			artist:    artist,
			album:     album,
			mediaType: media_type,
			source:    source,
		})
	}

	return mediaList, nil
}

func (m *music) Id() int {
	return m.id
}

func (m *music) Title() string {
	return m.title
}

func (m *music) Description() string {
	return ""
}

func (m *music) MediaType() string {
	return m.mediaType
}

// obtener url del objeto
func (m *music) Source() string {
	return m.source
}

func (m *music) Artist() string {
	return m.artist
}

func (m *music) Album() string {
	return m.album
}
