/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'upgrade': '#10b981',
        'downgrade': '#ef4444', 
        'reiterate': '#6b7280',
      }
    },
  },
  plugins: [],
}

