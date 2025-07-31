import React, { useState, useMemo, useRef } from 'react';

function App() {
  const [currentDate, setCurrentDate] = useState(new Date());
  const [settingsState, setSettingsState] = useState({
    calculation_method: 'Drik',
    locale: 'en',
    region: 'California',
    time_format: '12',
    location: {
      name: 'Milpitas, California',
      latitude: 37.4323,
      longitude: -121.9066,
      timezone: 'America/Los_Angeles',
      region: 'California'
    }
  });

  console.log('🔄 App render:', new Date().toISOString());

  return (
    <div className="min-h-screen bg-gradient-to-br from-orange-50 to-yellow-50">
      <div className="container mx-auto px-4 py-6 max-w-7xl">
        <div className="text-center mb-8">
          <h1 className="text-4xl md:text-5xl font-bold text-orange-800 mb-2">
            🕉️ Panchangam
          </h1>
          <p className="text-orange-600 text-lg">
            Hindu Calendar & Astronomical Almanac
          </p>
        </div>
        
        <div className="bg-white rounded-lg shadow-lg p-4">
          <p>Minimal test version - checking for infinite render loop</p>
          <p>Current time: {new Date().toLocaleTimeString()}</p>
        </div>
      </div>
    </div>
  );
}

export default App;