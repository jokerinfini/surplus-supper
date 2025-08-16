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
    console.log('Search button clicked, location:', location)
    try {
      setLoading(true)
      setError(null)
      console.log('Making API request to:', `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/restaurants`)
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
      navigator.geolocation.getCurrentPosition(
        (position) => {
          console.log('Location:', position.coords)
          // TODO: Reverse geocode and search
        },
        (error) => {
          console.error('Error getting location:', error)
        }
      )
    }
  }

  return (
    <div className="max-w-md mx-auto space-y-4">
      <form onSubmit={async (e) => {
        e.preventDefault();
        console.log('Form submitted');
        await handleSearch();
      }}>
        <div className="flex shadow-lg rounded-lg overflow-hidden">
          <input
            type="text"
            placeholder="Enter your location..."
            value={location}
            onChange={(e) => setLocation(e.target.value)}
            className="flex-1 px-4 py-3 text-gray-900 focus:outline-none focus:ring-2 focus:ring-green-400"
          />
          <Button
            type="submit"
            className="bg-green-600 hover:bg-green-700 px-6 py-3 rounded-none"
            disabled={loading}
          >
            <Search className="w-4 h-4" />
          </Button>
        </div>
        
        <div className="flex space-x-2 mt-2">
          <Button
            onClick={getCurrentLocation}
            variant="outline"
            type="button"
            className="flex items-center space-x-2 text-gray-700 hover:bg-gray-100"
          >
            <MapPin className="w-4 h-4" />
            <span>Use my location</span>
          </Button>
        </div>
      </form>
      
      {loading && (
        <div className="text-center py-4">
          <p className="text-gray-600">Searching restaurants...</p>
        </div>
      )}
      
      {error && (
        <div className="text-center py-4">
          <p className="text-red-600">{error}</p>
        </div>
      )}
    </div>
  );
}
