module.exports = {
	proxy: "http://localhost:8080",
	files: ["html/**/*.templ", "public/**/*.*", "main.go"],
	// notify: false,
	open: false,
	ui: false,
	reloadDelay: 2500,
};
