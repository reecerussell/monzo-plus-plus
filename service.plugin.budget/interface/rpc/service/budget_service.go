package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/interface/rpc/proto"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/registry"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/usecase"
)

type BudgetService struct {
	pu usecase.PreferencesUsecase
	hc *http.Client
}

func NewBudgetService(ctn *di.Container) *BudgetService {
	pu := ctn.Resolve(registry.PreferencesUsecaseService).(usecase.PreferencesUsecase)

	return &BudgetService{
		pu: pu,
		hc: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (bs *BudgetService) Calculate(ctx context.Context, data *proto.CalculateData) (*proto.CalculateResponse, error) {
	accountID := data.GetAccountID()
	accessToken := data.GetAccessToken()

	log.Printf("Received calculating request:\nAccountID: %s\nAccess Token: %s\n", accountID, accessToken)

	totalSpent, err := bs.sumOfTransactions(accountID, accessToken)
	if err != nil {
		return nil, fmt.Errorf("sof: %v", err)
	}
	log.Printf("Sum of transactions: %d\n", totalSpent)

	spentToday, err := bs.getSpentToday(accountID, accessToken)
	if err != nil {
		return nil, fmt.Errorf("st: %v", err)
	}
	log.Printf("Spent today: %d\n", spentToday)

	monthlyBudget, mbErr := bs.pu.GetMonthlyBudget(data.GetUserID())
	if mbErr != nil {
		return nil, fmt.Errorf("mb: %v", mbErr.Text())
	}
	log.Printf("Monthly budget: %d\n", monthlyBudget)

	if err != nil {
		log.Printf("Failed to calulate: %v\n", err)
		return nil, err
	}

	if monthlyBudget == 0 {
		log.Printf("Failed to calculate: no monthly budget.\n")
		return &proto.CalculateResponse{}, nil
	}

	result := float64(monthlyBudget) - float64(totalSpent)
	result = result / float64(daysLeftInMonth())
	result = result - float64(spentToday)

	err = bs.sendFeedItem(accountID, accessToken, result/float64(100))
	if err != nil {
		return nil, fmt.Errorf("feed: %v", err)
	}

	log.Printf("Caluclated: %.2f\n", result/float64(100))

	return &proto.CalculateResponse{}, nil
}

func (bs *BudgetService) sumOfTransactions(accountID, accessToken string) (int, error) {
	now := time.Now().UTC()
	if now.Day() == 1 {
		return 0, nil // nothing spent before today, this month
	}

	monthStart := time.Date(now.Year(), time.Month(now.Month()), 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	todayStart := time.Date(now.Year(), time.Month(now.Month()), now.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	url := fmt.Sprintf("https://api.monzo.com/transactions?account_id=%s&since=%s&before=%s", accountID, monthStart, todayStart)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("req: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	log.Printf("Attempting to fetch transactions...\n")

	resp, err := bs.hc.Do(req)
	if err != nil {
		log.Printf("\tFailed to fetch transactions: %v\n", err)
		return 0, fmt.Errorf("res: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp errorResponse
		_ = json.NewDecoder(resp.Body).Decode(&errorResp)
		log.Printf("\tFailed to fetch transactions: %s\n", errorResp.Message)
		return 0, fmt.Errorf("monzo: %s", errorResp.Message)
	}

	var data transactions
	_ = json.NewDecoder(resp.Body).Decode(&data)

	var sum int

	// Add up the amounts from each transaction, where they're not
	// excluded from the user's summary or if the amount is not pending.
	// Only transactions which result in money, subracted from the user's
	// account is counted.
	for _, t := range data.Transactions {
		if !t.Include || t.IsPending || t.Amount > 0 {
			continue
		}

		sum += t.Amount
	}

	log.Printf("\tSuccessfully fetchs transactions: sum: %.2f\n", float64(sum*-1)/float64(100))

	return sum * -1, nil
}

func (bs *BudgetService) getSpentToday(accountID, accessToken string) (int, error) {
	url := fmt.Sprintf("https://api.monzo.com/balance?account_id=%s", accountID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("req: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	log.Printf("Attempted to fetch balance...\n")
	resp, err := bs.hc.Do(req)
	if err != nil {
		return 0, fmt.Errorf("res: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp errorResponse
		_ = json.NewDecoder(resp.Body).Decode(&errorResp)

		log.Printf("\tFailed fetch balance: %s\n", errorResp.Message)

		return 0, fmt.Errorf("monzo: %s", errorResp.Message)
	}

	var data balance
	_ = json.NewDecoder(resp.Body).Decode(&data)

	log.Printf("\tFetched balance: %.2f\n", float64(data.SpentToday)/float64(100))

	return data.SpentToday * -1, nil
}

func (bs *BudgetService) sendFeedItem(accountID, accessToken string, budget float64) error {
	var message, more string
	textColour := "#333333"
	switch true {
	case budget < 0:
		budget = budget * -1
		message = fmt.Sprintf("You've spent £%.2f over your daily budget!", budget)
		more = "It looks like you've spent too much today!"
		textColour = "#cc4e4e"
		break
	case budget == 0:
		message = fmt.Sprintf("You've spent exactly your daily budget today!")
		more = "You're cutting this close!"
		break
	case budget > 0:
		message = fmt.Sprintf("You've got £%.2f left of your daily budget.", budget)
		more = "You're within budget; looking good!"
		break
	}

	body := url.Values{}
	body.Set("type", "basic")
	body.Set("account_id", accountID)
	body.Set("params[title]", message)
	body.Set("params[image_url]", "https://media.giphy.com/media/ND6xkVPaj8tHO/giphy.gif")
	body.Set("params[body]", more)
	body.Set("params[title_color]", textColour)
	body.Set("params[body_color]", textColour)

	req, err := http.NewRequest(http.MethodPost, "https://api.monzo.com/feed", strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("req: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	log.Printf("Attempting to post feed item...\n")
	resp, err := bs.hc.Do(req)
	if err != nil {
		log.Printf("\tFailed to post feed item: %v\n", err)
		return fmt.Errorf("post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("\tSucecssfully posted feed item.\n")
		return nil
	}

	var errorResp errorResponse
	_ = json.NewDecoder(resp.Body).Decode(&errorResp)

	log.Printf("\tFailed to post feed item: %s\n", errorResp.Message)

	return fmt.Errorf("monzo: %v", errorResp.Message)
}

func daysLeftInMonth() int {
	now := time.Now().UTC()

	nextMonth := now.AddDate(0, 1, 0)
	nextMonth = time.Date(nextMonth.Year(), time.Month(nextMonth.Month()), 1, 0, 0, 0, 0, time.UTC)
	today := time.Date(now.Year(), time.Month(now.Month()), now.Day(), 0, 0, 0, 0, time.UTC)

	return int(nextMonth.Sub(today).Hours() / float64(24))
}

// transaction is used to read transaction data from monzo. This only
// contains properties that are required for the calculation.
type transaction struct {
	Amount    int  `json:"amount"`
	IsPending bool `json:"amount_is_pending"`
	Include   bool `json:"include_in_spending"`
}

// transactions is a wrapper object, arround an array of transactions.
type transactions struct {
	Transactions []transaction `json:"transactions"`
}

// balance is used to read the user's balance.
type balance struct {
	SpentToday int `json:"spend_today"`
}

// error response is used to read a monzo error response message.
type errorResponse struct {
	Message string `json:"message"`
}
