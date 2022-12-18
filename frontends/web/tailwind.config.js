/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{html,js,ts,jsx,tsx}",
    "./index.html"
  ],
  theme: {
    extend: {},
  },
  safelist : [
    "good-error",
    "bad-error",
  ],
  plugins: [],
}
