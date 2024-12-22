# Servicio Streaming UIDE

El proyecto consiste en un servicio de streaming de películas y canciones para usuarios registrados. El proyecto esta desarrollado en Go.

---

El alcance del proyecto comprende las siguientes funcionalidades programadas en los módulos en la carpeta `microservices` especificados a continuación:
- `authentication`: servicio Web para manejar autenticación al servidor
- `database`: servicio para manejar conexión a la base de datos. Todas las interacciones con la DB pasan por este servicio
- `frontend`: servicio que ofrece el frontend al usuario. Se comunica con los otros servicios a través de JSON pero devuelve HTML al usuario. Utiliza plantillas escritas en [templ](https://templ.guide/), un DSL creado para Go. El estilado fue realizado con [Tailwind CSS](https://tailwindcss.com/)
- `medialist`: servicio que obtiene la lista de las películas disponibles
- `movie`: servicio que obtiene los datos y la fuente de una película específica
- `passwordhashing`: servicio que se encarga de generar hashing de contraseñas y verificación de hashes para posterior autenticación
- `registration`: servicio para registrar usuarios nuevos a la base de datos. Se comunica con el servicio de `passwordhashing` para enviar un nuevo usuario a `database`
- `shared`: **este no es un servicio,** es un módulo con tipos, clases e implementación compartida por los servicios del sistema
- `song`: servicio que obtiene los datos y la fuente de una canción específica

La comunicación entre los servicios se realiza a través de JSON.

## Instalación

Para compilar el servicio. **Es necesario instalar [templ](https://templ.guide/)** y generar las plantillas primero.
La herramienta se instala con el siguiente código:

```sh
go install github.com/a-h/templ/cmd/templ@latest
```

La herramienta debe estar en el `PATH` del sistema.

## Compilación y uso

El repositorio provee un Makefile que automáticamente genera las plantillas, compila e inicia simultáneamente todos los servicios web con el comando:

```sh
make start-servers
```

Se podrá acceder al frontend con la dirección [localhost:8080](http://localhost:8080/)

---

Si se desea realizar la compilación de forma manual, se deben generar las plantillas con:
```sh
templ generate
```
y ejecutar cada modulo con:
```sh
go run api/main.go
```
desde la raíz de cada módulo de la carpeta `microservices`.

