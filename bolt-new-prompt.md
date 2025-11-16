# Panchangam Calendar Web Application - bolt.new Design Prompt

## Project Overview
Create a modern, responsive web application for displaying Panchangam (Hindu calendar) data with daily astronomical calculations, auspicious times, and traditional Hindu temporal information.

## Core Application Structure

### Backend API Integration
Your app should integrate with a Go-based gRPC API with the following structure:

**Primary Endpoint**: Panchangam service with `Get(GetPanchangamRequest)` method
**Request Parameters**:
```json
{
  "date": "2024-01-15",           // ISO 8601 format (YYYY-MM-DD)
  "latitude": 13.0827,            // Location coordinates
  "longitude": 80.2707,
  "timezone": "Asia/Kolkata",     // Timezone identifier
  "region": "Tamil Nadu",         // Regional calculation system
  "calculation_method": "Drik",   // Drik (modern) or Vakya (traditional)
  "locale": "en"                  // Language (en, ta, hi, etc.)
}
```

**Response Structure**:
```json
{
  "panchangam_data": {
    "date": "2024-01-15",
    "tithi": "Shukla Panchami",        // Lunar day (1-30)
    "nakshatra": "Rohini",             // Lunar mansion (1-27)
    "yoga": "Siddha",                  // Combination of sun-moon (1-27)
    "karana": "Bava",                  // Half-tithi period (1-11)
    "sunrise_time": "06:45:30",        // Local sunrise time
    "sunset_time": "18:15:45",         // Local sunset time
    "events": [
      {
        "name": "Rahu Kalam",
        "time": "09:30:00-11:00:00",
        "event_type": "RAHU_KALAM"
      },
      {
        "name": "Abhijit Muhurta",
        "time": "12:00:00-12:48:00",
        "event_type": "MUHURTA"
      }
    ]
  }
}
```

## UI Design Requirements

### 1. Calendar Layout
- **Primary View**: Monthly calendar grid similar to Google Calendar
- **Date Cells**: Each day displays compact panchangam information
- **Navigation**: Month/Year navigation with today button
- **Responsive**: Mobile-first design with tablet/desktop adaptations

### 2. Daily Cell Content (Compact View)
Each calendar day should display:
```
[Date Number]
𑚕 Tithi Name
⭐ Nakshatra 
🧘 Yoga
⚡ Karana
🌅 06:45 🌅 18:15
```

### 3. Detailed Day View (Expandable)
When clicking a date, show expanded information:

#### The Five Angas (पञ्चाङ्ग)
- **Tithi (तिथि)**: Lunar day with completion percentage
- **Nakshatra (नक्षत्र)**: Current star constellation with deity and characteristics  
- **Yoga (योग)**: Auspicious combination with quality rating
- **Karana (करण)**: Half-tithi period with attributes
- **Vara (वार)**: Weekday with planetary ruler and significance

#### Astronomical Times
- **Sunrise/Sunset**: Precise local times with civil twilight
- **Moonrise/Moonset**: If available from backend
- **Solar Noon**: Midday time for muhurta calculations

#### Muhurta (Auspicious Times)
Display as timeline with color coding:
- **Green**: Highly auspicious (Abhijit, Brahma Muhurta)
- **Yellow**: Mildly auspicious 
- **Orange**: Neutral periods
- **Red**: Inauspicious (Rahu Kalam, Yamagandam, Gulika Kalam)

#### Events & Festivals
- Religious festivals and observances
- Regional celebrations
- Ekadashi, Amavasya, Purnima markers
- Fasting days and special occasions

### 4. Visual Design Elements

#### Color Scheme
```css
:root {
  --primary-saffron: #FF9933;
  --secondary-green: #138808;
  --accent-blue: #000080;
  --background-cream: #FFF8DC;
  --text-dark: #2C3E50;
  --auspicious-green: #28A745;
  --neutral-yellow: #FFC107;
  --inauspicious-red: #DC3545;
}
```

#### Typography
- **Headers**: Devanagari-friendly fonts (Noto Sans Devanagari, Sanskrit fonts)
- **Body**: Clean, readable fonts (Inter, Roboto) with good multilingual support
- **Sanskrit/Tamil**: Appropriate font rendering for regional scripts

#### Icons & Symbols
- Traditional Hindu symbols (ॐ, 🕉️, ⚛️)
- Astronomical symbols (☀️, 🌙, ⭐, 🌅, 🌄)
- Time indicators (⏰, ⏳, 🕐)

### 5. Interactive Features

#### Location Services
- GPS-based location detection
- Manual location search with city names
- Location favorites and recent locations
- Coordinate display for verification

#### Settings Panel
- **Calculation Method**: Toggle between Drik (modern astronomical) and Vakya (traditional)
- **Language**: Support for English, Hindi, Tamil, Sanskrit
- **Region**: North India, South India, specific states
- **Display Options**: 12/24 hour format, temperature units
- **Notifications**: Daily panchangam, festival alerts

#### Export & Sharing
- Share daily panchangam via social media
- Export monthly calendar as PDF
- Generate panchangam for specific date ranges
- Print-friendly layouts

## Technical Implementation Guide

### Frontend Technology Stack
**Recommended**: React/Next.js with TypeScript
**Styling**: Tailwind CSS or styled-components
**State Management**: Zustand or Redux Toolkit
**HTTP Client**: Axios or fetch with proper error handling
**Date Handling**: date-fns or Day.js for timezone-aware operations

### Backend Integration
The application integrates with the Go-based gRPC API through HTTP endpoints. The API client handles:
- Request validation (date format, coordinates, timezone)
- Response validation and transformation
- Caching for performance optimization
- Error handling with graceful degradation
- Progressive loading for date ranges

See `ui/src/services/api/panchangamApiClient.ts` for the complete implementation.

### Component Architecture
```
src/
├── components/
│   ├── Calendar/
│   │   ├── CalendarGrid.tsx
│   │   ├── DateCell.tsx
│   │   └── MonthNavigation.tsx
│   ├── DayDetail/
│   │   ├── DayDetailModal.tsx
│   │   ├── FiveAngas.tsx
│   │   ├── MuhurtaTimeline.tsx
│   │   └── EventsList.tsx
│   ├── LocationPicker/
│   │   └── LocationSelector.tsx
│   └── Settings/
│       └── SettingsPanel.tsx
├── services/
│   ├── panchangamApi.ts
│   └── locationService.ts
├── types/
│   └── panchangam.ts
└── utils/
    ├── dateHelpers.ts
    └── timeFormatters.ts
```

## Data Definitions & Context

### The Five Angas Explained
1. **Tithi**: 30 lunar days per lunar month, each ~19-26 hours
2. **Nakshatra**: 27 lunar mansions, each 13°20' of ecliptic
3. **Yoga**: 27 combinations of sun-moon positions
4. **Karana**: 11 half-tithis cycling through lunar month
5. **Vara**: 7 weekdays with planetary associations

### Important Muhurtas
- **Brahma Muhurta**: 96 minutes before sunrise (spiritual practices)
- **Abhijit Muhurta**: 24 min before/after noon (universally auspicious)
- **Godhuli Muhurta**: Twilight period (cow-dust time)
- **Rahu Kalam**: Inauspicious period (varies by weekday)
- **Yamagandam**: Another inauspicious period
- **Gulika Kalam**: Third inauspicious time

### Regional Variations
- **North India**: Follows traditional Purnimanta (full moon ending) months
- **South India**: Follows Amanta (new moon ending) months  
- **Tamil Nadu**: Unique Tamil calendar integration with solar months
- **Kerala**: Malayalam calendar with specific local traditions

## Accessibility & Performance

### Accessibility
- WCAG 2.1 AA compliance
- Screen reader support for Sanskrit/Hindi text
- High contrast mode for visually impaired
- Keyboard navigation for all interactive elements
- Alt text for symbolic representations

### Performance  
- Lazy loading for calendar months
- Service worker for offline panchangam data
- Optimized bundle splitting
- Progressive Web App (PWA) capabilities
- Efficient re-rendering with proper React patterns

## Testing Requirements
- Unit tests for date calculations and formatting
- Integration tests for API calls and data transformation
- Visual regression tests for calendar layouts
- Cross-browser testing (Chrome, Firefox, Safari, Edge)
- Mobile responsiveness testing

## Deployment Considerations
- Environment variables for API endpoints
- CDN integration for static assets  
- SEO optimization with proper meta tags
- Social media preview cards
- Error boundary components for graceful failures

## Success Metrics
- Clean, intuitive calendar interface
- Accurate panchangam data display
- Responsive across all device sizes
- Fast loading times (<3 seconds)
- Accessible to users with disabilities
- Culturally sensitive and authentic presentation

Build this as a production-ready application that respects the sacred and traditional nature of panchangam while providing a modern, user-friendly experience.