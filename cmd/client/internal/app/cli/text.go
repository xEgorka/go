package cli

import (
	"encoding/json"
	"fmt"

	"github.com/xEgorka/project3/cmd/client/internal/app/models"
)

func (a *App) ptext(j string) {
	x, err := unmarshal[models.Text](j)
	if err != nil {
		a.exitInternal(err)
	}
	if len(x.Text) > 0 {
		if _, err := fmt.Fprintf(a.tw, "Text:\n%s\n", x.Text); err != nil {
			a.exitInternal(err)
		}
	}
	if len(x.Note) > 0 {
		if _, err := fmt.Fprintf(a.tw, "Note:\t%s\t\n", x.Note); err != nil {
			a.exitInternal(err)
		}
	}
}

func (a *App) text() []byte {
	a.screen("Text: \n")
	var lines []string
	for {
		a.r.Scan()
		line := a.r.Text()
		if len(line) == 0 {
			break
		}
		lines = append(lines, line)
	}
	if err := a.r.Err(); err != nil {
		a.exitInternal(err)
	}
	var text string
	for i := 0; i < len(lines); i++ {
		text += lines[i]
		if i < len(lines)-1 {
			text += "\n"
		}
	}
	j, err := json.Marshal(models.Text{Text: text, Note: a.input("Note: ")})
	if err != nil {
		a.exitInternal(err)
	}
	return j
}
