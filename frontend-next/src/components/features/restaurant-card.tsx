'use client'

import { type Restaurant } from '@/lib/api'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { MapPin, Star } from 'lucide-react'

interface RestaurantCardProps {
  restaurant: Restaurant;
  onSelect?: (restaurant: Restaurant) => void;
}

export function RestaurantCard({ restaurant, onSelect }: RestaurantCardProps) {
  return (
    <Card 
      className="overflow-hidden hover:shadow-xl transition-all duration-300 transform hover:-translate-y-2 card-3d group relative bg-white rounded-lg shadow-lg"
    >
      <div className="relative">
        <img
          src={'https://placehold.co/600x400/2A9D8F/ffffff?text=' + encodeURIComponent(restaurant.name)}
          alt={restaurant.name}
          className="h-56 w-full object-cover transition-transform duration-300 group-hover:scale-105"
        />
        <div className="absolute top-4 right-4 bg-green-500 text-white text-sm font-bold px-3 py-1 rounded-full flex items-center gap-1 shadow-lg">
          <Star className="w-4 h-4 fill-current" />
          <span>{restaurant.rating}</span>
        </div>
        {/* Hover overlay */}
        <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all duration-300"></div>
      </div>
      
      <div className="p-6 relative z-10">
        <h3 className="text-2xl font-display font-bold text-dark-text group-hover:text-primary-green transition-colors duration-300">{restaurant.name}</h3>
        <p className="mt-2 text-gray-600">{restaurant.description}</p>
        
        <div className="mt-4 space-y-2 text-sm text-gray-700">
          <p className="flex items-center">
            <MapPin className="mr-2 h-4 w-4 text-primary-green" />
            <span className="truncate">{restaurant.address}</span>
          </p>
          <p className="flex items-center">
            <span className="mr-2 h-4 w-4 text-primary-green">🏷️</span>
            <span>{restaurant.cuisine_type}</span>
          </p>
          {restaurant.distance && (
            <p className="flex items-center">
              <span className="mr-2 h-4 w-4 text-primary-green">⏰</span>
              <span>{restaurant.distance.toFixed(1)} km away</span>
            </p>
          )}
        </div>

        <div className="mt-6 flex justify-between items-center">
          <div className="text-xl font-bold text-dark-text">
            $7.99 <span className="text-sm font-normal text-gray-500 line-through">~$25</span>
          </div>
          <div className="flex space-x-2">
            <Button
              onClick={() => onSelect?.(restaurant)}
              className="bg-primary-green text-white font-semibold px-4 py-2 rounded-md hover:bg-opacity-90 transition duration-200 transform hover:scale-105"
            >
              View Details
            </Button>
            <Button
              variant="outline"
              className="border-primary-green text-primary-green hover:bg-primary-green hover:text-white transform hover:scale-105 transition-all duration-200"
            >
              Order Now
            </Button>
          </div>
        </div>
      </div>
    </Card>
  )
}
