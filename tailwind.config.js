/** @type {import('tailwindcss').Config} */
module.exports = {
	content: [
		"./html/components/**/*.templ",
		"./html/pages/**/*.templ",
		"./html/layouts/**/*.templ",
		"./html/partials/**/*.templ",
	],
	theme: {
		extend: {},
	},
	plugins: [require("@tailwindcss/typography")],
};
