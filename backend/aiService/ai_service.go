package aiService

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// PriceRecommendation represents a price recommendation from the Profit Advisor
type PriceRecommendation struct {
	ItemID           string  `json:"item_id"`
	OriginalPrice    float64 `json:"original_price"`
	RecommendedPrice float64 `json:"recommended_price"`
	ConfidenceScore  float64 `json:"confidence_score"`
	Reasoning        string  `json:"reasoning"`
}

// Recipe represents a generated recipe from Creative Kitchen
type Recipe struct {
	Name        string   `json:"name"`
	Ingredients []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	PrepTime    string   `json:"prep_time"`
	CookTime    string   `json:"cook_time"`
	Difficulty  string   `json:"difficulty"`
	Tags        []string `json:"tags"`
}

// GetPriceRecommendation takes an item ID and returns a suggested discount price
// This is a stub function that returns hardcoded values for now
func GetPriceRecommendation(itemID string) (*PriceRecommendation, error) {
	// Simulate AI processing time
	time.Sleep(100 * time.Millisecond)
	
	// For now, return a hardcoded recommendation
	// In a real implementation, this would call an AI service
	originalPrice := 15.99 + rand.Float64()*10.0
	discountPercentage := 0.4 + rand.Float64()*0.3 // 40-70% discount
	recommendedPrice := originalPrice * (1 - discountPercentage)
	
	recommendation := &PriceRecommendation{
		ItemID:           itemID,
		OriginalPrice:    originalPrice,
		RecommendedPrice: recommendedPrice,
		ConfidenceScore:  0.85 + rand.Float64()*0.1,
		Reasoning:        "Based on market analysis, competitor pricing, and demand patterns, this price optimizes for quick sale while maintaining profitability.",
	}
	
	log.Printf("Generated price recommendation for item %s: $%.2f (original: $%.2f)", 
		itemID, recommendation.RecommendedPrice, recommendation.OriginalPrice)
	
	return recommendation, nil
}

// GenerateRecipe takes a list of surplus ingredients and user preference
// and prepares a structured prompt for an LLM to generate a recipe
func GenerateRecipe(ingredients string, preference string) (string, error) {
	// Create a structured prompt for the LLM
	prompt := fmt.Sprintf(`Generate a creative recipe using the following surplus ingredients: %s

User Preference: %s

Please provide the recipe in the following JSON format:
{
  "name": "Recipe Name",
  "ingredients": ["ingredient1", "ingredient2", ...],
  "instructions": ["step1", "step2", ...],
  "prep_time": "X minutes",
  "cook_time": "X minutes", 
  "difficulty": "easy/medium/hard",
  "tags": ["tag1", "tag2", ...]
}

Requirements:
- Use all the provided ingredients
- Consider the user preference: %s
- Make it creative and appealing
- Include clear, step-by-step instructions
- Suggest appropriate cooking times and difficulty level
- Add relevant tags for categorization

Please respond with only the JSON object, no additional text.`, ingredients, preference, preference)
	
	log.Printf("Generated recipe prompt for ingredients: %s, preference: %s", ingredients, preference)
	
	return prompt, nil
}

// ProcessRecipeResponse processes the LLM response and converts it to a Recipe struct
func ProcessRecipeResponse(llmResponse string) (*Recipe, error) {
	var recipe Recipe
	err := json.Unmarshal([]byte(llmResponse), &recipe)
	if err != nil {
		return nil, fmt.Errorf("failed to parse recipe response: %w", err)
	}
	
	return &recipe, nil
}

// GetBulkPriceRecommendations generates price recommendations for multiple items
func GetBulkPriceRecommendations(itemIDs []string) ([]*PriceRecommendation, error) {
	var recommendations []*PriceRecommendation
	
	for _, itemID := range itemIDs {
		recommendation, err := GetPriceRecommendation(itemID)
		if err != nil {
			log.Printf("Failed to get recommendation for item %s: %v", itemID, err)
			continue
		}
		recommendations = append(recommendations, recommendation)
	}
	
	return recommendations, nil
} 