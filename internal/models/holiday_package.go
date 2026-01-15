package models

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// HolidayPackage represents a holiday package
type HolidayPackage struct {
	ID             uint                   `json:"id" gorm:"primaryKey"`
	Title          string                 `json:"title" gorm:"not null"`
	Description    string                 `json:"description"`
	PackageType    string                 `json:"package_type" gorm:"type:enum('spiritual','wildlife','adventure','cultural');default:'spiritual'"`
	DurationDays   int                    `json:"duration_days" gorm:"default:1"`
	DurationNights int                    `json:"duration_nights" gorm:"default:0"`
	PricePerPerson float64                `json:"price_per_person" gorm:"type:decimal(10,2)"`
	MaxPassengers  int                    `json:"max_passengers" gorm:"default:6"`
	Status         int                    `json:"status" gorm:"default:1;comment:1=Active, 0=Inactive"`
	ImageURL       string                 `json:"image_url"`
	Inclusions     datatypes.JSON         `json:"inclusions"`
	Exclusions     datatypes.JSON         `json:"exclusions"`
	Itinerary      datatypes.JSON         `json:"itinerary"`
	TermsConditions string                `json:"terms_conditions"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	
	// Associations
	PackageSchedules []PackageSchedule `json:"package_schedules,omitempty" gorm:"foreignKey:PackageID"`
	Bookings         []PackageBooking  `json:"bookings,omitempty" gorm:"foreignKey:PackageID"`
}

// PackageSchedule links packages to flight/helicopter schedules
type PackageSchedule struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	PackageID     uint      `json:"package_id" gorm:"not null"`
	ScheduleType  string    `json:"schedule_type" gorm:"type:enum('flight','helicopter');not null"`
	ScheduleID    int       `json:"schedule_id" gorm:"not null;comment:References flight_schedules.id or helicopter_schedules.id from Node.js backend"`
	SequenceOrder int       `json:"sequence_order" gorm:"default:1;comment:Order of this schedule in the package"`
	DayNumber     int       `json:"day_number" gorm:"default:1;comment:Which day of the package this schedule is on"`
	IsReturn      bool      `json:"is_return" gorm:"default:false;comment:true=Return journey, false=Onward journey"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	
	// Associations
	Package HolidayPackage `json:"package,omitempty" gorm:"foreignKey:PackageID"`
}

// PackageBooking represents a booking for a holiday package
type PackageBooking struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	PackageID        uint      `json:"package_id" gorm:"not null"`
	BookingReference string    `json:"booking_reference" gorm:"uniqueIndex;size:20"`
	PNR              string    `json:"pnr" gorm:"uniqueIndex;size:10"`
	GuestName        string    `json:"guest_name" gorm:"not null"`
	GuestEmail       string    `json:"guest_email" gorm:"not null"`
	GuestPhone       string    `json:"guest_phone" gorm:"not null;size:20"`
	NumPassengers    int       `json:"num_passengers" gorm:"not null"`
	TravelDate       time.Time `json:"travel_date" gorm:"type:date;not null;comment:Start date of the package"`
	TotalAmount      float64   `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	BookingStatus    string    `json:"booking_status" gorm:"type:enum('pending','confirmed','cancelled','completed');default:'pending'"`
	PaymentStatus    string    `json:"payment_status" gorm:"type:enum('pending','paid','failed','refunded');default:'pending'"`
	PaymentID        string    `json:"payment_id"`
	PaymentMethod    string    `json:"payment_method"`
	SpecialRequests  string    `json:"special_requests"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	
	// Associations
	Package                HolidayPackage            `json:"package,omitempty" gorm:"foreignKey:PackageID"`
	Passengers             []PackagePassenger        `json:"passengers,omitempty" gorm:"foreignKey:BookingID"`
	PackageScheduleBookings []PackageScheduleBooking `json:"schedule_bookings,omitempty" gorm:"foreignKey:PackageBookingID"`
}

// PackagePassenger represents a passenger in a package booking
type PackagePassenger struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	BookingID     uint      `json:"booking_id" gorm:"not null"`
	Title         string    `json:"title" gorm:"type:enum('Mr','Mrs','Ms','Dr','Master','Miss');not null"`
	FirstName     string    `json:"first_name" gorm:"not null;size:100"`
	LastName      string    `json:"last_name" gorm:"not null;size:100"`
	Age           int       `json:"age" gorm:"not null"`
	Gender        string    `json:"gender" gorm:"type:enum('Male','Female','Other');not null"`
	PassengerType string    `json:"passenger_type" gorm:"type:enum('Adult','Child','Infant');default:'Adult'"`
	Email         string    `json:"email" gorm:"size:255;comment:Email for primary passenger (contact person)"`
	Phone         string    `json:"phone" gorm:"size:20;comment:Phone for primary passenger (contact person)"`
	IsPrimary     bool      `json:"is_primary" gorm:"default:false;comment:true=Primary passenger (contact person), false=Additional passenger"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	
	// Associations
	Booking PackageBooking `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
}

// PackageScheduleBooking links package bookings to individual schedule bookings in Node.js backend
type PackageScheduleBooking struct {
	ID                uint                `json:"id" gorm:"primaryKey"`
	PackageBookingID  uint                `json:"package_booking_id" gorm:"not null"`
	PackageScheduleID uint                `json:"package_schedule_id" gorm:"not null"`
	NodeBookingID     *int                `json:"node_booking_id" gorm:"comment:References bookings.id or helicopter_bookings.id from Node.js backend"`
	BookingType       string              `json:"booking_type" gorm:"type:enum('flight','helicopter');not null"`
	BookingDate       time.Time           `json:"booking_date" gorm:"type:date;not null"`
	SeatAssignments   datatypes.JSON      `json:"seat_assignments" gorm:"comment:Array of seat assignments for passengers"`
	BookingStatus     string              `json:"booking_status" gorm:"type:enum('pending','confirmed','cancelled');default:'pending'"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	
	// Associations
	PackageBooking  PackageBooking  `json:"package_booking,omitempty" gorm:"foreignKey:PackageBookingID"`
	PackageSchedule PackageSchedule `json:"package_schedule,omitempty" gorm:"foreignKey:PackageScheduleID"`
}

// TableName methods for custom table names
func (HolidayPackage) TableName() string {
	return "holiday_packages"
}

func (PackageSchedule) TableName() string {
	return "package_schedules"
}

func (PackageBooking) TableName() string {
	return "package_bookings"
}

func (PackagePassenger) TableName() string {
	return "package_passengers"
}

func (PackageScheduleBooking) TableName() string {
	return "package_schedule_bookings"
}

// Helper methods for JSON handling
func (p *HolidayPackage) GetInclusions() []string {
	var inclusions []string
	if p.Inclusions != nil {
		json.Unmarshal(p.Inclusions, &inclusions)
	}
	return inclusions
}

func (p *HolidayPackage) SetInclusions(inclusions []string) error {
	data, err := json.Marshal(inclusions)
	if err != nil {
		return err
	}
	p.Inclusions = data
	return nil
}

func (p *HolidayPackage) GetItinerary() []map[string]interface{} {
	var itinerary []map[string]interface{}
	if p.Itinerary != nil {
		json.Unmarshal(p.Itinerary, &itinerary)
	}
	return itinerary
}

func (p *HolidayPackage) SetItinerary(itinerary []map[string]interface{}) error {
	data, err := json.Marshal(itinerary)
	if err != nil {
		return err
	}
	p.Itinerary = data
	return nil
}

// BeforeCreate hook to generate booking reference and PNR
func (pb *PackageBooking) BeforeCreate(tx *gorm.DB) error {
	if pb.BookingReference == "" {
		pb.BookingReference = generateBookingReference()
	}
	if pb.PNR == "" {
		pb.PNR = generatePNR()
	}
	return nil
}

// Helper functions for generating unique identifiers
func generateBookingReference() string {
	// Generate unique booking reference (PKG + timestamp + random)
	return "PKG" + time.Now().Format("060102") + generateRandomString(6)
}

func generatePNR() string {
	// Generate 6-character PNR
	return generateRandomString(6)
}

func generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	
	for i := range b {
		// Use crypto/rand for cryptographically secure randomness
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return string(b)
}