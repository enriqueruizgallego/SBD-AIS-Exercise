package model

import (
	"fmt"
	"time"
)

const (
	orderFilename = "order_%d.md"

	// todo create markdown emplate, fields should be able to be populated with fmt.Sprintf
	markdownTemplate = `# Order: %d

| Created At      | Drink ID | Amount |
|-----------------|----------|--------|
| %s | %d        | %d      |

<<<<<<< HEAD
Thanks for drinking with us!
=======
Thanks for drinking with us! 
>>>>>>> a04fbb2 (Adding solution assignment 6)
`
)

type Order struct {
	Base
	Amount uint64 `json:"amount"`
	// Relationships
	// foreign key
	DrinkID uint  `json:"drink_id" gorm:"not null"`
	Drink   Drink `json:"drink"`
}

// this function generates the receipt content in Markdown format.
func (o *Order) ToMarkdown() string {
	return fmt.Sprintf(markdownTemplate, o.ID, o.CreatedAt.Format(time.Stamp), o.DrinkID, o.Amount)
}

// returns the filename for each order
func (o *Order) GetFilename() string {
	return fmt.Sprintf(orderFilename, o.ID)
}
