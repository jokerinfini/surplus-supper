'use client'

import { type Restaurant } from '@/lib/api'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { MapPin, Clock, Tag } from 'lucide-react'

interface RestaurantCardProps {
  restaurant: Restaurant;
  onSelect?: (restaurant: Restaurant) => void;
}

export function RestaurantCard({ restaurant, onSelect }: RestaurantCardProps) {
  return (
    <Card className="overflow-hidden hover:shadow-lg transition-shadow duration-300">
      <div className="relative h-48">
        <img
          src={'https://placehold.co/600x400?text=' + encodeURIComponent(restaurant.name)}
          alt={restaurant.name}
          className="w-full h-full object-cover"
        />
        <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/70 to-transparent p-4">
          <h3 className="text-xl font-semibold text-white">{restaurant.name}</h3>
        </div>
      </div>
      
      <div className="p-4 space-y-4">
        <div className="space-y-2">
          <div className="flex items-center text-gray-600">
            <MapPin className="w-4 h-4 mr-2" />
            <span>{restaurant.address}</span>
          </div>
          <div className="flex items-center text-gray-600">
            <Tag className="w-4 h-4 mr-2" />
            <span>{restaurant.cuisine_type}</span>
          </div>
          {restaurant.distance && (
            <div className="flex items-center text-gray-600">
              <Clock className="w-4 h-4 mr-2" />
              <span>{restaurant.distance.toFixed(1)} km away</span>
            </div>
          )}
        </div>

        <Button
          onClick={() => onSelect?.(restaurant)}
          className="w-full bg-green-600 hover:bg-green-700 text-white"
        >
          View Details
        </Button>
      </div>
    </Card>
  )
}
