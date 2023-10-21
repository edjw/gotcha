## Gotcha

A stack for web development

- Go ✅
- Templ ✅
- Tailwind
- Chi ✅
- Htmx
- Air
- Browser sync

No repeated firewall issue on macOS ✅



## Run

If you're developing on macOS, set `export GO_ENV=development` in your `.zshrc` or `.bashrc` file. This is to allow "FriendlyServer" to set the IP address of your development server to `127.0.0.1` which stops macOS asking you to allow your app through the firewall on every single rebuild of your app.

Run `air` to start in development


## Future

Show how to deploy to fly (gh action), render, railway

Maybe make it into a cli like t3 but make the cli in cobra-cli