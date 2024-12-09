package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/vicentebravocabezas/streamingUIDE/server/authentication"
	"github.com/vicentebravocabezas/streamingUIDE/server/media"
	"github.com/vicentebravocabezas/streamingUIDE/web/templates"
	"github.com/vicentebravocabezas/streamingUIDE/web/templates/stream"
)

func streamPage(c echo.Context) error {
	user, err := authentication.ReadUserFromCookies(c)
	if err != nil {
		return err
	}

	mediaId := c.QueryParam("id")
	mediaType := c.QueryParam("media-type")

	var component templ.Component

	if mediaId == "" {
		mediaList, err := media.MediaList()
		if err != nil {
			return err
		}

		component = stream.MediaList(mediaList)
	}

	if mediaId != "" {
		media, err := media.GetMedia(mediaId, mediaType)
		if err != nil {
			return err
		}

		component = stream.Play(media)
	}

	return render(c, http.StatusOK, templates.Layout(stream.Stream(user.Username(), component)))
}
