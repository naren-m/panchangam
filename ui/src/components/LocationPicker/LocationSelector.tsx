import React, { useState, useEffect } from 'react';
import { MapPin, Search, Navigation, X } from 'lucide-react';
import { Location } from '../../types/panchangam';
import { locationService } from '../../services/locationService';

interface LocationSelectorProps {
  currentLocation: Location;
  onLocationSelect: (location: Location) => void;
  onClose: () => void;
}

export const LocationSelector: React.FC<LocationSelectorProps> = ({
  currentLocation,
  onLocationSelect,
  onClose
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<Location[]>([]);
  const [popularLocations, setPopularLocations] = useState<Location[]>([]);
  const [loading, setLoading] = useState(false);
  const [gpsLoading, setGpsLoading] = useState(false);

  useEffect(() => {
    setPopularLocations(locationService.getPopularLocations());
  }, []);

  useEffect(() => {
    const searchLocations = async () => {
      if (searchQuery.trim().length < 2) {
        setSearchResults([]);
        return;
      }

      setLoading(true);
      try {
        const results = await locationService.searchLocations(searchQuery);
        setSearchResults(results);
      } catch (error) {
        console.error('Search error:', error);
      } finally {
        setLoading(false);
      }
    };

    const debounceTimer = setTimeout(searchLocations, 300);
    return () => clearTimeout(debounceTimer);
  }, [searchQuery]);

  const handleGpsLocation = async () => {
    setGpsLoading(true);
    try {
      const location = await locationService.getCurrentLocation();
      onLocationSelect(location);
      onClose();
    } catch (error) {
      console.error('GPS error:', error);
      alert('Unable to get your location. Please select manually.');
    } finally {
      setGpsLoading(false);
    }
  };

  const displayLocations = searchQuery.trim() ? searchResults : popularLocations;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-xl shadow-2xl max-w-md w-full max-h-[80vh] overflow-hidden">
        {/* Header */}
        <div className="bg-blue-500 text-white p-4">
          <div className="flex items-center justify-between">
            <h2 className="text-xl font-bold">Select Location</h2>
            <button
              onClick={onClose}
              className="p-1 hover:bg-blue-600 rounded-full transition-colors"
            >
              <X className="w-5 h-5" />
            </button>
          </div>
        </div>

        {/* Current location */}
        <div className="p-4 bg-blue-50 border-b">
          <div className="flex items-center space-x-2 text-blue-800">
            <MapPin className="w-4 h-4" />
            <span className="text-sm font-medium">Current: {currentLocation.name}</span>
          </div>
        </div>

        {/* Search */}
        <div className="p-4 border-b">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search for a city..."
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
            {loading && (
              <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-500"></div>
              </div>
            )}
          </div>
        </div>

        {/* GPS Button */}
        <div className="p-4 border-b">
          <button
            onClick={handleGpsLocation}
            disabled={gpsLoading}
            className="w-full flex items-center justify-center space-x-2 py-2 px-4 bg-green-500 text-white rounded-lg hover:bg-green-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <Navigation className="w-4 h-4" />
            <span>{gpsLoading ? 'Getting location...' : 'Use current location'}</span>
          </button>
        </div>

        {/* Locations list */}
        <div className="max-h-64 overflow-y-auto">
          <div className="p-4">
            <h3 className="text-sm font-semibold text-gray-600 mb-3">
              {searchQuery.trim() ? 'Search Results' : 'Popular Locations'}
            </h3>
            <div className="space-y-2">
              {displayLocations.map((location, index) => (
                <button
                  key={index}
                  onClick={() => {
                    onLocationSelect(location);
                    onClose();
                  }}
                  className="w-full text-left p-3 rounded-lg hover:bg-gray-100 transition-colors border border-gray-200"
                >
                  <div className="flex items-center space-x-3">
                    <MapPin className="w-4 h-4 text-gray-400" />
                    <div>
                      <div className="font-medium text-gray-800">{location.name}</div>
                      <div className="text-xs text-gray-500">
                        {location.latitude.toFixed(4)}, {location.longitude.toFixed(4)}
                      </div>
                    </div>
                  </div>
                </button>
              ))}
            </div>

            {displayLocations.length === 0 && searchQuery.trim() && !loading && (
              <div className="text-center text-gray-500 py-8">
                No locations found for "{searchQuery}"
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};