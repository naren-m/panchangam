# Panchangam UI

A modern React TypeScript application for displaying Hindu calendar (Panchangam) data with real-time astronomical calculations.

## Features

- **Real-time Panchangam Data**: Connected to the real gRPC service via HTTP gateway
- **5 Core Elements**: Tithi, Nakshatra, Yoga, Karana, and Vara calculations
- **Location Support**: GPS location detection and manual location selection
- **Responsive Design**: Works on desktop and mobile devices
- **Health Monitoring**: API connection status and debugging tools
- **Settings Panel**: Customizable calculation methods, language, and regional preferences

## Architecture

```
UI (React/TypeScript) → HTTP Gateway (Go) → gRPC Service (Go) → Astronomy Calculations
```

## Prerequisites

- Node.js 18+ and npm
- Go 1.21+ (for backend services)
- The Panchangam gRPC service and HTTP gateway running

## Quick Start

### 1. Start Backend Services

First, start the gRPC service:
```bash
# From project root
go run main.go
```

Then start the HTTP gateway:
```bash
# From project root
go run cmd/gateway/main.go
```

The services will be available at:
- gRPC Service: `localhost:50051`
- HTTP Gateway: `http://localhost:8080`
- Health Check: `http://localhost:8080/api/v1/health`

### 2. Start Frontend

```bash
# Install dependencies
npm install

# Start development server
npm run dev
```

The UI will be available at: `http://localhost:5173`

## Environment Variables

Create `.env.development` for local development:

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_API_TIMEOUT=30000
VITE_DEBUG_API=true
```

For production, create `.env.production`:

```env
VITE_API_BASE_URL=https://api.panchangam.app
VITE_API_TIMEOUT=30000
VITE_DEBUG_API=false
```

## API Integration

The UI connects to the real Panchangam service through the HTTP gateway:

### API Endpoint
```
GET /api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946&tz=Asia/Kolkata
```

### Response Format
```json
{
  "date": "2024-01-15",
  "tithi": "Chaturthi (4)",
  "nakshatra": "Uttara Bhadrapada (26)",
  "yoga": "Siddha (21)",
  "karana": "Gara (6)",
  "sunrise_time": "01:15:32",
  "sunset_time": "12:41:47",
  "events": [
    {
      "name": "Tithi: Chaturthi",
      "time": "01:15:32",
      "event_type": "TITHI"
    }
  ]
}
```

## Development

### Project Structure

```
src/
├── components/          # React components
│   ├── Calendar/       # Calendar grid and navigation
│   ├── DayDetail/      # Day details modal and events
│   ├── LocationPicker/ # Location selection
│   └── Settings/       # Settings panel and API health check
├── hooks/              # Custom React hooks
├── services/           # API services and location handling
├── types/              # TypeScript type definitions
└── utils/              # Utility functions
```

### Key Components

- **ApiHealthCheck**: Monitors API connection status
- **PanchangamApi**: Service layer for API communication
- **usePanchangam**: Hook for fetching Panchangam data
- **LocationSelector**: GPS and manual location selection

### Error Handling

The app includes comprehensive error handling:

1. **API Connection Issues**: Falls back to placeholder data
2. **Network Timeouts**: 30-second timeout with retry logic
3. **Invalid Responses**: Graceful error messages
4. **Location Errors**: Falls back to popular locations

### Health Monitoring

The app includes a built-in health check system:

- Real-time API status monitoring
- Connection debugging information
- Development mode helpers
- Automatic fallback when API is unavailable

## Building

### Development Build
```bash
npm run build:dev
```

### Production Build
```bash
npm run build
```

### Preview Production Build
```bash
npm run preview
```

## Testing

```bash
# Run tests
npm test

# Run tests with coverage
npm run test:coverage

# Run E2E tests
npm run test:e2e
```

## Deployment

### Environment Setup

1. **Development**: Uses `localhost:8080` for API
2. **Production**: Configure `VITE_API_BASE_URL` to your production API

### Build and Deploy

```bash
# Build for production
npm run build

# Deploy dist/ folder to your hosting service
```

## Troubleshooting

### API Connection Issues

1. **Check API Health**: Use the Settings panel to view API status
2. **Verify Backend**: Ensure gRPC service and gateway are running
3. **Check Network**: Verify no firewall blocking localhost:8080
4. **Review Logs**: Check browser console for detailed error messages

### Common Issues

- **"API Unavailable"**: Backend services not running
- **"Connection Failed"**: Network or firewall issues
- **"Invalid Response"**: API version mismatch or malformed data

### Development Mode

In development mode, the app shows additional debugging information:

- API endpoint details
- Connection status
- Instructions for starting backend services
- Detailed error messages

## Contributing

1. Follow the existing code style
2. Add proper TypeScript types
3. Include error handling
4. Test API integration thoroughly
5. Update documentation as needed

## License

This project is part of the Panchangam system. See the main project LICENSE file.