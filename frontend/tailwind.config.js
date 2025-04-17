/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'coffee': {
          100: '#F5EBE0',
          200: '#E6D7C3',
          300: '#C8B6A6',
          400: '#A4907C',
          500: '#8D7B68',  // Primary brand color
          600: '#735F4D',
          700: '#5A4D3D',
          800: '#4A3F32',
          900: '#332A21',
        },
        'cream': '#FFF8E1',
        'espresso': '#3A2618',
        'toasted': '#D4A762', // Your brand gold/amber
      }
    },
  },
  plugins: [],
}

