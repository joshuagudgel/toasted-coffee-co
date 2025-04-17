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
        'sage': '#D7E6E0',
        
        // Add shades for gradient and variation use
        'coffee': {
          100: '#F5F5EF', // parchment
          200: '#FDE4CD', // latte
          300: '#FDD39D', // caramel
          400: '#DD9D79', // peach
          500: '#BF7454', // terracotta (primary brand)
          600: '#945F48', // mocha
          700: '#594A47', // espresso
          800: '#3A2618', // dark espresso (keeping this from original)
        },
        
        // Keep the toasted color for backward compatibility
        'toasted': '#FDD39D', // updated to match your new caramel color
      }
    },
  },
  plugins: [],
}

