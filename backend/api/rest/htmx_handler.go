package rest

import (
	"database/sql"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// HTMXHandler handles HTMX requests for server-side rendering
type HTMXHandler struct {
	db *sql.DB
}

// NewHTMXHandler creates a new HTMX handler
func NewHTMXHandler(db *sql.DB) *HTMXHandler {
	return &HTMXHandler{db: db}
}

// Restaurant represents a restaurant for the frontend
type Restaurant struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	CuisineType string  `json:"cuisine_type"`
	Rating      float64 `json:"rating"`
	Distance    float64 `json:"distance"`
}

// InventoryItem represents an inventory item for the frontend
type InventoryItem struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	OriginalPrice float64   `json:"original_price"`
	SurplusPrice  float64   `json:"surplus_price"`
	Quantity      int       `json:"quantity"`
	Category      string    `json:"category"`
	ExpiryTime    time.Time `json:"expiry_time"`
	Discount      float64   `json:"discount"`
}

// Offer represents an offer for the frontend
type Offer struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	OriginalPrice float64 `json:"original_price"`
	SurplusPrice  float64 `json:"surplus_price"`
	OfferType     string  `json:"offer_type"`
	Discount      float64 `json:"discount"`
}

// Order represents an order for the frontend
type Order struct {
	ID                  int       `json:"id"`
	RestaurantName      string    `json:"restaurant_name"`
	TotalAmount         float64   `json:"total_amount"`
	Status              string    `json:"status"`
	PickupTime          time.Time `json:"pickup_time"`
	SpecialInstructions string    `json:"special_instructions"`
	CreatedAt           time.Time `json:"created_at"`
}

// HomePageData represents data for the home page
type HomePageData struct {
	Title       string
	Description string
	SearchQuery string
}

// RestaurantListData represents data for the restaurant list page
type RestaurantListData struct {
	Restaurants []Restaurant
	SearchQuery string
	Latitude    float64
	Longitude   float64
}

// RestaurantDetailData represents data for the restaurant detail page
type RestaurantDetailData struct {
	Restaurant     Restaurant
	InventoryItems []InventoryItem
	Offers         []Offer
}

// DashboardData represents data for the restaurant dashboard
type DashboardData struct {
	Restaurant Restaurant
	Orders     []Order
	Stats      DashboardStats
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalOrders     int
	PendingOrders   int
	CompletedOrders int
	TotalRevenue    float64
}

// HandleHome handles the home page
func (h *HTMXHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	data := HomePageData{
		Title:       "Surplus Supper",
		Description: "Find delicious surplus food from restaurants near you",
		SearchQuery: r.URL.Query().Get("q"),
	}

	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>{{.Title}}</title>
		<script src="https://unpkg.com/htmx.org@1.9.6"></script>
		<script src="https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js" defer></script>
		<script src="https://cdn.tailwindcss.com"></script>
	</head>
	<body class="bg-gray-50">
		<div class="min-h-screen">
			<!-- Header -->
			<header class="bg-green-600 text-white shadow-lg">
				<div class="container mx-auto px-4 py-6">
					<div class="flex items-center justify-between">
						<h1 class="text-3xl font-bold">🍽️ Surplus Supper</h1>
						<nav class="space-x-4">
							<a href="/" class="hover:text-green-200">Home</a>
							<a href="/restaurants" class="hover:text-green-200">Restaurants</a>
							<a href="/restaurant/login" class="hover:text-green-200">Restaurant Login</a>
						</nav>
					</div>
				</div>
			</header>

			<!-- Hero Section -->
			<section class="bg-gradient-to-r from-green-500 to-green-600 text-white py-20">
				<div class="container mx-auto px-4 text-center">
					<h2 class="text-5xl font-bold mb-6">Reduce Food Waste, Save Money</h2>
					<p class="text-xl mb-8">Discover delicious surplus food from restaurants near you at amazing discounts.</p>
					
					<!-- Location and Search Section -->
					<div class="max-w-md mx-auto space-y-4">
						<!-- Current Location Display -->
						<div id="location-display" class="bg-green-700 rounded-lg p-3 hidden">
							<p class="text-sm">📍 <span id="current-location">Detecting your location...</span></p>
						</div>
						
						<!-- Search Form -->
						<form hx-get="/restaurants" hx-target="#restaurants-list" hx-trigger="submit" class="flex">
							<input 
								type="text" 
								name="location" 
								id="location-input"
								placeholder="Enter your location or use current location..." 
								class="flex-1 px-4 py-3 rounded-l-lg text-gray-900 focus:outline-none focus:ring-2 focus:ring-green-400"
								value="{{.SearchQuery}}"
							>
							<button 
								type="submit" 
								class="bg-green-700 hover:bg-green-800 px-6 py-3 rounded-r-lg font-semibold transition-colors"
							>
								Search
							</button>
						</form>
						
						<!-- Location Buttons -->
						<div class="flex space-x-2">
							<button 
								type="button" 
								onclick="getCurrentLocation()"
								class="flex-1 bg-green-700 hover:bg-green-800 px-4 py-2 rounded-lg font-semibold transition-colors"
							>
								📍 Use My Location
							</button>
							<button 
								type="button" 
								onclick="searchNearby()"
								class="flex-1 bg-green-700 hover:bg-green-800 px-4 py-2 rounded-lg font-semibold transition-colors"
							>
								🌍 Find Nearby
							</button>
						</div>
					</div>
				</div>
			</section>

			<!-- Features Section -->
			<section class="py-16 bg-white">
				<div class="container mx-auto px-4">
					<h3 class="text-3xl font-bold text-center mb-12 text-gray-800">How It Works</h3>
					<div class="grid md:grid-cols-3 gap-8">
						<div class="text-center">
							<div class="bg-green-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
								<span class="text-2xl">🔍</span>
							</div>
							<h4 class="text-xl font-semibold mb-2">Find Restaurants</h4>
							<p class="text-gray-600">Discover restaurants with surplus food near your location.</p>
						</div>
						<div class="text-center">
							<div class="bg-green-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
								<span class="text-2xl">💰</span>
							</div>
							<h4 class="text-xl font-semibold mb-2">Save Money</h4>
							<p class="text-gray-600">Get delicious food at 40-70% off regular prices.</p>
						</div>
						<div class="text-center">
							<div class="bg-green-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
								<span class="text-2xl">🌍</span>
							</div>
							<h4 class="text-xl font-semibold mb-2">Help the Planet</h4>
							<p class="text-gray-600">Reduce food waste and help save the environment.</p>
						</div>
					</div>
				</div>
			</section>

			<!-- Restaurants List -->
			<section class="py-16 bg-gray-50">
				<div class="container mx-auto px-4">
					<h3 class="text-3xl font-bold text-center mb-12 text-gray-800">Nearby Restaurants</h3>
					<div id="restaurants-list">
						<!-- Restaurants will be loaded here via HTMX -->
						<div class="text-center text-gray-500">
							<p>Use the buttons above to find restaurants with surplus food near you.</p>
						</div>
					</div>
				</div>
			</section>

			<!-- Footer -->
			<footer class="bg-gray-800 text-white py-8">
				<div class="container mx-auto px-4 text-center">
					<p>&copy; 2024 Surplus Supper. All rights reserved.</p>
				</div>
			</footer>
		</div>

		<script>
			let userLatitude = null;
			let userLongitude = null;
			let userLocation = null;

			// Get user's current location
			function getCurrentLocation() {
				if (navigator.geolocation) {
					document.getElementById('location-display').classList.remove('hidden');
					document.getElementById('current-location').textContent = 'Getting your location...';
					
					navigator.geolocation.getCurrentPosition(
						function(position) {
							userLatitude = position.coords.latitude;
							userLongitude = position.coords.longitude;
							
							// Reverse geocode to get address
							reverseGeocode(userLatitude, userLongitude);
							
							// Search for restaurants with current location
							searchWithCoordinates(userLatitude, userLongitude);
						},
						function(error) {
							console.error('Error getting location:', error);
							document.getElementById('current-location').textContent = 'Location access denied. Please enter your location manually.';
							document.getElementById('location-display').classList.add('bg-red-600');
						}
					);
				} else {
					alert('Geolocation is not supported by this browser. Please enter your location manually.');
				}
			}

			// Reverse geocode coordinates to address
			function reverseGeocode(lat, lng) {
				fetch('https://nominatim.openstreetmap.org/reverse?format=json&lat=' + lat + '&lon=' + lng + '&zoom=10')
					.then(function(response) { return response.json(); })
					.then(function(data) {
						userLocation = data.display_name;
						document.getElementById('current-location').textContent = userLocation;
						document.getElementById('location-input').value = userLocation;
					})
					.catch(function(error) {
						console.error('Error reverse geocoding:', error);
						document.getElementById('current-location').textContent = 'Location: ' + lat.toFixed(4) + ', ' + lng.toFixed(4);
					});
			}

			// Search with coordinates
			function searchWithCoordinates(lat, lng) {
				var form = document.querySelector('form[hx-get="/restaurants"]');
				var input = document.getElementById('location-input');
				
				// Add coordinates as hidden fields
				var coordInput = document.getElementById('lat-input');
				if (!coordInput) {
					coordInput = document.createElement('input');
					coordInput.type = 'hidden';
					coordInput.name = 'lat';
					coordInput.id = 'lat-input';
					form.appendChild(coordInput);
				}
				
				var lngInput = document.getElementById('lng-input');
				if (!lngInput) {
					lngInput = document.createElement('input');
					lngInput.type = 'hidden';
					lngInput.name = 'lng';
					lngInput.id = 'lng-input';
					form.appendChild(lngInput);
				}
				
				coordInput.value = lat;
				lngInput.value = lng;
				
				// Trigger HTMX request
				htmx.trigger(form, 'submit');
			}

			// Search nearby (uses current location if available)
			function searchNearby() {
				if (userLatitude && userLongitude) {
					searchWithCoordinates(userLatitude, userLongitude);
				} else {
					getCurrentLocation();
				}
			}

			// Auto-detect location on page load (optional)
			window.addEventListener('load', function() {
				// Uncomment the line below to auto-detect location on page load
				// getCurrentLocation();
			});
		</script>
	</body>
	</html>
	`

	tmplParsed, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmplParsed.Execute(w, data)
}

// HandleRestaurantList handles the restaurant list page
func (h *HTMXHandler) HandleRestaurantList(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	// Use provided coordinates or default to New York
	latitude, longitude := 40.7128, -74.0060 // Default: New York coordinates

	if latStr != "" && lngStr != "" {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			latitude = lat
		}
		if lng, err := strconv.ParseFloat(lngStr, 64); err == nil {
			longitude = lng
		}
		log.Printf("Using coordinates: %f, %f (Location: %s)", latitude, longitude, location)
	} else if location != "" {
		log.Printf("Location provided: %s, using default coordinates", location)
	}

	// Get nearby restaurants with distance calculation
	query := `
		SELECT id, name, description, address, latitude, longitude, phone, email, cuisine_type, rating,
		       (6371 * acos(
		           cos(radians($1)) * cos(radians(latitude)) * cos(radians(longitude) - radians($2)) + 
		           sin(radians($1)) * sin(radians(latitude))
		       )) as distance
		FROM restaurants 
		WHERE is_active = true 
		ORDER BY distance
		LIMIT 10
	`

	rows, err := h.db.Query(query, latitude, longitude)
	if err != nil {
		log.Printf("Database error: %v", err)
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<div class="text-center text-red-600 p-4">Failed to fetch restaurants. Please try again.</div>`))
		} else {
			http.Error(w, "Failed to fetch restaurants", http.StatusInternalServerError)
		}
		return
	}
	defer rows.Close()

	var restaurants []Restaurant
	for rows.Next() {
		var restaurant Restaurant
		var distance float64
		err := rows.Scan(
			&restaurant.ID, &restaurant.Name, &restaurant.Description, &restaurant.Address,
			&restaurant.Latitude, &restaurant.Longitude, &restaurant.Phone, &restaurant.Email,
			&restaurant.CuisineType, &restaurant.Rating, &distance,
		)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}

		restaurant.Distance = distance
		restaurants = append(restaurants, restaurant)
	}

	log.Printf("Found %d restaurants near coordinates %f, %f", len(restaurants), latitude, longitude)

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Return just the restaurants list
		tmpl := `
		{{if .Restaurants}}
			{{range .Restaurants}}
			<div class="bg-white rounded-lg shadow-md p-6 mb-4 hover:shadow-lg transition-shadow">
				<div class="flex items-center justify-between">
					<div class="flex-1">
						<h4 class="text-xl font-semibold text-gray-800 mb-2">{{.Name}}</h4>
						<p class="text-gray-600 mb-2">{{.Description}}</p>
						<div class="flex items-center space-x-4 text-sm text-gray-500">
							<span> {{.Address}}</span>
							<span>🍽️ {{.CuisineType}}</span>
							<span>⭐ {{.Rating}}</span>
							<span>📏 {{printf "%.1f" .Distance}} km</span>
						</div>
					</div>
					<div class="ml-4">
						<a 
							href="/restaurant/{{.ID}}" 
							class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-lg font-semibold transition-colors"
						>
							View Offers
						</a>
					</div>
				</div>
			</div>
			{{end}}
		{{else}}
			<div class="text-center p-8">
				<div class="bg-gray-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
					<span class="text-2xl"></span>
				</div>
				<h3 class="text-xl font-semibold text-gray-800 mb-2">No Restaurants Found</h3>
				<p class="text-gray-600 mb-4">Sorry, there are no restaurants with surplus food in your area yet.</p>
				<div class="space-y-2 text-sm text-gray-500">
					<p>📍 Location: {{if .SearchQuery}}{{.SearchQuery}}{{else}}Your current location{{end}}</p>
					<p>📏 Search radius: 10 km</p>
				</div>
				<div class="mt-6">
					<button onclick="getCurrentLocation()" class="bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded-lg font-semibold transition-colors mr-2">
						📍 Try My Location
					</button>
					<button onclick="document.getElementById('location-input').focus()" class="bg-gray-600 hover:bg-gray-700 text-white px-6 py-2 rounded-lg font-semibold transition-colors">
						🔍 Search Different Area
					</button>
				</div>
			</div>
		{{end}}
		`

		tmplParsed, err := template.New("restaurants").Parse(tmpl)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		tmplParsed.Execute(w, RestaurantListData{
			Restaurants: restaurants,
			SearchQuery: location,
			Latitude:    latitude,
			Longitude:   longitude,
		})
		return
	}

	// Full page request
	// Return full page template
	h.HandleHome(w, r)
}

// HandleRestaurantDetail handles the restaurant detail page
func (h *HTMXHandler) HandleRestaurantDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	restaurantID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	// Get restaurant details
	var restaurant Restaurant
	err = h.db.QueryRow(`
		SELECT id, name, description, address, latitude, longitude, phone, email, cuisine_type, rating
		FROM restaurants WHERE id = $1 AND is_active = true
	`, restaurantID).Scan(
		&restaurant.ID, &restaurant.Name, &restaurant.Description, &restaurant.Address,
		&restaurant.Latitude, &restaurant.Longitude, &restaurant.Phone, &restaurant.Email,
		&restaurant.CuisineType, &restaurant.Rating,
	)
	if err != nil {
		http.Error(w, "Restaurant not found", http.StatusNotFound)
		return
	}

	// Get inventory items
	rows, err := h.db.Query(`
		SELECT id, name, description, original_price, surplus_price, quantity, category, expiry_time
		FROM inventory_items 
		WHERE restaurant_id = $1 AND is_available = true AND expiry_time > NOW()
		ORDER BY created_at DESC
	`, restaurantID)
	if err != nil {
		http.Error(w, "Failed to fetch inventory", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var inventoryItems []InventoryItem
	for rows.Next() {
		var item InventoryItem
		err := rows.Scan(
			&item.ID, &item.Name, &item.Description, &item.OriginalPrice, &item.SurplusPrice,
			&item.Quantity, &item.Category, &item.ExpiryTime,
		)
		if err != nil {
			continue
		}
		item.Discount = ((item.OriginalPrice - item.SurplusPrice) / item.OriginalPrice) * 100
		inventoryItems = append(inventoryItems, item)
	}

	// Get offers
	offerRows, err := h.db.Query(`
		SELECT id, name, description, original_price, surplus_price, offer_type
		FROM offers 
		WHERE restaurant_id = $1 AND is_available = true
		ORDER BY created_at DESC
	`, restaurantID)
	if err != nil {
		http.Error(w, "Failed to fetch offers", http.StatusInternalServerError)
		return
	}
	defer offerRows.Close()

	var offers []Offer
	for offerRows.Next() {
		var offer Offer
		err := offerRows.Scan(
			&offer.ID, &offer.Name, &offer.Description, &offer.OriginalPrice, &offer.SurplusPrice, &offer.OfferType,
		)
		if err != nil {
			continue
		}
		offer.Discount = ((offer.OriginalPrice - offer.SurplusPrice) / offer.OriginalPrice) * 100
		offers = append(offers, offer)
	}

	data := RestaurantDetailData{
		Restaurant:     restaurant,
		InventoryItems: inventoryItems,
		Offers:         offers,
	}

	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>{{.Restaurant.Name}} - Surplus Supper</title>
		<script src="https://unpkg.com/htmx.org@1.9.6"></script>
		<script src="https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js" defer></script>
		<script src="https://cdn.tailwindcss.com"></script>
	</head>
	<body class="bg-gray-50">
		<div class="min-h-screen">
			<!-- Header -->
			<header class="bg-green-600 text-white shadow-lg">
				<div class="container mx-auto px-4 py-6">
					<div class="flex items-center justify-between">
						<h1 class="text-3xl font-bold">🍽️ Surplus Supper</h1>
						<nav class="space-x-4">
							<a href="/" class="hover:text-green-200">Home</a>
							<a href="/restaurants" class="hover:text-green-200">Restaurants</a>
						</nav>
					</div>
				</div>
			</header>

			<!-- Restaurant Info -->
			<section class="bg-white py-8">
				<div class="container mx-auto px-4">
					<div class="bg-white rounded-lg shadow-md p-6">
						<h2 class="text-3xl font-bold text-gray-800 mb-4">{{.Restaurant.Name}}</h2>
						<p class="text-gray-600 mb-4">{{.Restaurant.Description}}</p>
						<div class="flex items-center space-x-6 text-sm text-gray-500">
							<span>📍 {{.Restaurant.Address}}</span>
							<span>🍽️ {{.Restaurant.CuisineType}}</span>
							<span>⭐ {{.Restaurant.Rating}}</span>
							<span>📞 {{.Restaurant.Phone}}</span>
						</div>
					</div>
				</div>
			</section>

			<!-- Inventory Items -->
			<section class="py-8">
				<div class="container mx-auto px-4">
					<h3 class="text-2xl font-bold text-gray-800 mb-6">Available Items</h3>
					<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
						{{range .InventoryItems}}
						<div class="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
							<div class="flex justify-between items-start mb-4">
								<h4 class="text-lg font-semibold text-gray-800">{{.Name}}</h4>
								<span class="bg-red-100 text-red-800 px-2 py-1 rounded-full text-sm font-semibold">
									-{{printf "%.0f" .Discount}}%
								</span>
							</div>
							<p class="text-gray-600 mb-4">{{.Description}}</p>
							<div class="flex justify-between items-center mb-4">
								<div>
									<span class="text-gray-500 line-through">${{printf "%.2f" .OriginalPrice}}</span>
									<span class="text-2xl font-bold text-green-600 ml-2">${{printf "%.2f" .SurplusPrice}}</span>
								</div>
								<span class="text-sm text-gray-500">Qty: {{.Quantity}}</span>
							</div>
							<div class="flex justify-between items-center text-sm text-gray-500 mb-4">
								<span>Category: {{.Category}}</span>
								<span>Expires: {{.ExpiryTime.Format "15:04"}}</span>
							</div>
							<button 
								class="w-full bg-green-600 hover:bg-green-700 text-white py-2 px-4 rounded-lg font-semibold transition-colors"
								onclick="addToCart({{.ID}}, 'inventory', {{.SurplusPrice}})"
							>
								Add to Cart
							</button>
						</div>
						{{end}}
					</div>
				</div>
			</section>

			<!-- Offers -->
			<section class="py-8 bg-gray-50">
				<div class="container mx-auto px-4">
					<h3 class="text-2xl font-bold text-gray-800 mb-6">Special Offers</h3>
					<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
						{{range .Offers}}
						<div class="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
							<div class="flex justify-between items-start mb-4">
								<h4 class="text-lg font-semibold text-gray-800">{{.Name}}</h4>
								<span class="bg-purple-100 text-purple-800 px-2 py-1 rounded-full text-sm font-semibold">
									{{.OfferType}}
								</span>
							</div>
							<p class="text-gray-600 mb-4">{{.Description}}</p>
							<div class="flex justify-between items-center mb-4">
								<div>
									<span class="text-gray-500 line-through">${{printf "%.2f" .OriginalPrice}}</span>
									<span class="text-2xl font-bold text-green-600 ml-2">${{printf "%.2f" .SurplusPrice}}</span>
								</div>
								<span class="bg-red-100 text-red-800 px-2 py-1 rounded-full text-sm font-semibold">
									-{{printf "%.0f" .Discount}}%
								</span>
							</div>
							<button 
								class="w-full bg-purple-600 hover:bg-purple-700 text-white py-2 px-4 rounded-lg font-semibold transition-colors"
								onclick="addToCart({{.ID}}, 'offer', {{.SurplusPrice}})"
							>
								Add to Cart
							</button>
						</div>
						{{end}}
					</div>
				</div>
			</section>

			<!-- Shopping Cart -->
			<div id="cart" class="fixed bottom-4 right-4 bg-white rounded-lg shadow-lg p-4 hidden">
				<h4 class="font-semibold mb-2">Shopping Cart</h4>
				<div id="cart-items"></div>
				<div class="flex justify-between items-center mt-4">
					<span class="font-semibold">Total: $<span id="cart-total">0.00</span></span>
					<button 
						class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-lg font-semibold transition-colors"
						onclick="checkout()"
					>
						Checkout
					</button>
				</div>
			</div>

			<script>
				let cart = [];
				let cartTotal = 0;

				function addToCart(id, type, price) {
					cart.push({id, type, price});
					cartTotal += price;
					updateCart();
					showCart();
				}

				function updateCart() {
					document.getElementById('cart-total').textContent = cartTotal.toFixed(2);
					document.getElementById('cart-items').innerHTML = cart.map(function(item) {
						return '<div class="text-sm">' + (item.type === 'inventory' ? 'Item' : 'Offer') + ' #' + item.id + ': $' + item.price.toFixed(2) + '</div>';
					}).join('');
				}

				function showCart() {
					document.getElementById('cart').classList.remove('hidden');
				}

				function checkout() {
					if (cart.length === 0) return;
					
					const orderData = {
						restaurant_id: {{.Restaurant.ID}},
						items: cart
					};

					fetch('/order/confirm', {
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
						},
						body: JSON.stringify(orderData)
					})
					.then(response => response.text())
					.then(html => {
						document.body.innerHTML = html;
					})
					.catch(error => {
						console.error('Error:', error);
						alert('Failed to place order. Please try again.');
					});
				}
			</script>
		</div>
	</body>
	</html>
	`

	tmplParsed, err := template.New("restaurant-detail").Parse(tmpl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmplParsed.Execute(w, data)
}

// HandleRestaurantLogin handles restaurant login/registration
func (h *HTMXHandler) HandleRestaurantLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Handle login form submission
		email := r.FormValue("email")
		password := r.FormValue("password")
		action := r.FormValue("action") // "login" or "register"

		if action == "register" {
			// Handle restaurant registration
			h.handleRestaurantRegistration(w, r, email, password)
		} else {
			// Handle restaurant login
			h.handleRestaurantLogin(w, r, email, password)
		}
		return
	}

	// Serve login page
	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Restaurant Login - Surplus Supper</title>
		<script src="https://unpkg.com/htmx.org@1.9.6"></script>
		<script src="https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js" defer></script>
		<script src="https://cdn.tailwindcss.com"></script>
	</head>
	<body class="bg-gray-50">
		<div class="min-h-screen flex items-center justify-center">
			<div class="max-w-md w-full space-y-8">
				<div class="text-center">
					<h2 class="text-3xl font-bold text-gray-900">Restaurant Portal</h2>
					<p class="mt-2 text-gray-600">Manage your surplus food inventory</p>
				</div>
				
				<div class="bg-white rounded-lg shadow-lg p-8">
					<div x-data="{ mode: 'login' }">
						<!-- Toggle Buttons -->
						<div class="flex rounded-lg bg-gray-100 p-1 mb-6">
							<button 
								@click="mode = 'login'" 
								:class="mode === 'login' ? 'bg-white shadow-sm' : ''"
								class="flex-1 py-2 px-4 rounded-md text-sm font-medium transition-colors"
							>
								Login
							</button>
							<button 
								@click="mode = 'register'" 
								:class="mode === 'register' ? 'bg-white shadow-sm' : ''"
								class="flex-1 py-2 px-4 rounded-md text-sm font-medium transition-colors"
							>
								Register
							</button>
						</div>

						<!-- Login Form -->
						<form x-show="mode === 'login'" method="POST" class="space-y-4">
							<input type="hidden" name="action" value="login">
							<div>
								<label class="block text-sm font-medium text-gray-700">Email</label>
								<input type="email" name="email" required class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-green-500 focus:border-green-500">
							</div>
							<div>
								<label class="block text-sm font-medium text-gray-700">Password</label>
								<input type="password" name="password" required class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-green-500 focus:border-green-500">
							</div>
							<button type="submit" class="w-full bg-green-600 hover:bg-green-700 text-white py-2 px-4 rounded-md font-medium transition-colors">
								Login
							</button>
						</form>

						<!-- Registration Form -->
						<form x-show="mode === 'register'" method="POST" class="space-y-4">
							<input type="hidden" name="action" value="register">
							<div>
								<label class="block text-sm font-medium text-gray-700">Restaurant Name</label>
								<input type="text" name="restaurant_name" required class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-green-500 focus:border-green-500">
							</div>
							<div>
								<label class="block text-sm font-medium text-gray-700">Email</label>
								<input type="email" name="email" required class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-green-500 focus:border-green-500">
							</div>
							<div>
								<label class="block text-sm font-medium text-gray-700">Password</label>
								<input type="password" name="password" required class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-green-500 focus:border-green-500">
							</div>
							<div>
								<label class="block text-sm font-medium text-gray-700">Address</label>
								<input type="text" name="address" required class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-green-500 focus:border-green-500">
							</div>
							<div>
								<label class="block text-sm font-medium text-gray-700">Cuisine Type</label>
								<select name="cuisine_type" required class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-green-500 focus:border-green-500">
									<option value="">Select cuisine type</option>
									<option value="Italian">Italian</option>
									<option value="Japanese">Japanese</option>
									<option value="American">American</option>
									<option value="Mexican">Mexican</option>
									<option value="Chinese">Chinese</option>
									<option value="Indian">Indian</option>
									<option value="Thai">Thai</option>
									<option value="Healthy">Healthy</option>
									<option value="Other">Other</option>
								</select>
							</div>
							<button type="submit" class="w-full bg-green-600 hover:bg-green-700 text-white py-2 px-4 rounded-md font-medium transition-colors">
								Register Restaurant
							</button>
						</form>
					</div>
				</div>

				<div class="text-center">
					<a href="/" class="text-green-600 hover:text-green-500">← Back to Home</a>
				</div>
			</div>
		</div>
	</body>
	</html>
	`

	tmplParsed, err := template.New("restaurant-login").Parse(tmpl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmplParsed.Execute(w, nil)
}

// Handle restaurant registration
func (h *HTMXHandler) handleRestaurantRegistration(w http.ResponseWriter, r *http.Request, email, password string) {
	restaurantName := r.FormValue("restaurant_name")
	address := r.FormValue("address")
	cuisineType := r.FormValue("cuisine_type")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	// Insert restaurant
	query := `
		INSERT INTO restaurants (name, email, cuisine_type, address, latitude, longitude, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, true)
		RETURNING id
	`

	// For demo, use default coordinates (in real app, geocode the address)
	latitude, longitude := 40.7128, -74.0060

	var restaurantID int
	err = h.db.QueryRow(query, restaurantName, email, cuisineType, address, latitude, longitude).Scan(&restaurantID)
	if err != nil {
		log.Printf("Registration error: %v", err)
		http.Error(w, "Registration failed - restaurant may already exist", http.StatusBadRequest)
		return
	}

	// Insert restaurant staff
	staffQuery := `
		INSERT INTO restaurant_staff (restaurant_id, email, password_hash, role)
		VALUES ($1, $2, $3, 'owner')
	`

	_, err = h.db.Exec(staffQuery, restaurantID, email, string(hashedPassword))
	if err != nil {
		log.Printf("Staff creation error: %v", err)
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	// Redirect to dashboard
	http.Redirect(w, r, "/restaurant/dashboard", http.StatusSeeOther)
}

// Handle restaurant login
func (h *HTMXHandler) handleRestaurantLogin(w http.ResponseWriter, r *http.Request, email, password string) {
	// Query for restaurant staff
	query := `
		SELECT rs.password_hash, r.id, r.name
		FROM restaurant_staff rs
		JOIN restaurants r ON rs.restaurant_id = r.id
		WHERE rs.email = $1
	`

	var hashedPassword, restaurantName string
	var restaurantID int
	err := h.db.QueryRow(query, email).Scan(&hashedPassword, &restaurantID, &restaurantName)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Set session (in real app, use proper session management)
	// For now, redirect to dashboard
	http.Redirect(w, r, "/restaurant/dashboard", http.StatusSeeOther)
}

// HandleRestaurantDashboard handles the restaurant dashboard
func (h *HTMXHandler) HandleRestaurantDashboard(w http.ResponseWriter, r *http.Request) {
	// In a real app, get restaurant ID from session
	// For demo, we'll show a sample dashboard

	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Restaurant Dashboard - Surplus Supper</title>
		<script src="https://unpkg.com/htmx.org@1.9.6"></script>
		<script src="https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js" defer></script>
		<script src="https://cdn.tailwindcss.com"></script>
	</head>
	<body class="bg-gray-50">
		<div class="min-h-screen">
			<!-- Header -->
			<header class="bg-green-600 text-white shadow-lg">
				<div class="container mx-auto px-4 py-6">
					<div class="flex items-center justify-between">
						<h1 class="text-3xl font-bold">🏪 Restaurant Dashboard</h1>
						<nav class="space-x-4">
							<a href="/" class="hover:text-green-200">Home</a>
							<a href="/restaurant/inventory" class="hover:text-green-200">Manage Inventory</a>
							<a href="/restaurant/orders" class="hover:text-green-200">Orders</a>
						</nav>
					</div>
				</div>
			</header>

			<!-- Dashboard Content -->
			<div class="container mx-auto px-4 py-8">
				<div class="grid md:grid-cols-3 gap-6 mb-8">
					<!-- Stats Cards -->
					<div class="bg-white rounded-lg shadow-md p-6">
						<div class="flex items-center">
							<div class="bg-green-100 p-3 rounded-full">
								<span class="text-2xl">📦</span>
							</div>
							<div class="ml-4">
								<h3 class="text-lg font-semibold text-gray-800">Active Items</h3>
								<p class="text-3xl font-bold text-green-600">12</p>
							</div>
						</div>
					</div>

					<div class="bg-white rounded-lg shadow-md p-6">
						<div class="flex items-center">
							<div class="bg-blue-100 p-3 rounded-full">
								<span class="text-2xl">📋</span>
							</div>
							<div class="ml-4">
								<h3 class="text-lg font-semibold text-gray-800">Pending Orders</h3>
								<p class="text-3xl font-bold text-blue-600">5</p>
							</div>
						</div>
					</div>

					<div class="bg-white rounded-lg shadow-md p-6">
						<div class="flex items-center">
							<div class="bg-purple-100 p-3 rounded-full">
								<span class="text-2xl">💰</span>
							</div>
							<div class="ml-4">
								<h3 class="text-lg font-semibold text-gray-800">Today's Revenue</h3>
								<p class="text-3xl font-bold text-purple-600">$245</p>
							</div>
						</div>
					</div>
				</div>

				<!-- Quick Actions -->
				<div class="bg-white rounded-lg shadow-md p-6 mb-8">
					<h2 class="text-2xl font-bold text-gray-800 mb-4">Quick Actions</h2>
					<div class="grid md:grid-cols-2 gap-4">
						<a href="/restaurant/inventory" class="bg-green-600 hover:bg-green-700 text-white p-4 rounded-lg text-center font-semibold transition-colors">
							➕ Add Inventory Item
						</a>
						<a href="/restaurant/orders" class="bg-blue-600 hover:bg-blue-700 text-white p-4 rounded-lg text-center font-semibold transition-colors">
							 View Orders
						</a>
					</div>
				</div>

				<!-- Recent Activity -->
				<div class="bg-white rounded-lg shadow-md p-6">
					<h2 class="text-2xl font-bold text-gray-800 mb-4">Recent Activity</h2>
					<div class="space-y-4">
						<div class="flex items-center p-3 bg-gray-50 rounded-lg">
							<div class="bg-green-100 p-2 rounded-full mr-3">
								<span>📦</span>
							</div>
							<div>
								<p class="font-semibold">New order received</p>
								<p class="text-sm text-gray-600">Order #1234 - 2 items</p>
							</div>
							<span class="ml-auto text-sm text-gray-500">2 min ago</span>
						</div>
						<div class="flex items-center p-3 bg-gray-50 rounded-lg">
							<div class="bg-yellow-100 p-2 rounded-full mr-3">
								<span>⚠️</span>
							</div>
							<div>
								<p class="font-semibold">Item expiring soon</p>
								<p class="text-sm text-gray-600">Margherita Pizza - expires in 1 hour</p>
							</div>
							<span class="ml-auto text-sm text-gray-500">15 min ago</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</body>
	</html>
	`

	tmplParsed, err := template.New("dashboard").Parse(tmpl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmplParsed.Execute(w, nil)
}

// calculateDistance calculates the distance between two points using Haversine formula
func (h *HTMXHandler) calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in kilometers

	lat1Rad := lat1 * 3.14159265359 / 180
	lon1Rad := lon1 * 3.14159265359 / 180
	lat2Rad := lat2 * 3.14159265359 / 180
	lon2Rad := lon2 * 3.14159265359 / 180

	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := (dlat/2)*(dlat/2) + (dlon/2)*(dlon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
