package validation

import (
	"encoding/json"
	"log"
	"errors"
)


type OrderMessage struct {
	OrderUID          string          `json:"order_uid"`
	TrackNumber       string          `json:"track_number"`
	Entry             string          `json:"entry"`
	Delivery          DeliveryInfo    `json:"delivery"`
	Payment           PaymentInfo     `json:"payment"`
	Items             []OrderItem     `json:"items"`
	Locale            string          `json:"locale"`
	InternalSignature string          `json:"internal_signature"`
	CustomerID        string          `json:"customer_id"`
	DeliveryService   string          `json:"delivery_service"`
	ShardKey          string          `json:"shardkey"`
	SMID              int             `json:"sm_id"`
	DateCreated       string          `json:"date_created"`
	OOFShard          string          `json:"oof_shard"`
}

type DeliveryInfo struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type PaymentInfo struct {
	Transaction     string `json:"transaction"`
	RequestID        string `json:"request_id"`
	Currency         string `json:"currency"`
	Provider         string `json:"provider"`
	Amount           int    `json:"amount"`
	PaymentDT        int64  `json:"payment_dt"`
	Bank             string `json:"bank"`
	DeliveryCost     int    `json:"delivery_cost"`
	GoodsTotal       int    `json:"goods_total"`
	CustomFee        int    `json:"custom_fee"`
}

type OrderItem struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NMID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}


func Validate(value interface{}) error {
	jsonString, ok := value.(string)
	if !ok {
		return errors.New("value is not a valid JSON string")
	}
	var orderMessage OrderMessage
		if err := json.Unmarshal([]byte(jsonString), &orderMessage); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			return err
		}

		if err := validateOrderMessage(orderMessage); err != nil {
			log.Printf("Invalid order message: %v", err)
			return err
		}
		return nil
}

func validateOrderMessage(orderMessage OrderMessage) error {
	if orderMessage.OrderUID == "" {
		return errors.New("missing required field: OrderUID")
	}
	if orderMessage.Delivery.Name == "" {
		return errors.New("missing required field: Delivery")
	}
	if orderMessage.Payment.Transaction == "" {
		return errors.New("missing required field: Transaction")
	}
	if len(orderMessage.Items) < 1 {
		return errors.New("missing required field: Items")
	}
	if orderMessage.TrackNumber == "" {
		return errors.New("missing required field: TrackNumber")
	}
	return nil
}