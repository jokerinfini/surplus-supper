const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export interface Restaurant {
  id: number
  name: string
  description: string
  address: string
  latitude: number
  longitude: number
  phone: string
  email: string
  cuisine_type: string
  rating: number
  distance: number
  surplus_items?: string[]
}

export interface InventoryItem {
  id: number
  name: string
  description: string
  original_price: number
  surplus_price: number
  quantity: number
  category: string
  expiry_time: string
  discount: number
}

export const api = {
  async getRestaurants(lat?: number, lng?: number, location?: string): Promise<Restaurant[]> {
    try {
      const params = new URLSearchParams()
      if (lat && lng) {
        params.append('lat', lat.toString())
        params.append('lng', lng.toString())
      }
      if (location) {
        params.append('location', location)
      }

      const url = `${API_BASE}/api/restaurants?${params}`
      console.log('Making request to:', url)

      const response = await fetch(url)
      console.log('Response status:', response.status)
      
      const text = await response.text()
      console.log('Response text:', text)

      if (!response.ok) {
        throw new Error(`Failed to fetch restaurants: ${response.status} ${text}`)
      }

      try {
        return JSON.parse(text)
      } catch (e) {
        console.error('Failed to parse JSON:', e)
        throw new Error('Invalid JSON response from server')
      }
    } catch (error) {
      console.error('API Error:', error)
      throw error
    }
  },

  async getRestaurant(id: number): Promise<{ restaurant: Restaurant; inventory: InventoryItem[] }> {
    const response = await fetch(`${API_BASE}/restaurant/${id}`)
    if (!response.ok) {
      throw new Error('Failed to fetch restaurant')
    }
    return response.json()
  },

  async createOrder(orderData: any): Promise<any> {
    const response = await fetch(`${API_BASE}/order/confirm`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(orderData),
    })
    if (!response.ok) {
      throw new Error('Failed to create order')
    }
    return response.json()
  }
}
