package cli

import (
	"encoding/json"
	"fmt"

	"github.com/xEgorka/project3/cmd/client/internal/app/models"
)

func (a *App) ppass(j string) {
	x, err := unmarshal[models.Pass](j)
	if err != nil {
		a.exitInternal(err)
	}
	if len(x.Website) > 0 {
		if _, err := fmt.Fprintf(a.tw, "Website:\t%s\t\n", x.Website); err != nil {
			a.exitInternal(err)
		}
	}
	if len(x.Login) > 0 {
		if _, err := fmt.Fprintf(a.tw, "Login:\t%s\t\n", x.Login); err != nil {
			a.exitInternal(err)
		}
	}
	if len(x.Pass) > 0 {
		if _, err := fmt.Fprintf(a.tw, "Pass:\t%s\t\n", x.Pass); err != nil {
			a.exitInternal(err)
		}
	}
	if len(x.Note) > 0 {
		if _, err := fmt.Fprintf(a.tw, "Note:\t%s\t\n", x.Note); err != nil {
			a.exitInternal(err)
		}
	}
}

func (a *App) pass() []byte {
	j, err := json.Marshal(models.Pass{
		Website: a.input("Website: "),
		Login:   a.input("Login: "),
		Pass:    a.input("Pass: "),
		Note:    a.input("Note: "),
	})
	if err != nil {
		a.exitInternal(err)
	}
	return j
}
