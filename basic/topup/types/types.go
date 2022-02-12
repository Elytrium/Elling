package types

import (
	"github.com/Elytrium/elling/basic/common"
	"github.com/Elytrium/elling/elling"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

var Instructions common.Instructions

type Method struct {
	Name                      string            `yaml:"name" json:"name,omitempty"`
	DisplayName               string            `yaml:"display-name" json:"display_name,omitempty"`
	CommissionRate            int               `yaml:"commission-rate" json:"commission-rate,omitempty"`
	TTL                       int64             `yaml:"ttl" json:"-"`
	PayString                 string            `yaml:"pay-string" json:"pay_string,omitempty"`
	CreateRequest             elling.NetRequest `yaml:"create-request" json:"-"`
	NeedToCheck               bool              `yaml:"need-to-check" json:"-"`
	CheckRequest              elling.NetRequest `yaml:"check-request" json:"-"`
	CheckRequestSuccessString string            `yaml:"check-request-success-string" json:"-,omitempty"`
	RejectRequest             elling.NetRequest `yaml:"reject-request" json:"-"`
}

type PendingPurchase struct {
	TopUpID          string
	BalanceID        uint64
	Amount           int
	InvalidationDate time.Time
	Method           string
}

func (m *Method) RequestTopUp(user *elling.User, amount int) (PendingPurchase, error) {
	balanceID := user.Balance.ID
	date := time.Now().Add(time.Second * time.Duration(m.TTL))
	topUpID := elling.NextID()

	id, err := m.CreateRequest.DoRequest(map[string]string{
		"{topUpID}":    strconv.FormatUint(topUpID, 10),
		"{amount}":     strconv.Itoa(amount),
		"{user_name}":  strconv.FormatUint(user.ID, 10),
		"{balance_id}": strconv.FormatUint(balanceID, 10),
		"{date}":       date.String(),
	})

	pendingPurchase := PendingPurchase{
		TopUpID:          id[0],
		BalanceID:        balanceID,
		Amount:           amount,
		InvalidationDate: date,
		Method:           m.Name,
	}

	if m.NeedToCheck {
		elling.DB.Create(pendingPurchase)
	}

	return pendingPurchase, err
}

func (m Method) Validate(purchase *PendingPurchase) bool {
	resp, err := m.CheckRequest.DoRequest(map[string]string{
		"{topUpID}": purchase.TopUpID,
	})

	if err != nil {
		log.Error().Err(err).Send()
	}

	return resp[0] == m.CheckRequestSuccessString
}

func (m Method) GetPayString(purchase *PendingPurchase) string {
	return strings.Replace(m.PayString, "{topUpID}", purchase.TopUpID, -1)
}

func (m Method) Reject(purchase *PendingPurchase) {
	_, err := m.RejectRequest.DoRequest(map[string]string{
		"{topUpID}": purchase.TopUpID,
	})

	if err != nil {
		log.Error().Err(err).Send()
	}
}

func (p *PendingPurchase) GetMethod() Method {
	return Instructions[p.Method].(Method)
}

func (p *PendingPurchase) Validate() bool {
	return p.GetMethod().Validate(p)
}

func (p *PendingPurchase) GetPayString() string {
	return p.GetMethod().GetPayString(p)
}

func (p *PendingPurchase) Reject() {
	p.GetMethod().Reject(p)
}
