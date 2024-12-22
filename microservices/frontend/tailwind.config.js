/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ["./web/templates/**/*.{html,js,templ}"],
	theme: {
		extend: {},
	},
	plugins: [
		require('@tailwindcss/forms'),
		require('@tailwindcss/typography'),
	]
}

