module.exports = {
	proxy: "http://localhost:8080",
	files: ["components/**/*.templ", "public/**/*.*", "main.go"],
	notify: false,
	open: false,
	ui: false,
	reloadDelay: 4000,
};
