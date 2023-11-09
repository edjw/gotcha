/** @type {import('tailwindcss').Config} */
module.exports = {
	content: [
		"./html/components/**/*.templ",
		"./html/pages/**/*.templ",
		"./html/layouts/**/*.templ",
		"./html/partials/**/*.templ",
	],
	plugins: [require("@tailwindcss/typography")],
};
