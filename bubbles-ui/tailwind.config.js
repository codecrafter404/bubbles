/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        "primary": {
          "100": "#F2FCFF",
          "200": "#B7F3FE",
          "300": "#7AEAF7",
          "400": "#3DDCE3",
          "500": "#08BDBA",
          "600": "#03978C",
          "700": "#017162",
          "800": "#004C3D",
          "900": "#00261D"
        },
        "accent": {
          "100": "#ffe4e1",
          "200": "#ffcec9",
          "300": "#feaca3",
          "400": "#fb7c6e",
          "500": "#f35240",
          "600": "#db321f",
          "700": "#bc2919",
          "800": "#9c2518",
          "900": "#81251b"
        },
        "natural": {
          "100": "#FAFCFC",
          "200": "#E6EAEA",
          "300": "#D2D7D7",
          "400": "#BFC5C5",
          "500": "#ABB3B2",
          "600": "#888F8E",
          "700": "#656C6B",
          "800": "#434947",
          "900": "#222625"
        }
      }
    },
  },
  plugins: [],
}

