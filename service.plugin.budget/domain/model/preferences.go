package model

import "fmt"

type Preferences struct {
	userID string

	// represents pennies.
	monthlyBudget int
}

func NewPreferences(userID string) *Preferences {
	return &Preferences{
		userID: userID,
	}
}

// GetUserID returns the user's id.
func (p *Preferences) GetUserID() string {
	return p.userID
}

// GetMonthlyBudget returns the monthly budget.
func (p *Preferences) GetMonthlyBudget() int {
	return p.monthlyBudget
}

// UpdateMonthlyBudget updates the monthly budget.
func (p *Preferences) UpdateMonthlyBudget(budget int) error {
	if budget < 100 {
		return fmt.Errorf("monthly budget cannot be less than 100")
	}

	p.monthlyBudget = budget

	return nil
}

func PreferencesFromScan() (*Preferences, []interface{}) {
	p := new(Preferences)

	return p, []interface{}{
		&p.userID,
		&p.monthlyBudget,
	}
}
