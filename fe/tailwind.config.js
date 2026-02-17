/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Plus Jakarta Sans', 'Inter', 'sans-serif'],
      },
      boxShadow: {
        'soft': '0 20px 40px -10px rgba(0,0,0,0.05)',
        'glow': '0 0 20px rgba(255, 71, 87, 0.3)',
      },
    },
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: [
      {
        foodie: {
          "primary": "#FF4757", // Vibrant Coral Red - Appetizing & Energetic
          "primary-content": "#ffffff",
          "secondary": "#2F3542", // Deep Charcoal - Sophisticated
          "secondary-content": "#ffffff",
          "accent": "#FFA502", // Golden Yellow - Cheese/Crust
          "accent-content": "#ffffff",
          "neutral": "#57606f", // Soft Slate
          "neutral-content": "#ffffff",
          "base-100": "#ffffff", // Pure White Cards
          "base-200": "#F1F2F6", // Light Gray Background
          "base-300": "#DFE4EA", // Borders/Dividers
          "info": "#70a1ff",
          "success": "#2ed573",
          "warning": "#eccc68",
          "error": "#ff4757",
          "--rounded-box": "1.5rem",
          "--rounded-btn": "1rem",
        },
      },
    ],
  },
}
