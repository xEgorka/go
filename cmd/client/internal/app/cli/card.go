package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/xEgorka/project3/cmd/client/internal/app/models"
)

func (a *App) pcard(j string) {
	x, err := unmarshal[models.Card](j)
	if err != nil {
		a.exitInternal(err)
	}
	if x.Num > 0 {
		if _, err := fmt.Fprintf(a.tw, "Num:\t%d\t\n", x.Num); err != nil {
			a.exitInternal(err)
		}
	}
	if x.Exp > 0 {
		if _, err := fmt.Fprintf(a.tw, "Exp:\t%d\t\n", x.Exp); err != nil {
			a.exitInternal(err)
		}
	}
	if x.CVV > 0 {
		if _, err := fmt.Fprintf(a.tw, "CVV:\t%d\t\n", x.CVV); err != nil {
			a.exitInternal(err)
		}
	}
	if len(x.Bank) > 0 {
		if _, err := fmt.Fprintf(a.tw, "Bank:\t%s\t\n", x.Bank); err != nil {
			a.exitInternal(err)
		}
	}
	if len(x.Note) > 0 {
		if _, err := fmt.Fprintf(a.tw, "Note:\t%s\t\n", x.Note); err != nil {
			a.exitInternal(err)
		}
	}
}

func (a *App) card() []byte {
	j, err := json.Marshal(models.Card{
		Num:  a.num(),
		CVV:  a.cvv(),
		Exp:  a.exp(),
		Bank: a.input("Bank: "),
		Note: a.input("Note: "),
	})
	if err != nil {
		a.exitInternal(err)
	}
	return j
}

func (a *App) num() int {
	num, err := strconv.Atoi(a.input("Card number: "))
	if err != nil || !valid(num) {
		a.exit("Invalid card number", failure)
	}
	return num
}

func (a *App) cvv() int {
	cvv, err := strconv.Atoi(a.input("CVV: "))
	if err != nil {
		a.exit("Invalid CVV", failure)
	}
	return cvv
}

func (a *App) exp() int {
	exp, err := strconv.Atoi(a.input("Expired date MMYY: "))
	if err != nil {
		a.exit("Invalid expired date", failure)
	}
	return exp
}
