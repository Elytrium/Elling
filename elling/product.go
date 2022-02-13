package elling

import (
	"errors"
	"github.com/Elytrium/elling/config"
	"github.com/rs/zerolog/log"
	"time"
)

var productMetas = make(map[string]*ProductMeta)

type Product struct {
	ID            uint64
	DisplayName   string
	Billing       string
	BalanceID     uint64 `json:"-"`
	Module        string
	Type          string
	Users         []*User `gorm:"many2many:user_products;"`
	PaidTill      time.Time
	Suspended     bool
	SuspendReason string
}

type ProductMeta struct {
	DisplayName       string
	Name              string
	Type              string
	AvailableBillings map[string]bool
	BillingCost       map[string]int64
}

type ProductOrderedEvent struct {
	Product *Product
}

type ProductSuspendEvent struct {
	Product *Product
	Reason  string
}

type ProductUnsuspendEvent struct {
	Product *Product
}

type ProductDeletedEvent struct {
	Product *Product
}

type ProductListener struct{}

var BillingTypeNotSupportedError = errors.New("specified billing type is not supported by the product")

func init() {
	if config.AppConfig.IsMaster {
		registerLocalListener(ProductListener{})
	}
}

func (u *User) AddProduct(displayName, billing, module, typeName string, user User) error {
	product := &Product{
		ID:            NextID(),
		DisplayName:   displayName,
		Billing:       billing,
		BalanceID:     user.Balance.ID,
		Module:        module,
		Type:          typeName,
		Users:         []*User{u},
		PaidTill:      time.Time{},
		Suspended:     false,
		SuspendReason: "",
	}

	DispatchEvent(ProductOrderedEvent{Product: product})
	return product.Update()
}

func (p *Product) SuspendProduct(reason string) error {
	p.Suspended = true
	p.SuspendReason = reason
	DispatchEvent(ProductSuspendEvent{Product: p, Reason: reason})
	return p.Update()
}

func (p *Product) UnsuspendProduct() error {
	p.Suspended = false
	p.SuspendReason = ""
	DispatchEvent(ProductUnsuspendEvent{Product: p})
	return p.Update()
}

func (p *Product) GetMeta() *ProductMeta {
	return productMetas[p.Type]
}

func (p *Product) ChangeBilling(billing *Billing) error {
	if !p.GetMeta().AvailableBillings[billing.Name] {
		return BillingTypeNotSupportedError
	}

	p.Billing = billing.Name
	return p.Update()
}

func (p *Product) GetCost() int64 {
	return p.GetMeta().BillingCost[p.Billing]
}

func (p *Product) Prolong(period int64) error {
	bal, err := p.GetBalance()
	if err != nil {
		return err
	}

	if bal.Amount < p.GetCost() {
		return NotEnoughFundsError
	}

	err = bal.Withdraw(period * p.GetCost())
	if err != nil {
		return err
	}

	return p.localProlong(period)
}

func (p *Product) GetBalance() (*Balance, error) {
	var balance Balance
	DB.Find(&balance, p.BalanceID)

	return &balance, DB.Error
}

func (p *Product) localProlong(period int64) error {
	billing := GetBilling(p.Billing)
	p.PaidTill.Add(time.Duration(period) * billing.Duration)
	return p.Update()
}

func (p *Product) Update() error {
	DB.Save(p)
	return DB.Error
}

func RegisterProduct(meta *ProductMeta) {
	productMetas[meta.Type] = meta
}

func (*ProductListener) OnSmallTick(_ SmallTickEvent) {
	var products []Product
	DB.Where("PaidTill > ?", time.Now()).Find(&products)

	localBalIdSet := make(map[uint64]bool)
	for _, product := range products {
		if _, ok := localBalIdSet[product.BalanceID]; !ok {
			localBalIdSet[product.BalanceID] = true
		}
	}

	localBalIdSlice := make([]uint64, len(localBalIdSet))

	i := 0
	for k := range localBalIdSet {
		localBalIdSlice[i] = k
		i++
	}

	var localBalCacheSlice []Balance
	DB.Find(&localBalCacheSlice, localBalIdSlice)

	localBalCache := make(map[uint64]*Balance)
	for _, bal := range localBalCacheSlice {
		localBalCache[bal.ID] = &bal
	}

	for _, product := range products {
		if localBalCache[product.BalanceID].Amount < product.GetCost() {
			err := product.SuspendProduct("unpaid")

			if err != nil {
				log.Error().Err(err).Msg("Suspending unpaid products")
			}
		} else {
			err := localBalCache[product.BalanceID].Withdraw(product.GetCost())

			if err != nil {
				log.Error().Err(err).Msg("Proceeding payment")
			}

			err = product.localProlong(1)

			if err != nil {
				log.Error().Err(err).Msg("Prolong of product")
			}
		}
	}
}
