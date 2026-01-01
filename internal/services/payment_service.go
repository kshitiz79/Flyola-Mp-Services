package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flyola-services/internal/models"
	"io"
	"net/http"

	"gorm.io/gorm"
)

type PaymentService struct {
	db *gorm.DB
}

func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{db: db}
}

func (s *PaymentService) CreateRazorpayOrder(amount int64, currency string, receipt string, keyID, keySecret string) (string, error) {
	orderData := map[string]interface{}{
		"amount":   amount,
		"currency": currency,
		"receipt":  receipt,
	}

	jsonData, err := json.Marshal(orderData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.razorpay.com/v1/orders", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	auth := keyID + ":" + keySecret
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", errors.New("failed to create razorpay order: " + string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	orderID, ok := result["id"].(string)
	if !ok {
		return "", errors.New("id not found in razorpay response")
	}

	return orderID, nil
}

func (s *PaymentService) VerifyRazorpaySignature(orderID, paymentID, signature, secret string) error {
	payload := orderID + "|" + paymentID
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	generatedSignature := hex.EncodeToString(h.Sum(nil))

	if generatedSignature != signature {
		return errors.New("invalid payment signature")
	}
	return nil
}

func (s *PaymentService) ProcessPayment(payment *models.HotelPayment) error {
	return s.db.Create(payment).Error
}

func (s *PaymentService) GetPaymentByID(id uint) (*models.HotelPayment, error) {
	var payment models.HotelPayment
	err := s.db.Preload("Booking").First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (s *PaymentService) GetPaymentByBooking(bookingID uint) (*models.HotelPayment, error) {
	var payment models.HotelPayment
	err := s.db.Preload("Booking").Where("booking_id = ?", bookingID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (s *PaymentService) UpdatePaymentStatus(id uint, status string) (*models.HotelPayment, error) {
	var payment models.HotelPayment
	if err := s.db.First(&payment, id).Error; err != nil {
		return nil, err
	}

	payment.Status = status
	if err := s.db.Save(&payment).Error; err != nil {
		return nil, err
	}

	return &payment, nil
}

func (s *PaymentService) GetPaymentsByStatus(status string) ([]models.HotelPayment, error) {
	var payments []models.HotelPayment
	err := s.db.Preload("Booking").Where("status = ?", status).Find(&payments).Error
	return payments, err
}

func (s *PaymentService) RefundPayment(id uint, refundAmount float64) (*models.HotelPayment, error) {
	var payment models.HotelPayment
	if err := s.db.First(&payment, id).Error; err != nil {
		return nil, err
	}

	payment.Status = "refunded"
	payment.RefundAmount = refundAmount
	if err := s.db.Save(&payment).Error; err != nil {
		return nil, err
	}

	return &payment, nil
}
