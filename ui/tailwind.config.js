/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        'saffron': {
          50: '#FFF8E7',
          100: '#FFE6B8',
          200: '#FFD689',
          300: '#FFC65A',
          400: '#FFB62B',
          500: '#FF9933',
          600: '#E6852E',
          700: '#CC7029',
          800: '#B35C24',
          900: '#99471F',
        },
        'sacred-green': {
          50: '#E7F5E7',
          100: '#C2E4C2',
          200: '#9DD39D',
          300: '#78C278',
          400: '#53B153',
          500: '#138808',
          600: '#117A07',
          700: '#0F6C06',
          800: '#0D5E05',
          900: '#0B5004',
        },
        'divine-blue': {
          50: '#E7E7F5',
          100: '#C2C2E4',
          200: '#9D9DD3',
          300: '#7878C2',
          400: '#5353B1',
          500: '#000080',
          600: '#000073',
          700: '#000066',
          800: '#000059',
          900: '#00004C',
        },
        'cream': {
          50: '#FFF8DC',
          100: '#FFF2C7',
          200: '#FFECB3',
          300: '#FFE69E',
          400: '#FFE089',
          500: '#FFDA74',
        }
      },
      fontFamily: {
        'sans': ['Inter', 'Noto Sans Devanagari', 'system-ui', 'sans-serif'],
        'devanagari': ['Noto Sans Devanagari', 'serif'],
      },
      spacing: {
        '18': '4.5rem',
        '88': '22rem',
        '128': '32rem',
      },
      animation: {
        'pulse-slow': 'pulse 3s ease-in-out infinite',
        'fade-in': 'fadeIn 0.5s ease-in-out',
        'slide-up': 'slideUp 0.3s ease-out',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(20px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
      boxShadow: {
        'hindu': '0 10px 25px -5px rgba(255, 153, 51, 0.1), 0 10px 10px -5px rgba(255, 153, 51, 0.04)',
        'sacred': '0 20px 25px -5px rgba(19, 136, 8, 0.1), 0 10px 10px -5px rgba(19, 136, 8, 0.04)',
      },
      backgroundImage: {
        'gradient-hindu': 'linear-gradient(135deg, #FF9933 0%, #138808 50%, #000080 100%)',
        'gradient-sunrise': 'linear-gradient(90deg, #FF9933 0%, #FFC107 100%)',
        'gradient-sacred': 'linear-gradient(135deg, #138808 0%, #28A745 100%)',
      }
    },
  },
  plugins: [],
};