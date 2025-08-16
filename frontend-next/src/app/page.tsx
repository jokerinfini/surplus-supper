'use client'

import { motion } from 'framer-motion'
import { useState, useEffect, useRef } from 'react'
import { Search, MapPin, Star, ShoppingCart, Facebook, Twitter, Instagram } from 'lucide-react'
import { SearchBar } from '@/components/features/search-bar'
import { RestaurantCard } from '@/components/features/restaurant-card'
import type { Restaurant } from '@/lib/api'
import AuthModal from '@/components/auth/auth-modal'
import { isAuthenticated, getUser, logout } from '@/lib/auth'
import type { User } from '@/lib/auth'

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

// Falling Food Animation Component
function FallingFoodAnimation() {
  const canvasRef = useRef<HTMLCanvasElement>(null)

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return

    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const PIXEL_SIZE = 5
    const GRAVITY = 0.05
    const SPAWN_INTERVAL = 200

    const PALETTE: { [key: number]: string | null } = {
      0: null, 1: '#ff6b6b', 2: '#f9d74c', 3: '#8d5524', 4: '#6aaa64',
      5: '#ffffff', 6: '#ff9f43', 7: '#e0aaff', 8: '#f08080', 9: '#ffdab9',
      10: '#2f3542', 11: '#ff7979'
    }

    const foodDesigns = [
      { width: 8, height: 7, data: [0, 3, 3, 3, 3, 3, 3, 0, 3, 5, 2, 2, 2, 2, 5, 3, 5, 4, 4, 4, 4, 4, 4, 5, 5, 1, 1, 1, 1, 1, 1, 5, 5, 3, 3, 3, 3, 3, 3, 5, 3, 5, 5, 5, 5, 5, 5, 3, 0, 3, 3, 3, 3, 3, 3, 0] },
      { width: 8, height: 8, data: [0, 0, 0, 0, 3, 3, 3, 3, 0, 0, 0, 3, 2, 2, 2, 3, 0, 0, 3, 2, 11, 2, 11, 3, 0, 3, 2, 2, 2, 11, 2, 3, 3, 2, 11, 2, 2, 2, 2, 3, 3, 2, 2, 2, 11, 2, 2, 3, 3, 2, 2, 2, 2, 2, 3, 0, 3, 3, 3, 3, 3, 3, 0, 0] },
      { width: 8, height: 7, data: [0, 5, 5, 5, 5, 5, 5, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 10, 10, 10, 10, 10, 10, 0, 0, 10, 5, 5, 5, 5, 10, 0, 0, 10, 10, 10, 10, 10, 10, 0] },
      { width: 8, height: 8, data: [0, 9, 9, 8, 8, 8, 9, 0, 9, 8, 8, 5, 5, 8, 8, 9, 9, 8, 5, 5, 5, 5, 8, 9, 8, 8, 5, 5, 5, 5, 8, 8, 8, 8, 5, 5, 5, 5, 8, 8, 9, 8, 5, 5, 5, 5, 8, 9, 9, 8, 8, 5, 5, 8, 8, 9, 0, 9, 9, 8, 8, 8, 9, 0] },
      { width: 7, height: 8, data: [0, 0, 4, 4, 4, 0, 0, 0, 4, 4, 4, 4, 4, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 5, 1, 1, 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 5, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0] },
      { width: 8, height: 8, data: [0, 6, 6, 6, 6, 6, 6, 0, 6, 5, 5, 5, 5, 5, 5, 6, 6, 5, 6, 5, 5, 6, 5, 6, 6, 5, 5, 6, 6, 5, 5, 6, 6, 5, 5, 6, 6, 5, 5, 6, 6, 5, 6, 5, 5, 6, 5, 6, 6, 5, 5, 5, 5, 5, 5, 6, 0, 6, 6, 6, 6, 6, 6, 0] },
      { width: 7, height: 8, data: [0, 0, 0, 4, 4, 0, 0, 0, 0, 4, 4, 0, 0, 0, 0, 0, 7, 7, 7, 0, 0, 0, 7, 7, 7, 7, 7, 0, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 0, 7, 7, 7, 7, 7, 0, 0, 0, 7, 7, 7, 0, 0] }
    ]

    let fallingFoods: any[] = []
    let lastSpawnTime = 0

    class Food {
      x: number
      y: number
      vy: number
      rotation: number
      rotationSpeed: number
      bounced: boolean
      alpha: number
      design: any

      constructor(x: number, y: number, design: any) {
        this.design = design
        this.x = x
        this.y = y
        this.vy = Math.random() * 1 + 0.5
        this.rotation = 0
        this.rotationSpeed = (Math.random() - 0.5) * 0.02
        this.bounced = false
        this.alpha = 1
      }

      update() {
        this.vy += GRAVITY
        this.y += this.vy
        this.rotation += this.rotationSpeed
        const groundY = canvas.height - this.design.height * PIXEL_SIZE
        if (this.y > groundY && !this.bounced) {
          this.y = groundY
          this.vy *= -0.4
          this.rotationSpeed *= 0.5
          this.bounced = true
        }
        if (this.bounced) {
          this.alpha -= 0.02
        }
      }

      draw() {
        if (!ctx) return
        ctx.save()
        ctx.globalAlpha = this.alpha
        const centerX = this.x + (this.design.width * PIXEL_SIZE) / 2
        const centerY = this.y + (this.design.height * PIXEL_SIZE) / 2
        ctx.translate(centerX, centerY)
        ctx.rotate(this.rotation)
        ctx.translate(-centerX, -centerY)
        for (let i = 0; i < this.design.data.length; i++) {
          const colorIndex = this.design.data[i]
          if (PALETTE[colorIndex]) {
            const col = i % this.design.width
            const row = Math.floor(i / this.design.width)
            ctx.fillStyle = PALETTE[colorIndex]
            ctx.fillRect(this.x + col * PIXEL_SIZE, this.y + row * PIXEL_SIZE, PIXEL_SIZE, PIXEL_SIZE)
          }
        }
        ctx.restore()
      }
    }

    function resizeCanvas() {
      if (canvas && canvas.parentElement) {
        canvas.width = canvas.parentElement.offsetWidth
        canvas.height = canvas.parentElement.offsetHeight
      }
    }

    function spawnFood() {
      if (!canvas) return
      const design = foodDesigns[Math.floor(Math.random() * foodDesigns.length)]
      const x = Math.random() * (canvas.width - design.width * PIXEL_SIZE)
      const y = -design.height * PIXEL_SIZE
      fallingFoods.push(new Food(x, y, design))
    }

    function animate(timestamp: number) {
      if (!ctx || !canvas) return
      ctx.clearRect(0, 0, canvas.width, canvas.height)
      if (timestamp - lastSpawnTime > SPAWN_INTERVAL) {
        spawnFood()
        lastSpawnTime = timestamp
      }
      for (let i = fallingFoods.length - 1; i >= 0; i--) {
        const food = fallingFoods[i]
        food.update()
        food.draw()
        if (food.alpha <= 0) {
          fallingFoods.splice(i, 1)
        }
      }
      requestAnimationFrame(animate)
    }

    const handleResize = () => resizeCanvas()
    window.addEventListener('resize', handleResize)
    resizeCanvas()
    requestAnimationFrame(animate)

    return () => {
      window.removeEventListener('resize', handleResize)
    }
  }, [])

  return <canvas ref={canvasRef} className="absolute inset-0 w-full h-full pointer-events-none" />
}

export default function HomePage() {
  const [searchResults, setSearchResults] = useState<Restaurant[]>([])
  const [showResults, setShowResults] = useState(false)
  const { isModalOpen, currentRestaurant, openModal, closeModal } = useRestaurantModal()
  const [isAuthModalOpen, setIsAuthModalOpen] = useState(false)
  const [user, setUser] = useState<User | null>(null)

  const handleSearch = (results: Restaurant[]) => {
    setSearchResults(results);
    setShowResults(true);
  };

  const handleHomeClick = () => {
    setShowResults(false);
    setSearchResults([]);
  };

  // Authentication effects
  useEffect(() => {
    const checkAuth = () => {
      console.log('Checking authentication...');
      if (isAuthenticated()) {
        const currentUser = getUser();
        console.log('User is authenticated:', currentUser);
        setUser(currentUser);
      } else {
        console.log('User is not authenticated');
      }
    };
    
    checkAuth();
  }, []);

  const handleAuthSuccess = () => {
    const currentUser = getUser();
    setUser(currentUser);
  };

  const handleLogout = () => {
    logout();
    setUser(null);
  };

  // Card 3D effects
  useEffect(() => {
    const cards = document.querySelectorAll('.card-3d');
    
    cards.forEach(card => {
      // Listen for the mouse moving over the card
      const handleMouseMove = (e: Event) => {
        const mouseEvent = e as MouseEvent;
        const rect = card.getBoundingClientRect();
        // Calculate mouse position relative to the card's top-left corner
        const x = mouseEvent.clientX - rect.left;
        const y = mouseEvent.clientY - rect.top;

        const centerX = rect.width / 2;
        const centerY = rect.height / 2;

        // Calculate the rotation angle based on how far the mouse is from the center
        const rotateX = ((y - centerY) / centerY) * -10; // Max tilt of 10 degrees
        const rotateY = ((x - centerX) / centerX) * 10;  // Max tilt of 10 degrees

        // Update the CSS variables with the new rotation values
        (card as HTMLElement).style.setProperty('--rotateX', `${rotateX}deg`);
        (card as HTMLElement).style.setProperty('--rotateY', `${rotateY}deg`);
      };

      // When the mouse leaves, reset the card to its flat state
      const handleMouseLeave = () => {
        (card as HTMLElement).style.setProperty('--rotateX', '0deg');
        (card as HTMLElement).style.setProperty('--rotateY', '0deg');
      };

      card.addEventListener('mousemove', handleMouseMove);
      card.addEventListener('mouseleave', handleMouseLeave);

      // Cleanup function
      return () => {
        card.removeEventListener('mousemove', handleMouseMove);
        card.removeEventListener('mouseleave', handleMouseLeave);
      };
    });
  }, [searchResults, showResults]); // Re-run when cards change

  return (
    <div className="bg-light-bg font-sans text-dark-text">
      {/* Header */}
      <header className="bg-white shadow-sm sticky top-0 z-50 w-full">
        <nav className="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex-shrink-0">
              <motion.a 
                href="#" 
                className="text-2xl font-display font-bold text-primary-green"
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
              >
                Surplus Supper
              </motion.a>
            </div>
            <div className="hidden md:block">
              <div className="ml-10 flex items-baseline space-x-4">
                <a href="#" onClick={handleHomeClick} className="text-dark-text hover:text-primary-green px-3 py-2 rounded-md text-sm font-medium transition-colors cursor-pointer">Home</a>
                <a href="#" className="text-dark-text hover:text-primary-green px-3 py-2 rounded-md text-sm font-medium transition-colors">Restaurants</a>
                <a href="#" className="text-dark-text hover:text-primary-green px-3 py-2 rounded-md text-sm font-medium transition-colors">About Us</a>
                <a href="#" className="text-dark-text hover:text-primary-green px-3 py-2 rounded-md text-sm font-medium transition-colors">Contact</a>
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <button className="p-2 rounded-full text-dark-text hover:bg-gray-100 focus:outline-none transition-colors">
                <ShoppingCart className="h-5 w-5" />
              </button>
              
              {/* Debug info */}
              <span className="text-xs text-gray-500">
                Auth: {user ? 'Logged In' : 'Not Logged In'}
              </span>
              
              {user ? (
                <div className="flex items-center space-x-2">
                  <span className="text-sm text-gray-700">
                    Hi, {user.first_name}!
                  </span>
                  <button
                    onClick={handleLogout}
                    className="text-sm text-gray-600 hover:text-primary-green transition-colors"
                  >
                    Logout
                  </button>
                </div>
              ) : (
                <button
                  onClick={() => {
                    console.log('Sign In button clicked!');
                    setIsAuthModalOpen(true);
                  }}
                  className="bg-green-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-green-700 transition-colors shadow-md border border-green-700 min-w-[100px]"
                >
                  Sign In
                </button>
              )}
            </div>
          </div>
        </nav>
      </header>

      {/* Hero Section */}
      <section className="relative bg-dark-text">
        {/* Canvas for falling food animation */}
        <FallingFoodAnimation />
        <div className="absolute inset-0 bg-black opacity-50 z-10"></div>
        
        <div className="hero-content w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 md:py-32 lg:py-40 text-center text-white relative z-20">
          <motion.h1 
            className="text-4xl sm:text-5xl md:text-6xl font-display font-bold tracking-tight"
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
          >
            Delicious Meals, Zero Waste
          </motion.h1>
          <motion.p 
            className="mt-4 max-w-2xl mx-auto text-lg sm:text-xl"
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.2 }}
          >
            Get surplus food from your favorite restaurants at a discounted price.
          </motion.p>
          
          {/* Glassmorphism Search Bar */}
          <motion.div 
            className="mt-8 max-w-2xl mx-auto"
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.4 }}
          >
            <div className="bg-white/10 backdrop-blur-md rounded-3xl shadow-2xl p-8 border border-white/20">
              <SearchBar onSearch={handleSearch} />
            </div>
          </motion.div>
        </div>
      </section>

      {/* Restaurants Section */}
      <section className="py-16 sm:py-20 w-full">
        <div className="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          {showResults ? (
            <>
              <h2 className="text-3xl font-display font-bold text-center text-dark-text mb-8">
                Restaurants Near You
              </h2>
              {searchResults.length === 0 ? (
                <div className="text-center py-12">
                  <p className="text-gray-600 text-lg mb-4">No restaurants found in your area.</p>
                  <button 
                    onClick={handleHomeClick}
                    className="bg-primary-green text-white px-6 py-3 rounded-md hover:bg-opacity-90 transition-colors"
                  >
                    Search Again
                  </button>
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 w-full">
                  {searchResults.map((restaurant, index) => (
                    <motion.div
                      key={restaurant.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ duration: 0.6, delay: index * 0.1 }}
                      className="w-full"
                    >
                      <RestaurantCard restaurant={restaurant} onSelect={openModal} />
                    </motion.div>
                  ))}
                </div>
              )}
            </>
          ) : (
            <>
              <h2 className="text-3xl font-display font-bold text-center text-dark-text mb-8">
                Featured Restaurants
              </h2>
              <p className="text-center text-gray-600 max-w-xl mx-auto mb-12">
                Discover surprise bags of delicious, unsold food from a variety of local spots.
              </p>
              
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 w-full">
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.6, delay: 0.2 }}
                >
                  <RestaurantCard 
                    restaurant={{
                      id: 1,
                      name: "Pizza Palace",
                      description: "Authentic Italian pizza and pasta.",
                      address: "123 Main St, New York, NY",
                      cuisine_type: "Italian",
                      rating: 4.5,
                      distance: 0.5,
                      surplus_items: ["Pizza Margherita", "Pasta Carbonara", "Garlic Bread"],
                      latitude: 40.7128,
                      longitude: -74.0060,
                      phone: "+1-555-0123",
                      email: "info@pizzapalace.com"
                    }}
                    onSelect={openModal}
                  />
                </motion.div>

                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.6, delay: 0.3 }}
                >
                  <RestaurantCard 
                    restaurant={{
                      id: 2,
                      name: "Sushi Express",
                      description: "Fresh sushi and Japanese cuisine.",
                      address: "456 Broadway, New York, NY",
                      cuisine_type: "Japanese",
                      rating: 4.8,
                      distance: 0.8,
                      surplus_items: ["Salmon Nigiri", "California Roll", "Miso Soup"],
                      latitude: 40.7589,
                      longitude: -73.9851,
                      phone: "+1-555-0456",
                      email: "info@sushiexpress.com"
                    }}
                    onSelect={openModal}
                  />
                </motion.div>

                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.6, delay: 0.4 }}
                >
                  <RestaurantCard 
                    restaurant={{
                      id: 3,
                      name: "Burger Joint",
                      description: "Classic American burgers and fries.",
                      address: "789 Oak Ave, New York, NY",
                      cuisine_type: "American",
                      rating: 4.3,
                      distance: 1.2,
                      surplus_items: ["Cheeseburger", "French Fries", "Onion Rings"],
                      latitude: 40.7505,
                      longitude: -73.9934,
                      phone: "+1-555-0789",
                      email: "info@burgerjoint.com"
                    }}
                    onSelect={openModal}
                  />
                </motion.div>
              </div>
            </>
          )}
        </div>
      </section>

      {/* Recipe Modal */}
      {isModalOpen && currentRestaurant && (
        <RecipeModal
          isOpen={isModalOpen}
          onClose={closeModal}
          ingredients={currentRestaurant.surplus_items?.join(', ') || ''}
          restaurantName={currentRestaurant.name}
        />
      )}

      {/* Authentication Modal */}
      <AuthModal
        isOpen={isAuthModalOpen}
        onClose={() => setIsAuthModalOpen(false)}
        onSuccess={handleAuthSuccess}
        initialMode="login"
      />

      {/* Footer */}
      <footer className="text-white py-12 w-full" style={{ backgroundColor: '#264653' }}>
        <div className="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
            <div>
              <h3 className="text-xl font-display font-bold text-primary-green">Surplus Supper</h3>
              <p className="mt-2 text-gray-300 text-sm">Eat Good. Do Good.</p>
            </div>
            <div>
              <h4 className="font-semibold tracking-wider uppercase text-white">Quick Links</h4>
              <ul className="mt-4 space-y-2 text-sm">
                <li><a href="#" className="text-gray-300 hover:text-primary-green transition-colors">Home</a></li>
                <li><a href="#" className="text-gray-300 hover:text-primary-green transition-colors">About Us</a></li>
                <li><a href="#" className="text-gray-300 hover:text-primary-green transition-colors">Become a Partner</a></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold tracking-wider uppercase text-white">Support</h4>
              <ul className="mt-4 space-y-2 text-sm">
                <li><a href="#" className="text-gray-300 hover:text-primary-green transition-colors">FAQ</a></li>
                <li><a href="#" className="text-gray-300 hover:text-primary-green transition-colors">Contact Us</a></li>
                <li><a href="#" className="text-gray-300 hover:text-primary-green transition-colors">Terms of Service</a></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold tracking-wider uppercase text-white">Follow Us</h4>
              <div className="mt-4 flex space-x-4">
                <a href="#" className="text-gray-300 hover:text-primary-green transition-colors"><Facebook className="h-5 w-5" /></a>
                <a href="#" className="text-gray-300 hover:text-primary-green transition-colors"><Instagram className="h-5 w-5" /></a>
                <a href="#" className="text-gray-300 hover:text-primary-green transition-colors"><Twitter className="h-5 w-5" /></a>
              </div>
            </div>
          </div>
          <div className="mt-8 border-t border-gray-700 pt-8 text-center text-sm text-gray-400">
            <p>&copy; 2024 Surplus Supper. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  )
}
