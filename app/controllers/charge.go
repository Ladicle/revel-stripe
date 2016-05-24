package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
	"time"
)

type Charge struct {
	*revel.Controller
}

func (c Charge) Pay(stripeToken string, stripeEmail string, planID string) revel.Result {
	fmt.Println(stripeToken)
	fmt.Println(stripeEmail)

	stripe.Key = revel.Config.StringDefault("stripe.secretkey", "dummy")
	if stripe.Key == "dummy" {
		fmt.Println("Can not load configuration value.")
	}
	stripe.SetBackend("api", nil)

	now := time.Now()
	endOfTheMonth := time.Date(
		now.Year(),
		now.Month()+1,
		1,
		23, 59, 59, 59,
		time.Local).AddDate(0, 0, -1).Unix()

	customerParams := stripe.CustomerParams{
		Token:    stripeToken,
		Email:    stripeEmail,
		Plan:     planID,
		TrialEnd: endOfTheMonth,
	}

	isCharge := true
	customer, err := customer.New(&customerParams)
	if err != nil {
		isCharge = false
		fmt.Println(err)
	}

	return c.Render(stripeToken, stripeEmail, isCharge, customer)
}

func (c Charge) Cancel(customerID string) revel.Result {
	stripe.Key = revel.Config.StringDefault("stripe.secretkey", "dummy")
	if stripe.Key == "dummy" {
		fmt.Println("Can not load configuration value.")
	}
	stripe.SetBackend("api", nil)

	cs, err := customer.Get(customerID, nil)
	if err != nil {
		fmt.Println(err)
	}

	subp := cs.Subs.Values[0]
	if cs.Subs.ListMeta.Count < 1 {
		fmt.Println("[ERROR] Subscription is not exist")
	}
	if subp.Status != "active" {
		fmt.Println("[ERROR] Subscription status is " + subp.Status)
	}

	subParam := &stripe.SubParams{
		Customer:  customerID,
		Plan:      subp.Plan.ID,
		EndCancel: true,
	}

	_, err = sub.Cancel(subp.ID, subParam)
	if err != nil {
		fmt.Println(err)
	}

	return c.Render()
}
