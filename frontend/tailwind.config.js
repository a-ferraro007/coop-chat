module.exports = {
  mode: 'jit',
  purge: ['./pages/**/*.{js,ts,jsx,tsx}', './components/**/*.{js,ts,jsx,tsx}'],
  darkMode: false, // or 'media' or 'class'
  theme: {
    fontFamily: {
      sans: 'Oswald, sans-serif'
    },
    extend: {
      colors: {
        'primary-dark': '#373F47',
        'login-bg-blue': 'rgba(0, 147, 233, .5)',
        'login-bg-green': 'rgba(128, 208, 99, .5)'
      },
      height: {
        400: '400px'
      },
      minHeight: {
        400: '400px'
      },
      width: {
        665: '565px'
      },
      borderRadius: {
        60: '60px'
      }
    }
  },
  variants: {
    extend: {}
  },
  plugins: []
}
