run:
	@templ generate
dev:
	@$(MAKE) -j 3 tailwind start-servers templ
tailwind:
	@tailwindcss -i ./web/assets/input.css -o ./web/assets/public/app.css -m --watch
templ:
	@templ generate --watch
database:
	@(cd ./microservices/database && go run api/main.go)
authentication:
	@(cd ./microservices/authentication && go run api/main.go)
passwordhashing:
	@(cd ./microservices/passwordhashing && go run api/main.go)
medialist:
	@(cd ./microservices/medialist && go run api/main.go)
movie:
	@(cd ./microservices/movie && go run api/main.go)
song:
	@(cd ./microservices/song && go run api/main.go)
registration:
	@(cd ./microservices/registration && go run api/main.go)
frontend:
	@(cd ./microservices/frontend && templ generate && go run api/main.go)
start-servers:
	@$(MAKE) -j 8 database authentication passwordhashing medialist movie song registration frontend
