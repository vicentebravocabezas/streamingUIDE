package server

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo) {
	e.GET("/", index)

	// ruta para insertar usuario
	e.GET("/insert-user", insertUser)

	//ruta para borrar usuario
	e.GET("/delete-user", deleteUser)

	//ruta para obtener lista de archivos multimedia
	e.GET("/media", mediaList)
}
