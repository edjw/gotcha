## Gotcha

A stack for web development

- Go
- Templ
- Tailwind
- Chi
- Htmx
- Air
- Browser sync


## Run locally

Install tailwind with `pnpm install`

If you're developing on macOS, set `export GO_ENV=development` in your `.zshrc` or `.bashrc` file. This is to allow "FriendlyServer" to set the IP address of your development server to `127.0.0.1` which stops macOS asking you to allow your app through the firewall on every single rebuild of your app.

Run `air` to start in development. If you want live-reload of the browser, run `browser-sync start --config browser-sync-config.js` in a separate terminal tab. It has to wait for the things in air to finish before reloading. You can also just manually refresh your browser.

## Deploy

Only deploying to render.com so far. Use `make deploy_render` as your build command on Render.


## Future

Show how to deploy to fly (gh action), railway

Maybe make it into a cli like t3 but make the cli in cobra-cli