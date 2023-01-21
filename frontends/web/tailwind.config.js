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
    "text-cyan-600",
    "text-green-500",
    "border-gray-400",
    "border-b",
    "program-definition"
  ],
  plugins: [],
}
