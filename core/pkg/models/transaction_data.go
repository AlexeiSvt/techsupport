package models

import "time"

// Transaction represents a financial or significant event within the system.
// In a graph database, it is a discrete node used to calculate risk velocity 
// and detect suspicious payment patterns across multiple accounts.
type Transaction struct {
    // TransactionID is the unique reference identifier from the payment provider or internal gateway.
    TransactionID string `json:"transaction_id" cypher:"id"`

    // Amount is the numerical value of the financial operation.
    Amount float64 `json:"amount" cypher:"amount"`

    // Currency is the ISO 4217 code (e.g., "USD", "EUR") used for the transaction.
    Currency string `json:"currency" cypher:"currency"`

    // PaymentMethod is the category of payment (e.g., "credit_card", "e-wallet").
    PaymentMethod string `json:"payment_method" cypher:"method"`

    // PaymentProvider is the third-party service processing the payment (e.g., "Stripe").
    PaymentProvider string `json:"payment_provider" cypher:"provider"`

    // DeviceID is the hardware fingerprint. Used to link this node to a specific Device node.
    DeviceID string `json:"device_id" cypher:"device_id"`

    // DeviceModel is the marketing name of the hardware used (e.g., "iPhone 15 Pro").
    DeviceModel string `json:"device_model" cypher:"device_model"`

    // IP is the network address from which the transaction was initiated.
    IP string `json:"ip" cypher:"ip"`

    // Country is the geographic origin determined from the transaction IP.
    Country string `json:"country" cypher:"country"`

    // City is the specific city determined from the transaction IP.
    City string `json:"city" cypher:"city"`

    // Timestamp is the RFC3339 time when the transaction was processed.
    Timestamp time.Time `json:"timestamp" cypher:"timestamp"`

    // ASN (Autonomous System Number) identifies the ISP network used for the transaction.
    ASN string `json:"asn" cypher:"asn"`
}