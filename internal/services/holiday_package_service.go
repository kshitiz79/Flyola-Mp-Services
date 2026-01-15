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

// GetAllPackageBookings retrieves all package bookings for admin
func (s *HolidayPackageService) GetAllPackageBookings() ([]models.PackageBooking, error) {
	var bookings []models.PackageBooking
	err := s.db.Preload("Package").
		Preload("Passengers").
		Order("created_at DESC").
		Find(&bookings).Error
	return bookings, err
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

// GetPackageByIDWithoutStatusFilter retrieves a specific package without status filtering (for internal use)
func (s *HolidayPackageService) GetPackageByIDWithoutStatusFilter(id uint) (*models.HolidayPackage, error) {
	var pkg models.HolidayPackage
	err := s.db.Where("id = ?", id).
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

// GetPackagesByDate retrieves packages available on a specific date
func (s *HolidayPackageService) GetPackagesByDate(dateStr string) ([]models.HolidayPackage, error) {
	// Parse the date
	targetDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	// Get the weekday for the target date (for future use)
	_ = targetDate.Weekday().String()

	// For now, we'll return packages that have schedules operating on this weekday
	// This is a simplified approach - in a real system, you'd want to check actual schedule availability
	var packages []models.HolidayPackage
	
	// Get all active packages with their schedules
	err = s.db.Where("status = ?", 1).
		Preload("PackageSchedules").
		Find(&packages).Error
	if err != nil {
		return nil, err
	}

	// Filter packages that have schedules available on the target date
	// This is a basic implementation - you might want to integrate with the Node.js backend
	// to check actual schedule availability using the getScheduleBetweenAirportDate endpoint
	var availablePackages []models.HolidayPackage
	for _, pkg := range packages {
		// For now, include all packages (you can enhance this logic)
		// In a real implementation, you'd check if the package's schedules are available on the target date
		availablePackages = append(availablePackages, pkg)
	}

	return availablePackages, nil
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
			schedule.ID = 0 // Clear any ID from request to avoid conflicts
			schedule.PackageID = pkg.ID
			schedule.SequenceOrder = i + 1
			if err := tx.Create(&schedule).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdatePackage updates an existing holiday package
func (s *HolidayPackageService) UpdatePackage(pkg *models.HolidayPackage) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get the existing package to preserve created_at
		var existingPkg models.HolidayPackage
		if err := tx.First(&existingPkg, pkg.ID).Error; err != nil {
			return err
		}

		// Preserve the original created_at
		pkg.CreatedAt = existingPkg.CreatedAt
		pkg.UpdatedAt = time.Now()

		// Update the main package (excluding associations)
		if err := tx.Omit("package_schedules").Save(pkg).Error; err != nil {
			return err
		}

		// Check if there are any bookings for this package's schedules
		var bookingCount int64
		if err := tx.Table("package_schedule_bookings").
			Joins("JOIN package_schedules ON package_schedule_bookings.package_schedule_id = package_schedules.id").
			Where("package_schedules.package_id = ?", pkg.ID).
			Count(&bookingCount).Error; err != nil {
			return err
		}

		// If there are existing bookings, don't allow schedule changes
		if bookingCount > 0 {
			// Only update the main package fields, preserve existing schedules
			// Return a custom error to indicate schedules weren't updated
			return fmt.Errorf("package updated successfully, but schedules cannot be modified as there are existing bookings")
		}

		// If no bookings exist, we can safely update schedules
		// Delete existing package schedules
		if err := tx.Where("package_id = ?", pkg.ID).Delete(&models.PackageSchedule{}).Error; err != nil {
			return err
		}

		// Create new package schedules if provided
		for i, schedule := range pkg.PackageSchedules {
			schedule.ID = 0 // Clear any ID from request to avoid conflicts
			schedule.PackageID = pkg.ID
			schedule.SequenceOrder = i + 1
			if err := tx.Create(&schedule).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// DeletePackage soft deletes a holiday package by setting status to 0
func (s *HolidayPackageService) DeletePackage(id uint) error {
	return s.db.Model(&models.HolidayPackage{}).
		Where("id = ?", id).
		Update("status", 0).Error
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
// UpdateBookingPaymentStatus updates the payment status and booking status after successful payment
func (s *HolidayPackageService) UpdateBookingPaymentStatus(bookingID uint, paymentID, paymentMethod string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get the booking
		var booking models.PackageBooking
		if err := tx.Preload("PackageScheduleBookings").First(&booking, bookingID).Error; err != nil {
			return err
		}

		// Update payment details
		booking.PaymentStatus = "paid"
		booking.BookingStatus = "confirmed"
		booking.PaymentID = paymentID
		booking.PaymentMethod = paymentMethod

		// Save the booking
		if err := tx.Save(&booking).Error; err != nil {
			return err
		}

		// Update all schedule bookings to confirmed
		for _, scheduleBooking := range booking.PackageScheduleBookings {
			scheduleBooking.BookingStatus = "confirmed"
			if err := tx.Save(&scheduleBooking).Error; err != nil {
				return err
			}
		}

		return nil
	})
}