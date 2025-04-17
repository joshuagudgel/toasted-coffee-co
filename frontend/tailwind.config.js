/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'terracotta': '#BF7454',
        'espresso': '#594A47',
        'peach': '#DD9D79',
        'parchment': '#F5F5EF',
        'latte': '#FDE4CD',
        'caramel': '#FDD39D',
        'mocha': '#945F48',
        'sage': '#D7E6E0'
      }
    },
  },
  plugins: [],
}

