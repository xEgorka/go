package cli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/xEgorka/project3/cmd/client/internal/app/models"
)

func (a *App) pfile(j string) {
	x, err := unmarshal[models.File](j)
	if err != nil {
		a.exitInternal(err)
	}
	if err := a.dec(x); err != nil {
		a.exit("File save error", failure)
	} else {
		_, err := fmt.Fprintf(a.w, "File saved: %s\n", x.Name)
		if err != nil {
			a.exitInternal(err)
		}
		if err := a.w.Flush(); err != nil {
			a.exitInternal(err)
		}
	}
	if len(x.Note) > 0 {
		if _, err := fmt.Fprintf(a.tw, "Note:\t%s\t\n", x.Note); err != nil {
			a.exitInternal(err)
		}
	}
}

func (a *App) dec(x *models.File) error {
	dec, err := base64.StdEncoding.DecodeString(x.Data)
	if err != nil {
		return err
	}
	f, err := os.Create(x.Name)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			a.exitInternal(err)
		}
	}()
	if _, err := f.Write(dec); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	return nil
}

func (a *App) file() []byte {
	fpath := a.input("File path: ")
	f, e := os.Open(fpath)
	if e != nil {
		a.exit("File open error", failure)
	}
	defer func() {
		if er := f.Close(); er != nil {
			a.exitInternal(er)
		}
	}()

	if fi, err := os.Stat(fpath); err != nil {
		a.exitInternal(err)
	} else if fi.Size() > int64(a.cfg.FileSizeLimit) {
		_, err := fmt.Fprintf(a.w, "File size exceeds limit %d MB\n",
			a.cfg.FileSizeLimit/1024/1024)
		if err != nil {
			a.exitInternal(err)
		}
		if err := a.w.Flush(); err != nil {
			a.exitInternal(err)
		}
		a.exit("File size error", failure)
	}

	j, e := json.Marshal(models.File{
		Name: filepath.Base(f.Name()),
		Data: a.enc(f),
		Note: a.input("Note: "),
	})
	if e != nil {
		a.exitInternal(e)
	}
	return j
}

const bufSize = 1024

func (a *App) enc(f *os.File) string {
	var b []byte
	for {
		buf := make([]byte, bufSize)
		if _, err := f.Read(buf); err == io.EOF {
			break
		} else if err != nil {
			a.exitInternal(err)
		}
		b = append(b, buf...)
	}
	return base64.StdEncoding.EncodeToString(b)
}
