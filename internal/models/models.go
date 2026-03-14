package models

import "time"

type Account struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID         int64     `json:"id"`
	AccountID  int64     `json:"account_id"`
	CategoryID *int64    `json:"category_id,omitempty"`
	Type       string    `json:"type"`
	Amount     float64   `json:"amount"`
	Note       string    `json:"note,omitempty"`
	Date       time.Time `json:"date"`
	CreatedAt  time.Time `json:"created_at"`

	// Joined fields for display
	AccountName  string `json:"account_name,omitempty"`
	CategoryName string `json:"category_name,omitempty"`
}

type Recurring struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	CategoryID  *int64    `json:"category_id,omitempty"`
	Frequency   string    `json:"frequency"`
	DueDay      int       `json:"due_day"`
	Type        string    `json:"type"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`

	CategoryName string `json:"category_name,omitempty"`
}

type Bucket struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Target    float64   `json:"target"`
	Current   float64   `json:"current"`
	CreatedAt time.Time `json:"created_at"`
}

type Person struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Debt struct {
	ID        int64     `json:"id"`
	PersonID  int64     `json:"person_id"`
	Amount    float64   `json:"amount"`
	Direction string    `json:"direction"`
	Note      string    `json:"note,omitempty"`
	Settled   bool      `json:"settled"`
	CreatedAt time.Time `json:"created_at"`

	PersonName string `json:"person_name,omitempty"`
}

type TransactionPerson struct {
	TransactionID int64  `json:"transaction_id"`
	PersonID      int64  `json:"person_id"`
	Note          string `json:"note,omitempty"`

	PersonName string `json:"person_name,omitempty"`
}

// PersonTransaction is a transaction linked to a person (for person history).
type PersonTransaction struct {
	ID           int64     `json:"id"`
	Type         string    `json:"type"`
	Amount       float64   `json:"amount"`
	Note         string    `json:"note,omitempty"`
	Date         time.Time `json:"date"`
	CategoryName string    `json:"category_name,omitempty"`
	AccountName  string    `json:"account_name,omitempty"`
	LinkNote     string    `json:"link_note,omitempty"`
}

// CategorySpend represents spending in a single category (for summary).
type CategorySpend struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

// PersonHistory aggregates a person's transactions and debts.
type PersonHistory struct {
	Transactions []PersonTransaction `json:"transactions"`
	Debts        []Debt              `json:"debts"`
}
