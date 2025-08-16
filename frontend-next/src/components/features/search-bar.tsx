'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Search, MapPin } from 'lucide-react'
import { api, type Restaurant } from '@/lib/api'

interface SearchBarProps {
  onSearch: (results: Restaurant[]) => void;
}

export function SearchBar({ onSearch }: SearchBarProps) {
  const [location, setLocation] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSearch = async () => {
    if (!location.trim()) {
      setError('Please enter a location')
      return
    }

    try {
      setLoading(true)
      setError(null)
      console.log('Searching for restaurants in:', location)
      
      const results = await api.getRestaurants(undefined, undefined, location)
      console.log('Search results:', results)
      onSearch(results)
    } catch (error) {
      console.error('Error searching:', error)
      setError(error instanceof Error ? error.message : 'Failed to search restaurants')
    } finally {
      setLoading(false)
    }
  }

  const getCurrentLocation = () => {
    if (navigator.geolocation) {
      setLoading(true)
      setError(null)
      
      navigator.geolocation.getCurrentPosition(
        async (position) => {
          try {
            console.log('Location obtained:', position.coords)
            // For now, we'll use a default location since we don't have reverse geocoding
            // In a real app, you'd reverse geocode the coordinates to get an address
            const results = await api.getRestaurants(position.coords.latitude, position.coords.longitude)
            onSearch(results)
          } catch (error) {
            console.error('Error with location search:', error)
            setError('Failed to search with current location')
          } finally {
            setLoading(false)
          }
        },
        (error) => {
          console.error('Error getting location:', error)
          setError('Unable to get your location. Please enter it manually.')
          setLoading(false)
        }
      )
    } else {
      setError('Geolocation is not supported by your browser')
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      handleSearch()
    }
  }

  return (
    <div className="w-full space-y-4">
      {/* Glassmorphism Search Bar */}
      <div className="p-2 rounded-lg flex flex-col sm:flex-row gap-2 glass-search">
        <div className="relative flex-grow">
          <span className="absolute inset-y-0 left-0 flex items-center pl-3">
            <MapPin className="text-white w-5 h-5" />
          </span>
          <input 
            type="text" 
            placeholder="Enter your city or address..." 
            value={location}
            onChange={(e) => setLocation(e.target.value)}
            onKeyPress={handleKeyPress}
            className="w-full pl-10 pr-4 py-3 focus:outline-none text-lg text-white placeholder-white/70 bg-transparent border-none"
            disabled={loading}
            style={{
              color: 'white',
              textShadow: '0 1px 2px rgba(0,0,0,0.3)'
            }}
          />
        </div>
        <button 
          onClick={handleSearch}
          disabled={loading || !location.trim()}
          className="w-full sm:w-auto flex-shrink-0 bg-primary-green text-white font-semibold px-6 py-3 rounded-md hover:bg-opacity-90 transition duration-200 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg"
        >
          {loading ? (
            <div className="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent"></div>
          ) : (
            <Search className="w-5 h-5" />
          )}
          <span>{loading ? 'Searching...' : 'Search'}</span>
        </button>
      </div>
      
      {/* Use My Location Button */}
      <button 
        onClick={getCurrentLocation}
        disabled={loading}
        className="text-white hover:underline flex items-center justify-center gap-2 mx-auto disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <span className="w-4 h-4">🎯</span>
        <span>Use my location</span>
      </button>
      
      {/* Error Message */}
      {error && (
        <div className="text-center py-3">
          <div className="inline-flex items-center space-x-2 bg-red-50 border border-red-200 text-red-700 px-4 py-2 rounded-lg">
            <div className="w-2 h-2 bg-red-500 rounded-full"></div>
            <p className="text-sm font-medium">{error}</p>
          </div>
        </div>
      )}
      
      {/* Loading Message */}
      {loading && !error && (
        <div className="text-center py-3">
          <div className="inline-flex items-center space-x-2 bg-blue-50 border border-blue-200 text-blue-700 px-4 py-2 rounded-lg">
            <div className="animate-spin rounded-full h-4 w-4 border-2 border-blue-500 border-t-transparent"></div>
            <p className="text-sm font-medium">Searching for restaurants near you...</p>
          </div>
        </div>
      )}
    </div>
  );
}
