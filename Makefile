deploy_render:
	export GOPATH=/opt/render/project/go && \
	export PATH=$$PATH:$$GOPATH/bin && \
	go install github.com/a-h/templ/cmd/templ@v0.2.408 && \
	$$GOPATH/bin/templ generate && \
	pnpm install --production && \
	pnpm exec tailwindcss -i ./assets/global.css -o ./public/global.css --minify && \
	go build -tags netgo -ldflags '-s -w' -o app
