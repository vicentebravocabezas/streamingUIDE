package frontend

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/frontend/cookies"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/frontend/web/templates"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/frontend/web/templates/stream"
	"github.com/vicentebravocabezas/streamingUIDE/microservices/shared"
)

func streamPage(c echo.Context) error {
	user, err := cookies.GetUser(c)
	if err != nil {
		return err
	}

	mediaId := c.QueryParam("id")
	mediaType := c.QueryParam("media-type")

	var component templ.Component

	if mediaId == "" {
		var mediaList []shared.MediaGlobalType
		resp, err := http.Get(shared.MediaListAddr.WithSchemeAndPath(""))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if err := json.Unmarshal(body, &mediaList); err != nil {
			return err
		}

		component = stream.MediaList(mediaList)
	}

	if mediaId != "" {
		var media shared.MediaGlobalType
		if mediaType == "Movie" {
			resp, err := http.Get(shared.MovieAddr.WithSchemeAndPath("/" + mediaId))
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			if err := json.Unmarshal(body, &media); err != nil {
				return err
			}
		} else if mediaType == "Music" {
			resp, err := http.Get(shared.SongAddr.WithSchemeAndPath("/" + mediaId))
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			if err := json.Unmarshal(body, &media); err != nil {
				return err
			}
		}

		component = stream.Play(media)
	}

	return render(c, http.StatusOK, templates.Layout(stream.Stream(user.GetUsername(), component)))
}
