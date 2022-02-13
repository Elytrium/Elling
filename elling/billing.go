package elling

import (
	"time"
)

type Billing struct {
	Name     string
	Duration time.Duration
}

var registeredBillings = make(map[string]*Billing)

var OneTimeBilling = RegisterBilling("onetime", 0)

var HourlyBilling = RegisterBilling("hourly", time.Hour)

var DailyBilling = RegisterBilling("daily", time.Hour*24)

var WeeklyBilling = RegisterBilling("weekly", time.Hour*24*7)

var MonthlyBilling = RegisterBilling("monthly", time.Hour*24*30)

var QuarterlyBilling = RegisterBilling("quarterly", time.Hour*24*90)

var YearlyBilling = RegisterBilling("yearly", time.Hour*24*365)

func RegisterBilling(name string, duration time.Duration) *Billing {
	billing := &Billing{
		Name:     name,
		Duration: duration,
	}

	registeredBillings[name] = billing

	return billing
}

func GetBilling(name string) *Billing {
	return registeredBillings[name]
}
