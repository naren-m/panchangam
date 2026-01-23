import React, { useState, useEffect } from 'react';
import { MapPin, Search, Navigation, X, Clock, Heart, Star } from 'lucide-react';
import { Location } from '../../types/panchangam';
import { locationService } from '../../services/locationService';

interface LocationSelectorProps {
  currentLocation: Location;
  onLocationSelect: (location: Location) => void;
  onClose: () => void;
}

interface LocationItemProps {
  location: Location;
  onSelect: () => void;
  showRegion?: boolean;
  isFavorite?: boolean;
  onToggleFavorite?: (location: Location) => void;
}

const LocationItem: React.FC<LocationItemProps> = ({ 
  location, 
  onSelect, 
  showRegion = false, 
  isFavorite = false,
  onToggleFavorite 
}) => {
  const getTimezoneDisplay = (timezone: string) => {
    // Convert timezone to a more readable format
    if (timezone.startsWith('America/')) {
      const city = timezone.split('/')[1].replace('_', ' ');
      return city;
    }
    if (timezone.startsWith('Asia/')) {
      return timezone.split('/')[1];
    }
    if (timezone.startsWith('Europe/')) {
      return timezone.split('/')[1];
    }
    if (timezone.startsWith('Pacific/')) {
      return timezone.split('/')[1];
    }
    return timezone;
  };

  return (
    <div className="flex items-center space-x-2">
      <button
        onClick={onSelect}
        className="flex-1 text-left p-3 rounded-lg hover:bg-gray-100 transition-colors border border-gray-200 hover:border-blue-300"
      >
        <div className="flex items-center space-x-3">
          <MapPin className="w-4 h-4 text-gray-400 flex-shrink-0" />
          <div className="flex-1 min-w-0">
            <div className="font-medium text-gray-800 truncate">{location.name}</div>
            <div className="flex items-center space-x-2 text-xs text-gray-500">
              <span>{location.latitude.toFixed(4)}, {location.longitude.toFixed(4)}</span>
              {showRegion && location.region && (
                <>
                  <span>•</span>
                  <span className="truncate">{location.region}</span>
                </>
              )}
            </div>
            <div className="flex items-center space-x-1 text-xs text-blue-600 mt-1">
              <Clock className="w-3 h-3" />
              <span>{getTimezoneDisplay(location.timezone)}</span>
            </div>
          </div>
        </div>
      </button>
      
      {onToggleFavorite && (
        <button
          onClick={(e) => {
            e.stopPropagation();
            onToggleFavorite(location);
          }}
          className={`p-2 rounded-full transition-colors ${
            isFavorite 
              ? 'bg-red-100 text-red-600 hover:bg-red-200' 
              : 'bg-gray-100 text-gray-400 hover:bg-gray-200 hover:text-red-500'
          }`}
          title={isFavorite ? 'Remove from favorites' : 'Add to favorites'}
        >
          <Heart className={`w-4 h-4 ${isFavorite ? 'fill-current' : ''}`} />
        </button>
      )}
    </div>
  );
};

export const LocationSelector: React.FC<LocationSelectorProps> = ({
  currentLocation,
  onLocationSelect,
  onClose
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<Location[]>([]);
  const [locationsByCategory, setLocationsByCategory] = useState<{
    favorites: Location[];
    usUk: Location[];
    popular: Location[];
  }>({ favorites: [], usUk: [], popular: [] });
  const [loading, setLoading] = useState(false);
  const [gpsLoading, setGpsLoading] = useState(false);

  useEffect(() => {
    setLocationsByCategory(locationService.getLocationsByCategory());
  }, []);

  const handleToggleFavorite = (location: Location) => {
    const isFav = locationService.isFavorite(location);
    if (isFav) {
      locationService.removeFromFavorites(location);
    } else {
      locationService.addToFavorites(location);
    }
    // Refresh categories
    setLocationsByCategory(locationService.getLocationsByCategory());
  };

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

  const displayLocations = searchQuery.trim() ? searchResults : 
    [...locationsByCategory.favorites, ...locationsByCategory.usUk, ...locationsByCategory.popular];

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
              placeholder="Search cities worldwide..."
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
        <div className="max-h-80 overflow-y-auto">
          <div className="p-4">
            <h3 className="text-sm font-semibold text-gray-600 mb-3">
              {searchQuery.trim() ? `Search Results for "${searchQuery.trim()}"` : 'Popular Locations'}
            </h3>
            
            {/* Categorized locations */}
            {!searchQuery.trim() && (
              <div className="mb-4">
                {/* Favorites */}
                {locationsByCategory.favorites.length > 0 && (
                  <div className="mb-6">
                    <div className="flex items-center space-x-2 text-xs text-pink-600 font-medium mb-2">
                      <Star className="w-3 h-3" />
                      <span>Your Favorites</span>
                    </div>
                    <div className="space-y-2">
                      {locationsByCategory.favorites.map((location, index) => (
                        <LocationItem 
                          key={`favorite-${index}`} 
                          location={location} 
                          onSelect={() => {
                            onLocationSelect(location);
                            onClose();
                          }}
                          isFavorite={true}
                          onToggleFavorite={handleToggleFavorite}
                        />
                      ))}
                    </div>
                  </div>
                )}
                
                {/* US/UK Locations */}
                <div className="mb-6">
                  <div className="text-xs text-blue-600 font-medium mb-2">🇺🇸🇬🇧 US & UK Cities</div>
                  <div className="space-y-2">
                    {locationsByCategory.usUk.map((location, index) => (
                      <LocationItem 
                        key={`us-uk-${index}`} 
                        location={location} 
                        onSelect={() => {
                          onLocationSelect(location);
                          onClose();
                        }}
                        isFavorite={locationService.isFavorite(location)}
                        onToggleFavorite={handleToggleFavorite}
                      />
                    ))}
                  </div>
                </div>
                
                {/* Popular (mainly Indian) Locations */}
                {locationsByCategory.popular.length > 0 && (
                  <div>
                    <div className="text-xs text-blue-600 font-medium mb-2">🇮🇳 Other Popular Cities</div>
                    <div className="space-y-2">
                      {locationsByCategory.popular.slice(0, 15).map((location, index) => (
                        <LocationItem 
                          key={`popular-${index}`} 
                          location={location} 
                          onSelect={() => {
                            onLocationSelect(location);
                            onClose();
                          }}
                          isFavorite={locationService.isFavorite(location)}
                          onToggleFavorite={handleToggleFavorite}
                        />
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}
            
            {/* Search results */}
            {searchQuery.trim() && (
              <div className="space-y-2">
                {searchResults.map((location, index) => (
                  <LocationItem 
                    key={`search-${index}`} 
                    location={location} 
                    onSelect={() => {
                      onLocationSelect(location);
                      onClose();
                    }}
                    showRegion={true}
                    isFavorite={locationService.isFavorite(location)}
                    onToggleFavorite={handleToggleFavorite}
                  />
                ))}
              </div>
            )}

            {searchResults.length === 0 && searchQuery.trim() && !loading && (
              <div className="text-center text-gray-500 py-8">
                <MapPin className="w-8 h-8 text-gray-300 mx-auto mb-2" />
                <div className="text-sm font-medium">No locations found</div>
                <div className="text-xs">Try searching for a city, state, or country</div>
              </div>
            )}

            {loading && (
              <div className="text-center text-gray-500 py-8">
                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-500 mx-auto mb-2"></div>
                <div className="text-sm">Searching...</div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};