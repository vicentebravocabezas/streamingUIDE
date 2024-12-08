build:
	@CGO_ENABLED=0 go build -o main . && upx main && chmod +x main
dev:
	@$(MAKE) -j 3 tailwind run templ
tailwind:
	@tailwindcss -i ./web/assets/input.css -o ./web/assets/public/app.css -m --watch
templ:
	@templ generate --watch
run:
	@air
