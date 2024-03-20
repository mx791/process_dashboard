package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mx791/process_dashboard"
)

func PriceCall(days int, strike float64, currentPrice float64, impliedVolatility float64) float64 {
	price := 0.0
	iters := 50_000.0
	for i := 0.0; i < iters; i += 1.0 {
		p := currentPrice
		for e := 0; e < days; e++ {
			p *= 1.0 + rand.NormFloat64()*impliedVolatility/15.87
		}
		if p > strike {
			price += p - strike
		}
	}
	return price / iters
}

func main() {
	IV := 0.1
	best_loss := -1.0
	target := 2.99
	lr := 0.5

	task := func() map[string]string {
		c_IV := IV*(1.0-lr) + rand.Float64()*lr
		lr *= 0.99
		p := PriceCall(5, 175.0, 176.08, c_IV)
		loss := (target - p) * (target - p)
		if (loss < best_loss) || (best_loss == -1.0) {
			IV = c_IV
			best_loss = loss
		}
		time.Sleep(2 * time.Second)
		out := make(map[string]string)
		out["last_best_loss"] = fmt.Sprintf("%f", best_loss)
		out["last_loss"] = fmt.Sprintf("%f", loss)
		out["current_iv"] = fmt.Sprintf("%f", IV)
		out["price"] = fmt.Sprintf("%f", p)
		return out
	}
	d := process_dashboard.DashBoard{task, 100, []process_dashboard.Callback{
		process_dashboard.LastValueCallback{"Implied volatility", "current_iv"},
		process_dashboard.LastValueCallback{"Loss", "last_best_loss"},
		process_dashboard.LastValueCallback{"Current iteration", "currentIteration"},
		process_dashboard.LineCallback{"Best loss", "last_best_loss"},
		process_dashboard.LineCallback{"Tested loss", "last_loss"},
		process_dashboard.LineCallback{"Implied volatility", "current_iv"},
		process_dashboard.LineCallback{"Current Price", "price"},
	}}
	d.Run()
}
