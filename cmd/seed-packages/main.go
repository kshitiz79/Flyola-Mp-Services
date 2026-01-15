package main

import (
	"encoding/json"
	"flyola-services/internal/config"
	"flyola-services/internal/database"
	"flyola-services/internal/models"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("✅ Database connected successfully")

	// Create sample holiday packages
	packages := []models.HolidayPackage{
		{
			Title:          "Maihar VIP Darshan - Jabalpur Return",
			Description:    "🔱 VIP Helicopter Tour to Maa Sharda Devi Temple with same day return from Jabalpur",
			PackageType:    "spiritual",
			DurationDays:   1,
			DurationNights: 0,
			PricePerPerson: 16000.00,
			MaxPassengers:  6,
			Status:         1,
			Inclusions:     createJSON([]string{
				"Jabalpur → Maihar → Jabalpur helicopter service",
				"Fast, safe & time-saving travel",
				"Maihar helipad to temple AC car transfer (to & fro)",
				"Maa Sharda Devi VIP Darshan (No queue, special arrangements)",
				"Special Prasad from temple",
				"Comfortable, secure & luxury journey",
				"Professional services & experienced staff",
			}),
			Exclusions: createJSON([]string{
				"Personal expenses",
				"Food and beverages",
				"Travel insurance",
			}),
			Itinerary: createJSON([]map[string]interface{}{
				{
					"day":        1,
					"title":      "Jabalpur to Maihar VIP Darshan",
					"duration":   "2 Hr 30 Min",
					"activities": []string{
						"Departure from Jabalpur",
						"Helicopter flight to Maihar",
						"AC car transfer to temple",
						"VIP Darshan at Maa Sharda Devi Temple",
						"Special Prasad collection",
						"Return helicopter flight to Jabalpur",
					},
				},
			}),
		},
		{
			Title:          "Maihar VIP Darshan - Chitrakoot",
			Description:    "🔱 VIP Helicopter Tour from Chitrakoot to Maa Sharda Temple",
			PackageType:    "spiritual",
			DurationDays:   1,
			DurationNights: 0,
			PricePerPerson: 10000.00,
			MaxPassengers:  6,
			Status:         1,
			Inclusions:     createJSON([]string{
				"Chitrakoot helicopter service",
				"Maa Sharda Temple VIP Darshan & special prasad",
				"Full AC luxury cab (arrival & departure)",
			}),
			Exclusions: createJSON([]string{
				"Personal expenses",
				"Food and beverages",
			}),
			Itinerary: createJSON([]map[string]interface{}{
				{
					"day":        1,
					"title":      "Chitrakoot to Maihar VIP Darshan",
					"duration":   "1 Hr 30 Min",
					"activities": []string{
						"Departure from Chitrakoot",
						"Helicopter flight to Maihar",
						"AC luxury cab to temple",
						"VIP Darshan at Maa Sharda Devi Temple",
						"Special Prasad collection",
						"Return journey",
					},
				},
			}),
		},
		{
			Title:          "Bandhavgarh Wildlife Safari",
			Description:    "🐅 1 Night / 2 Days Wildlife Helicopter Tour with jungle safari",
			PackageType:    "wildlife",
			DurationDays:   2,
			DurationNights: 1,
			PricePerPerson: 25000.00,
			MaxPassengers:  6,
			Status:         1,
			Inclusions:     createJSON([]string{
				"Helicopter travel (Jabalpur ⇄ Bandhavgarh)",
				"1 Night stay in Bandhavgarh",
				"All meals (Lunch, Dinner & Breakfast)",
				"Jungle Safari with all necessary permits",
				"Cab transfers (Helipad ⇄ Resort ⇄ Safari Gate)",
			}),
			Exclusions: createJSON([]string{
				"Personal expenses",
				"Alcoholic beverages",
				"Travel insurance",
			}),
			Itinerary: createJSON([]map[string]interface{}{
				{
					"day":        1,
					"title":      "Jabalpur to Bandhavgarh",
					"duration":   "Full Day",
					"activities": []string{
						"Helicopter departure from Jabalpur",
						"Arrival at Bandhavgarh",
						"Check-in at resort",
						"Lunch",
						"Evening jungle safari",
						"Dinner at resort",
					},
				},
				{
					"day":        2,
					"title":      "Bandhavgarh Safari & Return",
					"duration":   "Half Day",
					"activities": []string{
						"Early morning jungle safari",
						"Breakfast at resort",
						"Check-out",
						"Helicopter return to Jabalpur",
					},
				},
			}),
		},
	}

	// Insert packages
	for _, pkg := range packages {
		if err := db.Create(&pkg).Error; err != nil {
			log.Printf("Failed to create package %s: %v", pkg.Title, err)
		} else {
			log.Printf("✅ Created package: %s", pkg.Title)
		}
	}

	log.Println("🎉 Sample holiday packages seeded successfully!")
}

func createJSON(data interface{}) []byte {
	jsonData, _ := json.Marshal(data)
	return jsonData
}