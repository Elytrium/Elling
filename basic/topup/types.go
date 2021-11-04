package topup

import (
	"Elling/elling"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

type Method struct {
	Name                      string            `yaml:"name" json:"name,omitempty"`
	DisplayName               string            `yaml:"display-name" json:"display_name,omitempty"`
	AccountLimit              int               `yaml:"account-limit" json:"account_limit,omitempty"`
	TTL                       int64             `yaml:"ttl" json:"ttl,omitempty"`
	PayString                 string            `yaml:"pay-string" json:"pay_string,omitempty"`
	CreateRequest             elling.NetRequest `yaml:"create-request" json:"create_request"`
	CheckRequest              elling.NetRequest `yaml:"check-request" json:"check_request"`
	CheckRequestSuccessString string            `yaml:"check-request-success-string" json:"check_request_success_string,omitempty"`
	RejectRequest             elling.NetRequest `yaml:"reject-request" json:"reject_request"`
}

type PendingPurchase struct {
	TopUpID          int64
	BalanceID        int64
	Amount           int
	InvalidationDate time.Time
	Method           string
}

func (m Method) RequestTopUp(user elling.User, amount int) (PendingPurchase, error) {
	balanceID := user.Balance.ID
	date := time.Now().Add(time.Second * time.Duration(m.TTL))
	topUpID := time.Now().UnixNano()

	pendingPurchase := PendingPurchase{
		TopUpID:          topUpID,
		BalanceID:        balanceID,
		Amount:           amount,
		InvalidationDate: date,
		Method:           m.Name,
	}

	elling.DB.Create(pendingPurchase)

	_, err := m.CreateRequest.DoRequest(map[string]string{
		"{topUpID}":    strconv.FormatInt(topUpID, 10),
		"{amount}":     strconv.Itoa(amount),
		"{user_name}":  user.Email,
		"{balance_id}": strconv.FormatInt(balanceID, 10),
		"{date}":       date.String(),
	})

	return pendingPurchase, err
}

func (m Method) Validate(purchase PendingPurchase) bool {
	resp, err := m.CheckRequest.DoRequest(map[string]string{
		"{topUpID}": strconv.FormatInt(purchase.TopUpID, 10),
	})

	if err != nil {
		log.Err(err)
	}

	return resp[0] == m.CheckRequestSuccessString
}

func (m Method) GetPayString(purchase PendingPurchase) string {
	return strings.Replace(m.PayString, "{topUpID}", strconv.FormatInt(purchase.TopUpID, 10), -1)
}

func (m Method) Reject(purchase PendingPurchase) {
	_, err := m.RejectRequest.DoRequest(map[string]string{
		"{topUpID}": strconv.FormatInt(purchase.TopUpID, 10),
	})

	if err != nil {
		log.Err(err)
	}
}

func (p PendingPurchase) GetMethod() Method {
	return Instructions[p.Method].(Method)
}

func (p PendingPurchase) Validate() bool {
	return p.GetMethod().Validate(p)
}

func (p PendingPurchase) GetPayString() string {
	return p.GetMethod().GetPayString(p)
}

func (p PendingPurchase) Reject() {
	p.GetMethod().Reject(p)
}
