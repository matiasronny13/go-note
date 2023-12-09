/** @type {import('tailwindcss').Config} */
module.exports = {
  //content: ["./internal/handler/web/template/**/*.{html,js}"],
  content: ["./web/template/**/*.html"],
  theme: {    
    colors: {
      'base': {
        '50': '#f6f6f9',
        '100': '#ebebf3',
        '200': '#d3d4e4',
        '300': '#adafcc',
        '400': '#8084b0',
        '500': '#606597',
        '600': '#4c4f7d',
        '700': '#3e4066',
        '800': '#363856',
        '900': '#313249',
        '950': '#0e0e15',
      }
    },
    textColor: {
      'base': {
        '50': '#f5f4ef',
        '100': '#e8e6d9',
        '200': '#d2cfb6',
        '300': '#b9b18d',
        '400': '#a3956c',
        '500': '#92825d',
        '600': '#7d6a4f',
        '700': '#665442',
        '800': '#59473b',
        '900': '#4e3f36',
        '950': '#15100e',
      },
      'cinnabar': {
        '300': '#f8a9a9',
        '400': '#f27777',
        '500': '#e63d3d',
        '600': '#d42e2e',
        '700': '#b22323',
        '800': '#942020',
        '900': '#7b2121',
      }
    },
    extend: {
      typography: ({ theme }) => ({
        DEFAULT: {
          css: {
            '--tw-prose-body': theme('textColor.base[100]'),
            '--tw-prose-headings': theme('textColor.base[100]'),
            h1: {
              margin: "0px"
            },
            p: {
              fontSize: '16px',
              lineHeight: "1.75",
              wordBreak: "break-all",
              whiteSpace: "normal"
            },
            ul: {
              fontSize: '16px',
              lineHeight: "1.5"
            },
            a: {
              color: theme('textColor.cinnabar[400]'),
              '&:hover': {
                color: '#007bff',
              },
            },
          },
        },
      }),
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}

