package controllers

import "github.com/revel/revel"
import "github.com/stripe/stripe-go"
import "github.com/stripe/stripe-go/plan"
import "fmt"

type App struct {
	*revel.Controller
}

var publishkey string = ""

func init() {
	revel.OnAppStart(func() {
		var success bool = false
		if publishkey, success = revel.Config.String("stripe.publishkey"); !success {
			revel.ERROR.Fatal("Not found stripe.publishkey.")
		}
		if stripe.Key, success = revel.Config.String("stripe.secretkey"); !success {
			revel.ERROR.Fatal("Not found stripe.secretkey.")
		}
	})
}

type Plan struct {
	ID       string
	Amount   uint64
	Name     string
	Key      string
	Currency string
}

func (c App) Index() revel.Result {

	stripe.SetBackend("api", nil)
	list := plan.List(&stripe.PlanListParams{})
	plans := []Plan{}
	for list.Next() {
		fmt.Println(list.Plan())
		plans = append(plans, Plan{
			ID:       list.Plan().ID,
			Amount:   list.Plan().Amount,
			Name:     list.Plan().Name,
			Currency: string(list.Plan().Currency),
			Key:      publishkey,
		})
	}

	return c.Render(plans, publishkey)
}
