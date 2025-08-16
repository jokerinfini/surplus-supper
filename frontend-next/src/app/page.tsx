'use client'

import { motion } from 'framer-motion'
import { useState } from 'react'
import { Search, Store, Utensils, Heart, ShoppingCart, Menu, Motorcycle, Facebook, Twitter, Instagram } from 'lucide-react'
import { SearchBar } from '@/components/features/search-bar'
import { RestaurantCard } from '@/components/features/restaurant-card'
import type { Restaurant } from '@/lib/api'

// Recipe Modal Component
function RecipeModal({ isOpen, onClose, ingredients, restaurantName }: any) {
  const [loading, setLoading] = useState(false)
  const [recipe, setRecipe] = useState('')

  const generateRecipe = async () => {
    setLoading(true)
    try {
      const prompt = `You are a creative chef helping a user reduce food waste. Create a simple and delicious recipe using the following surplus ingredients from "${restaurantName}": ${ingredients}. 
      
      The recipe should be easy for a home cook to follow. 
      
      Please format the response in Markdown with the following sections:
      - A catchy recipe title.
      - A brief, encouraging introduction.
      - A list of additional ingredients the user might need (if any).
      - Step-by-step instructions.
      - A concluding tip on food storage or another way to use one of the ingredients.`

      // Simulate API call (replace with actual Gemini API)
      await new Promise(resolve => setTimeout(resolve, 2000))
      
      const mockRecipe = `# 🍝 Surplus Pasta Delight

**A delicious way to use your surplus ingredients!**

## Additional Ingredients Needed:
- Olive oil (2 tbsp)
- Salt and pepper to taste
- Fresh herbs (optional)

## Instructions:
1. Heat olive oil in a large pan over medium heat
2. Add chopped onions and garlic, sauté until fragrant
3. Add ground beef and cook until browned
4. Stir in canned tomatoes and simmer for 10 minutes
5. Cook spaghetti according to package directions
6. Combine pasta with sauce and serve hot

## Storage Tip:
Leftover sauce can be frozen for up to 3 months!`

      setRecipe(mockRecipe)
    } catch (error) {
      console.error('Error generating recipe:', error)
      setRecipe('Sorry, I couldn\'t generate a recipe at this moment. Please try again later.')
    } finally {
      setLoading(false)
    }
  }

  if (!isOpen) return null

  return (
    <div 
      className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
      onClick={onClose}
    >
      <div 
        className="bg-white rounded-lg shadow-2xl w-full max-w-2xl max-h-[90vh] overflow-y-auto"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="p-6 border-b sticky top-0 bg-white z-10">
          <div className="flex justify-between items-center">
            <h3 className="text-2xl font-bold text-gray-800">✨ Your Recipe Idea</h3>
            <button 
              onClick={onClose}
              className="text-gray-500 hover:text-gray-800 text-2xl"
            >
              &times;
            </button>
          </div>
        </div>
        <div className="p-6 space-y-4">
          {loading ? (
            <div className="text-center py-10">
              <div className="animate-spin text-4xl text-green-600 mb-4">🍽️</div>
              <p className="text-gray-600">Generating your recipe...</p>
            </div>
          ) : recipe ? (
            <div 
              className="prose max-w-none"
              dangerouslySetInnerHTML={{ 
                __html: recipe
                  .replace(/^# (.*$)/gim, '<h1 class="text-2xl font-bold mb-2">$1</h1>')
                  .replace(/^## (.*$)/gim, '<h2 class="text-xl font-semibold mt-4 mb-2">$1</h2>')
                  .replace(/^### (.*$)/gim, '<h3 class="text-lg font-semibold mt-3 mb-1">$1</h3>')
                  .replace(/\*\*(.*)\*\*/g, '<strong>$1</strong>')
                  .replace(/\*(.*)\*/g, '<em>$1</em>')
                  .replace(/^- (.*$)/gim, '<li class="ml-4 list-disc">$1</li>')
                  .replace(/^\d+\. (.*$)/gim, '<li class="ml-4 list-decimal">$1</li>')
                  .replace(/\n/g, '<br>')
              }}
            />
          ) : (
            <div className="text-center">
              <button 
                onClick={generateRecipe}
                className="bg-gradient-to-r from-purple-500 to-indigo-600 text-white py-3 px-6 rounded-lg font-bold hover:from-purple-600 hover:to-indigo-700 transition-all"
              >
                Generate Recipe
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

// Modal state management
function useRestaurantModal() {
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [currentRestaurant, setCurrentRestaurant] = useState<Restaurant | null>(null)

  const openModal = (restaurant: Restaurant) => {
    setCurrentRestaurant(restaurant)
    setIsModalOpen(true)
  }

  const closeModal = () => {
    setIsModalOpen(false)
    setCurrentRestaurant(null)
  }

  return {
    isModalOpen,
    currentRestaurant,
    openModal,
    closeModal
  }
}

export default function HomePage() {
  const [searchResults, setSearchResults] = useState<Restaurant[]>([])
  const [showResults, setShowResults] = useState(false)
  const { isModalOpen, currentRestaurant, openModal, closeModal } = useRestaurantModal()

  const handleSearch = (results: Restaurant[]) => {
    setSearchResults(results);
    setShowResults(true);
  };

  return (
    <div className="bg-green-50 font-['Poppins']">
      {/* Header */}
      <header className="bg-white shadow-md sticky top-0 z-50">
        <nav className="container mx-auto px-6 py-3">
          <div className="flex justify-between items-center">
            <motion.a 
              href="#" 
              className="text-2xl font-bold text-green-600"
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
            >
              Surplus Supper
            </motion.a>
            <div className="hidden md:flex items-center space-x-6">
              <a href="#" className="text-gray-600 hover:text-green-600 transition-colors">Home</a>
              <a href="#" className="text-gray-600 hover:text-green-600 transition-colors">Restaurants</a>
              <a href="#" className="text-gray-600 hover:text-green-600 transition-colors">About Us</a>
              <a href="#" className="text-gray-600 hover:text-green-600 transition-colors">Contact</a>
            </div>
            <div className="flex items-center space-x-4">
              <button className="text-gray-600 hover:text-green-600 transition-colors">
                <ShoppingCart className="h-5 w-5" />
              </button>
              <button className="md:hidden text-gray-600 hover:text-green-600 transition-colors">
                <Menu className="h-5 w-5" />
              </button>
            </div>
          </div>
        </nav>
      </header>

      {/* Search Results */}
      {showResults && searchResults.length > 0 && (
        <section className="py-8">
          <div className="container mx-auto px-6">
            <h2 className="text-2xl font-bold mb-6">Search Results</h2>
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
              {searchResults.map((restaurant) => (
                <div key={restaurant.id} className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
                  <div className="p-6">
                    <h3 className="text-xl font-bold mb-2">{restaurant.name}</h3>
                    <p className="text-gray-600 mb-4">{restaurant.description}</p>
                    <div className="text-sm text-gray-500">
                      <p className="mb-1">{restaurant.address}</p>
                      <p className="mb-1">Cuisine: {restaurant.cuisine_type}</p>
                      <p>Rating: {restaurant.rating}⭐</p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </section>
      )}

      {/* Hero Section */}
      <section className="relative text-white py-20 overflow-hidden">
        <div 
          className="absolute inset-0 bg-cover bg-center bg-no-repeat"
          style={{
            backgroundImage: `linear-gradient(rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0.5)), url('https://images.unsplash.com/photo-1504674900247-0877df9cc836?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4-0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')`
          }}
        />
        
        {/* Animated Road */}
        <div className="absolute bottom-0 left-0 w-full h-6 bg-black bg-opacity-60 border-t-4 border-dashed border-white border-opacity-70" />
        
        {/* Animated Scooter - Fixed direction and using white motorcycle */}
        <motion.div 
          className="absolute bottom-4 left-0 text-6xl text-white"
          animate={{ x: ['-150px', 'calc(100vw + 200px)'] }}
          transition={{ duration: 12, repeat: Infinity, ease: 'linear' }}
          style={{ transform: 'scaleX(-1)' }} // Flip horizontally to face right direction
        >
          🏍️
        </motion.div>

        <div className="container mx-auto text-center relative z-10">
          <motion.h1 
            className="text-4xl md:text-6xl font-bold mb-4"
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
          >
            Delicious Meals, Zero Waste
          </motion.h1>
          <motion.p 
            className="text-lg md:text-xl mb-8"
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.2 }}
          >
            Get surplus food from your favorite restaurants at a discounted price.
          </motion.p>
          <motion.div 
            className="w-full max-w-xl mx-auto"
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.4 }}
          >
            <div className="bg-white/90 backdrop-blur-sm rounded-lg shadow-lg p-4">
              <SearchBar onSearch={handleSearch} />
            </div>
          </motion.div>
        </div>
      </section>

      {/* How it works */}
      <section className="py-16 bg-green-100">
        <div className="container mx-auto text-center px-4">
          <motion.h2 
            className="text-3xl font-bold mb-8 text-gray-800"
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
          >
            How It Works
          </motion.h2>
          <div className="grid md:grid-cols-3 gap-8">
            <motion.div 
              className="bg-white p-6 rounded-lg shadow-lg hover:shadow-xl transition-shadow"
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.1 }}
            >
              <div className="text-5xl text-green-600 mb-4">🏪</div>
              <h3 className="text-xl font-bold mb-2">Find Restaurants</h3>
              <p className="text-gray-600">Browse local restaurants with surplus food.</p>
            </motion.div>
            <motion.div 
              className="bg-white p-6 rounded-lg shadow-lg hover:shadow-xl transition-shadow"
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
            >
              <div className="text-5xl text-green-600 mb-4">🍽️</div>
              <h3 className="text-xl font-bold mb-2">Choose Your Meal</h3>
              <p className="text-gray-600">Select from a variety of delicious surplus meals.</p>
            </motion.div>
            <motion.div 
              className="bg-white p-6 rounded-lg shadow-lg hover:shadow-xl transition-shadow"
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.3 }}
            >
              <div className="text-5xl text-green-600 mb-4">💚</div>
              <h3 className="text-xl font-bold mb-2">Save Food & Money</h3>
              <p className="text-gray-600">Enjoy your meal and feel good about reducing food waste.</p>
            </motion.div>
          </div>
        </div>
      </section>

      {/* Featured Restaurants */}
      <section className="py-16 bg-green-800">
        <div className="container mx-auto px-4">
          <motion.h2 
            className="text-3xl font-bold text-center mb-8 text-white"
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
          >
            {showResults ? 'Search Results' : 'Featured Restaurants'}
          </motion.h2>
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            {showResults ? (
              searchResults.length > 0 ? (
                searchResults.map((restaurant, index) => (
                  <RestaurantCard 
                    key={restaurant.id || index} 
                    restaurant={restaurant}
                    onSelect={openModal}
                  />
                ))
              ) : (
                <div className="col-span-full text-center text-white">
                  <p>No restaurants found. Try another location.</p>
                </div>
              )
            ) : (
              // Placeholder for featured restaurants or loading state
              <div className="col-span-full text-center text-white">
                <p>Search for restaurants to see results.</p>
              </div>
            )}
          </div>
        </div>
      </section>

      {isModalOpen && currentRestaurant && (
        <RecipeModal
          isOpen={isModalOpen}
          onClose={closeModal}
          ingredients={currentRestaurant.surplus_items?.join(', ') || ''}
          restaurantName={currentRestaurant.name}
        />
      )}

      {/* Mobile App Promo */}
      <section className="py-16 bg-green-100">
        <div className="container mx-auto px-4">
          <motion.div 
            className="flex flex-col md:flex-row items-center bg-white p-8 rounded-lg shadow-lg"
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
          >
            <div className="md:w-1/2 text-center md:text-left mb-8 md:mb-0">
              <h2 className="text-3xl font-bold mb-4">Get the Surplus Supper App</h2>
              <p className="text-gray-600 mb-6">Download our app to easily order surplus food on the go. Available on both iOS and Android.</p>
              <div className="flex justify-center md:justify-start space-x-4">
                <a href="#" className="hover:opacity-80 transition-opacity">
                  <img 
                    src="https://upload.wikimedia.org/wikipedia/commons/thumb/3/3c/Download_on_the_App_Store_Badge.svg/1200px-Download_on_the_App_Store_Badge.svg.png" 
                    alt="App Store" 
                    className="h-12"
                  />
                </a>
                <a href="#" className="hover:opacity-80 transition-opacity">
                  <img 
                    src="https://upload.wikimedia.org/wikipedia/commons/thumb/7/78/Google_Play_Store_badge_EN.svg/1200px-Google_Play_Store_badge_EN.svg.png" 
                    alt="Google Play" 
                    className="h-12"
                  />
                </a>
              </div>
            </div>
            <div className="md:w-1/2 flex justify-center">
              <img 
                src="https://cdn.dribbble.com/users/1043333/screenshots/11649021/media/cc12c2a681375744837f19b049360c61.png" 
                alt="Mobile App" 
                className="max-w-xs"
              />
            </div>
          </motion.div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-green-900 text-white py-8">
        <div className="container mx-auto text-center px-4">
          <p>&copy; 2024 Surplus Supper. All rights reserved.</p>
          <div className="flex justify-center space-x-4 mt-4">
            <a href="#" className="hover:text-green-400 transition-colors">
              <Facebook className="h-5 w-5" />
            </a>
            <a href="#" className="hover:text-green-400 transition-colors">
              <Twitter className="h-5 w-5" />
            </a>
            <a href="#" className="hover:text-green-400 transition-colors">
              <Instagram className="h-5 w-5" />
            </a>
          </div>
        </div>
      </footer>
    </div>
  )
}
