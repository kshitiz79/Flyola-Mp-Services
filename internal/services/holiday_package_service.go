package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"flyola-services/internal/models"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type HolidayPackageService struct {
	db              *gorm.DB
	nodeBackendURL  string
}

func NewHolidayPackageService(db *gorm.DB, nodeBackendURL string) *HolidayPackageService {
	return &HolidayPackageService{
		db:             db,
		nodeBackendURL: nodeBackendURL,
	}
}

// GetAllPackages retrieves all active holiday packages
func (s *HolidayPackageService) GetAllPackages() ([]models.HolidayPackage, error) {
	var packages []models.HolidayPackage
	err := s.db.Where("status = ?", 1).
		Preload("PackageSchedules").
		Find(&packages).Error
	return packages, err
}

// GetPackageByID retrieves a specific package with all details
func (s *HolidayPackageService) GetPackageByID(id uint) (*models.HolidayPackage, error) {
	var pkg models.HolidayPackage
	err := s.db.Where("id = ? AND status = ?", id, 1).
		Preload("PackageSchedules").
		First(&pkg).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

// GetPackagesByType retrieves packages by type (spiritual, wildlife, etc.)
func (s *HolidayPackageService) GetPackagesByType(packageType string) ([]models.HolidayPackage, error) {
	var packages []models.HolidayPackage
	err := s.db.Where("package_type = ? AND status = ?", packageType, 1).
		Preload("PackageSchedules").
		Find(&packages).Error
	return packages, err
}

// CreatePackageBooking creates a new package booking
func (s *HolidayPackageService) CreatePackageBooking(booking *models.PackageBooking) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Create the main booking
		if err := tx.Create(booking).Error; err != nil {
			return err
		}

		// Get package details with schedules
		var pkg models.HolidayPackage
		if err := tx.Preload("PackageSchedules").First(&pkg, booking.PackageID).Error; err != nil {
			return err
		}

		// Create schedule bookings for each schedule in the package
		for _, schedule := range pkg.PackageSchedules {
			scheduleBooking := models.PackageScheduleBooking{
				PackageBookingID:  booking.ID,
				PackageScheduleID: schedule.ID,
				BookingType:       schedule.ScheduleType,
				BookingDate:       s.calculateBookingDate(booking.TravelDate, schedule.DayNumber),
				BookingStatus:     "pending",
			}
			if err := tx.Create(&scheduleBooking).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// BookPackageSchedules books individual flight/helicopter schedules for a package
func (s *HolidayPackageService) BookPackageSchedules(bookingID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get package booking with all details
		var packageBooking models.PackageBooking
		if err := tx.Preload("Package").
			Preload("Passengers").
			Preload("PackageScheduleBookings.PackageSchedule").
			First(&packageBooking, bookingID).Error; err != nil {
			return err
		}

		// Book each schedule through Node.js backend
		for _, scheduleBooking := range packageBooking.PackageScheduleBookings {
			nodeBookingID, err := s.bookIndividualSchedule(packageBooking, scheduleBooking)
			if err != nil {
				return fmt.Errorf("failed to book schedule %d: %v", scheduleBooking.PackageScheduleID, err)
			}

			// Update schedule booking with Node.js booking ID
			scheduleBooking.NodeBookingID = &nodeBookingID
			scheduleBooking.BookingStatus = "confirmed"
			if err := tx.Save(&scheduleBooking).Error; err != nil {
				return err
			}
		}

		// Update main booking status
		packageBooking.BookingStatus = "confirmed"
		return tx.Save(&packageBooking).Error
	})
}

// bookIndividualSchedule books a single flight/helicopter schedule through Node.js backend
func (s *HolidayPackageService) bookIndividualSchedule(packageBooking models.PackageBooking, scheduleBooking models.PackageScheduleBooking) (int, error) {
	// Prepare booking data for Node.js backend
	bookingData := map[string]interface{}{
		"bookedSeat": map[string]interface{}{
			"schedule_id": scheduleBooking.PackageSchedule.ScheduleID,
			"bookDate":    scheduleBooking.BookingDate.Format("2006-01-02"),
			"seat_labels": s.generateSeatLabels(len(packageBooking.Passengers)),
		},
		"booking": map[string]interface{}{
			"pnr":             packageBooking.PNR,
			"bookingNo":       packageBooking.BookingReference,
			"contact_no":      packageBooking.GuestPhone,
			"email_id":        packageBooking.GuestEmail,
			"noOfPassengers":  packageBooking.NumPassengers,
			"bookDate":        scheduleBooking.BookingDate.Format("2006-01-02"),
			"totalFare":       packageBooking.TotalAmount,
			"bookedUserId":    1, // Default user ID for package bookings
			"schedule_id":     scheduleBooking.PackageSchedule.ScheduleID,
		},
		"billing": map[string]interface{}{
			"user_id": 1,
		},
		"payment": map[string]interface{}{
			"user_id":        1,
			"payment_amount": packageBooking.TotalAmount,
			"payment_status": "SUCCESS",
			"transaction_id": packageBooking.BookingReference,
			"payment_mode":   "PACKAGE",
			"payment_id":     packageBooking.PaymentID,
		},
		"passengers": s.convertPassengersForNodeBackend(packageBooking.Passengers),
	}

	// Send booking request to Node.js backend
	jsonData, err := json.Marshal(bookingData)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(s.nodeBackendURL+"/bookings/complete", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("node backend returned status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	// Extract booking ID from response
	if booking, ok := result["booking"].(map[string]interface{}); ok {
		if id, ok := booking["id"].(float64); ok {
			return int(id), nil
		}
	}

	return 0, errors.New("failed to extract booking ID from node backend response")
}

// Helper functions
func (s *HolidayPackageService) calculateBookingDate(travelDate time.Time, dayNumber int) time.Time {
	return travelDate.AddDate(0, 0, dayNumber-1)
}

func (s *HolidayPackageService) generateSeatLabels(numPassengers int) []string {
	seats := make([]string, numPassengers)
	for i := 0; i < numPassengers; i++ {
		seats[i] = fmt.Sprintf("A%d", i+1) // Simple seat assignment
	}
	return seats
}

func (s *HolidayPackageService) convertPassengersForNodeBackend(passengers []models.PackagePassenger) []map[string]interface{} {
	result := make([]map[string]interface{}, len(passengers))
	for i, p := range passengers {
		result[i] = map[string]interface{}{
			"title": p.Title,
			"name":  p.FirstName + " " + p.LastName,
			"age":   p.Age,
			"type":  p.PassengerType,
		}
	}
	return result
}

// CreatePackage creates a new holiday package
func (s *HolidayPackageService) CreatePackage(pkg *models.HolidayPackage) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Create the main package
		if err := tx.Create(pkg).Error; err != nil {
			return err
		}

		// Create package schedules if provided
		for i, schedule := range pkg.PackageSchedules {
			schedule.PackageID = pkg.ID
			schedule.SequenceOrder = i + 1
			if err := tx.Create(&schedule).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
func (s *HolidayPackageService) GetBookingByID(id uint) (*models.PackageBooking, error) {
	var booking models.PackageBooking
	err := s.db.Preload("Package").
		Preload("Passengers").
		Preload("PackageScheduleBookings.PackageSchedule").
		First(&booking, id).Error
	return &booking, err
}

// GetBookingByReference retrieves a package booking by booking reference
func (s *HolidayPackageService) GetBookingByReference(reference string) (*models.PackageBooking, error) {
	var booking models.PackageBooking
	err := s.db.Where("booking_reference = ?", reference).
		Preload("Package").
		Preload("Passengers").
		Preload("PackageScheduleBookings.PackageSchedule").
		First(&booking).Error
	return &booking, err
}

// UpdateBookingStatus updates the booking status
func (s *HolidayPackageService) UpdateBookingStatus(id uint, status string) error {
	return s.db.Model(&models.PackageBooking{}).
		Where("id = ?", id).
		Update("booking_status", status).Error
}

// UpdatePaymentStatus updates the payment status
func (s *HolidayPackageService) UpdatePaymentStatus(id uint, status string, paymentID string) error {
	updates := map[string]interface{}{
		"payment_status": status,
	}
	if paymentID != "" {
		updates["payment_id"] = paymentID
	}
	return s.db.Model(&models.PackageBooking{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// CancelPackageBooking cancels a package booking and all associated schedule bookings
func (s *HolidayPackageService) CancelPackageBooking(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get booking with schedule bookings
		var booking models.PackageBooking
		if err := tx.Preload("PackageScheduleBookings").First(&booking, id).Error; err != nil {
			return err
		}

		// Cancel each individual schedule booking in Node.js backend
		for _, scheduleBooking := range booking.PackageScheduleBookings {
			if scheduleBooking.NodeBookingID != nil {
				if err := s.cancelIndividualSchedule(*scheduleBooking.NodeBookingID, scheduleBooking.BookingType); err != nil {
					// Log error but continue with other cancellations
					fmt.Printf("Failed to cancel schedule booking %d: %v\n", *scheduleBooking.NodeBookingID, err)
				}
			}
			
			// Update schedule booking status
			scheduleBooking.BookingStatus = "cancelled"
			if err := tx.Save(&scheduleBooking).Error; err != nil {
				return err
			}
		}

		// Update main booking status
		booking.BookingStatus = "cancelled"
		return tx.Save(&booking).Error
	})
}

// cancelIndividualSchedule cancels a single booking in Node.js backend
func (s *HolidayPackageService) cancelIndividualSchedule(nodeBookingID int, bookingType string) error {
	var endpoint string
	if bookingType == "helicopter" {
		endpoint = fmt.Sprintf("/bookings/helicopter/%d/cancel", nodeBookingID)
	} else {
		endpoint = fmt.Sprintf("/bookings/%d/cancel", nodeBookingID)
	}

	req, err := http.NewRequest("DELETE", s.nodeBackendURL+endpoint, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("node backend returned status %d", resp.StatusCode)
	}

	return nil
}