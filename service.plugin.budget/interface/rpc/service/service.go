package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/libraries/monzo"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/interface/rpc/proto"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/registry"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/usecase"
)

var (
	errNoMonthlyBudget = errors.New(fmt.Errorf("no monthly budget"))
)

// Service is a implementation of the RPC plugin service.
type Service struct {
	pu usecase.PreferencesUsecase
	hc *http.Client
}

// NewService returns a new instance of Service.
func NewService(ctn *di.Container) *Service {
	pu := ctn.Resolve(registry.PreferencesUsecaseService).(usecase.PreferencesUsecase)

	return &Service{
		pu: pu,
		hc: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// Send is a handle for the gRPC plugin service. It is the entrypoint
// for the plugin handler and calculates the daily budget.
func (s *Service) Send(ctx context.Context, in *proto.SendRequest) (*proto.EmptySendResponse, error) {
	userID, accountID, accessToken := in.GetUserID(), in.GetAccountID(), in.GetAccessToken()

	budget, err := s.calculateBudget(userID, accountID, accessToken)
	if err != nil {
		if err == errNoMonthlyBudget {
			return &proto.EmptySendResponse{}, nil
		}

		return &proto.EmptySendResponse{}, fmt.Errorf(err.Text())
	}

	err = s.sendFeedItem(accountID, accessToken, budget)
	if err != nil {
		return &proto.EmptySendResponse{}, fmt.Errorf(err.Text())
	}

	return &proto.EmptySendResponse{}, nil
}

func (s *Service) calculateBudget(userID, accountID, accessToken string) (float64, errors.Error) {
	total, today, monthly := make(chan float64, 1), make(chan float64, 1), make(chan float64, 1)

	var eg errors.Group
	eg.Go(func() errors.Error {
		sum, err := s.sumOfTransactions(accountID, accessToken)
		if err != nil {
			return err
		}

		total <- float64(sum)

		return nil
	})
	eg.Go(func() errors.Error {
		spent, err := s.getSpentToday(accountID, accessToken)
		if err != nil {
			return err
		}

		today <- float64(spent)

		return nil
	})
	eg.Go(func() errors.Error {
		mb, err := s.pu.GetMonthlyBudget(userID)
		if err != nil {
			return err
		}

		monthly <- float64(mb)

		return nil
	})

	if err := eg.Wait(); err != nil {
		return 0, err
	}

	totalSpent, spentToday, monthlyBudget := <-total, <-today, <-monthly

	if monthlyBudget == 0 {
		return 0, errNoMonthlyBudget
	}

	r := monthlyBudget / totalSpent
	r = r / float64(daysLeftInMonth())
	r = r - spentToday
	r = r / 100.0

	return r, nil
}

func (s *Service) sumOfTransactions(accountID, accessToken string) (int, errors.Error) {
	now := time.Now().UTC()
	if now.Day() == 1 {
		return 0, nil // nothing spent before today, this month
	}

	monthStart := time.Date(now.Year(), time.Month(now.Month()), 1, 0, 0, 0, 0, time.UTC)
	todayStart := time.Date(now.Year(), time.Month(now.Month()), now.Day(), 0, 0, 0, 0, time.UTC)

	trans, err := monzo.GetTransactions(accountID, accessToken, &monzo.TransactionOpts{
		Since:  &monthStart,
		Before: &todayStart,
	})
	if err != nil {
		return 0, err
	}

	var sum int

	// Add up the amounts from each transaction, where they're not
	// excluded from the user's summary or if the amount is not pending.
	// Only transactions which result in money, subracted from the user's
	// account is counted.
	for _, t := range trans {
		if !t.IncludeInSpending || t.AmountIsPending || t.Amount > 0 {
			continue
		}

		sum += t.Amount
	}

	return sum * -1, nil
}

func (s *Service) getSpentToday(accountID, accessToken string) (int, errors.Error) {
	balance, err := monzo.GetBalance(accountID, accessToken)
	if err != nil {
		return 0, err
	}

	return balance.SpendToday * -1, nil
}

func (s *Service) sendFeedItem(accountID, accessToken string, budget float64) errors.Error {
	var message, more string
	textColor := "#333333"
	switch true {
	case budget < 0:
		budget = budget * -1
		message = fmt.Sprintf("You've spent £%.2f over your daily budget!", budget)
		more = "It looks like you've spent too much today!"
		textColor = "#cc4e4e"
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

	imgURL := "https://media.giphy.com/media/ND6xkVPaj8tHO/giphy.gif"
	err := monzo.CreateFeedItem(accountID, accessToken, message, imgURL, &monzo.FeedItemOpts{
		BodyColor:  textColor,
		TitleColor: textColor,
		Body:       more,
	})
	if err != nil {
		return err
	}

	return nil
}

func daysLeftInMonth() int {
	now := time.Now().UTC()

	nextMonth := now.AddDate(0, 1, 0)
	nextMonth = time.Date(nextMonth.Year(), time.Month(nextMonth.Month()), 1, 0, 0, 0, 0, time.UTC)
	today := time.Date(now.Year(), time.Month(now.Month()), now.Day(), 0, 0, 0, 0, time.UTC)

	return int(nextMonth.Sub(today).Hours() / float64(24))
}
