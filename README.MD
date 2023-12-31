## Gotcha

A stack for web development using the following projects

- **Go**
- **T**empl for templating
- **T**ailwind for CSS
- **C**hi for routing
- **H**TMX for fetching server state as HTML and updating the DOM
- **A**lpine for clientside Javascript in your HTML

---

- Air for live-reloading of Go app
- Browser sync for live-reloading of the browser

Includes deployment setups for Render and Fly.io and a Dockerfile. Successfully deployed to Railway and AWS via Flightcontrol with the Dockerfile.

This is a learning project for me to learn Go. I'm not a Go developer so I'm sure there are better ways to do things.

## Run locally

Install Tailwind with `pnpm install`

If you're developing on macOS, set `export GO_ENV=development` in your `.zshrc` or `.bashrc` file. This is to allow "FriendlyServer" to set the IP address of your development server to `127.0.0.1`. This stops macOS asking you to allow your app through the firewall on every single rebuild of your app.

Run `air` to start in development.

If you want live-reload of the browser, run `browser-sync start --config browser-sync-config.js` in a separate terminal tab. It has to wait for the things in `air` to finish before reloading so it is not a very fast automatic live-reload. You can also just manually refresh your browser.

Optionally, you can set an environment variable called `DEPLOYMENT_SITE_URL` and give it the full URL of your site eg `https://gotcha-zain.onrender.com`. If you do, this will try to prevent requests to your partials/fragments of HTML from other sites.

## Deploy

### Render
- Use `make deploy_render` as your build command on <https://render.com>
- Some of the commands in the Makefile for deploying to Render make my brain hurt.

### Fly.io
- Use the included `fly.toml`, `Dockerfile` and Github action to deploy to <https://fly.io>

- The Github action is optional but it is nice to have a build on every push to `main` which Fly doesn't support natively.

### Railway

- Works with just the Dockerfile

### AWS via Flightcontrol

- Works with just the Dockerfile


### Digital Ocean Apps Platform

- Works with just the Dockerfile. Make sure you disable the other auto-detected 'service' and just use the Dockerfile.


## Future

- Make it into a cli where you can choose what parts you want like [create-t3-app](https://create.t3.gg) but make the cli in cobra-cli. Mainly to practise cobra-cli.