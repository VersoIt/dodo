package domain

type OrderItemPayload struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
}

type OrderPaidEvent struct {
	OrderID     string             `json:"order_id"`
	OrderNumber string             `json:"order_number"`
	Items       []OrderItemPayload `json:"items"`
}

type OrderReadyEvent struct {
	OrderID     string             `json:"order_id"`
	OrderNumber string             `json:"order_number"`
	Address     DeliveryAddress    `json:"address"`
	Items       []OrderItemPayload `json:"items"`
}

type OrderStatusChangedEvent struct {
	CustomerID string `json:"customer_id"`
	OrderID    string `json:"order_id"`
	Status     string `json:"status"`
}
