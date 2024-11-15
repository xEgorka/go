package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xEgorka/project3/cmd/client/internal/app/config"
	"github.com/xEgorka/project3/cmd/client/internal/app/crypto"
	"github.com/xEgorka/project3/cmd/client/internal/app/logger"
	"github.com/xEgorka/project3/cmd/client/internal/app/service"
)

// App provides application.
type App struct {
	cfg *config.Config
	s   *service.Service
	c   crypto.Crypto
	r   *bufio.Scanner
	w   *bufio.Writer
	tw  *tabwriter.Writer
}

// New creates App.
func New(config *config.Config, service *service.Service, crypto crypto.Crypto) *App {
	return &App{cfg: config, s: service, c: crypto,
		r:  bufio.NewScanner(os.Stdin),
		w:  bufio.NewWriter(os.Stdout),
		tw: tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)}
}

// Run runs application command line user interface.
func (a *App) Run(ctx context.Context) error {
	a.auth(ctx)
	a.main(ctx)
	return nil
}

const (
	success = 0
	failure = 1
)

func (a *App) exit(msg string, code int) {
	if _, err := fmt.Fprintln(a.w, msg); err != nil {
		logger.Log.Error("Internal error", zap.Error(err))
		os.Exit(failure)
	}
	if err := a.w.Flush(); err != nil {
		logger.Log.Error("Internal error", zap.Error(err))
		os.Exit(failure)
	}
	os.Exit(code)
}

func (a *App) exitSuccess() { a.exit("Exit", success) }

func (a *App) exitInternal(err error) {
	logger.Log.Error("Internal error", zap.Error(err))
	a.exit("Internal error", failure)
}

func (a *App) exitInvalidInput() { a.exit("Invalid input", failure) }

func (a *App) exitTryAgain() { a.exit("Try again", failure) }

func (a *App) exitTimeout() { a.exit("Sign in timeout", failure) }

func (a *App) screen(s string) {
	if _, err := fmt.Fprint(a.w, s); err != nil {
		a.exitInternal(err)
	}
	if err := a.w.Flush(); err != nil {
		a.exitInternal(err)
	}
}

func (a *App) newLine() { a.screen("\n") }

const (
	signIn = 1
	signUp = 2
)

func (a *App) auth(ctx context.Context) {
	a.screen("AUTHORIZATION SCREEN\n\n1. SIGN IN\n2. SIGN UP\n\n")
	c := a.choice("Option number: ")
	email, pass := a.sign(c)
	switch c {
	case signIn:
		a.signIn(ctx, email, pass)
	case signUp:
		a.signUp(ctx, email, pass)
	default:
		a.exitInvalidInput()
	}
	a.newLine()
}

func (a *App) sign(c int) (email, pass string) {
	if c == signIn || c == signUp {
		email = a.input("Email: ")
		if _, err := mail.ParseAddress(email); err != nil {
			a.exit("Invalid email", failure)
		}
		pass = a.input("Password: ")
		if c == signUp {
			if pass != a.input("Repeat password: ") {
				a.exit("Passwords mismatch", failure)
			}
		}
	}
	return
}

func (a *App) signIn(ctx context.Context, email, pass string) {
	if err := a.s.Login(ctx, email, pass); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unavailable:
				a.screen("\nService unavailable\n")
			case codes.InvalidArgument:
				a.exitInvalidInput()
			case codes.Unauthenticated:
				a.exitTryAgain()
			default:
				a.exitTryAgain()
			}
		}
	}
}

func (a *App) signUp(ctx context.Context, email, pass string) {
	if err := a.s.Register(ctx, email, pass); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unavailable:
				a.exit("Service unavailable", failure)
			case codes.InvalidArgument:
				a.exitInvalidInput()
			case codes.Unauthenticated:
				a.exitTryAgain()
			case codes.AlreadyExists:
				a.exit("Already signed up", failure)
			default:
				a.exitTryAgain()
			}
		}
	} else {
		_, err := fmt.Fprintf(a.w, "Signed up: %s\n", email)
		if err != nil {
			a.exitInternal(err)
		}
		if err := a.w.Flush(); err != nil {
			a.exitInternal(err)
		}
	}
}

const (
	print = 1
	enter = 2
)

func (a *App) main(ctx context.Context) {
	if len(a.cfg.Token) > 0 {
		a.screen("MAIN SCREEN ONLINE ACCESS\n\n")
		logger.Log.Info("running sync...", zap.String("address", a.cfg.AppAddr))
		go a.sync(ctx)
	} else {
		a.screen("MAIN SCREEN OFFLINE ACCESS\n\n")
	}
	a.screen("1. PRINT DATA\n2. ENTER DATA\n\n")
	switch a.choice("Option number: ") {
	case print:
		a.print(ctx)
	case enter:
		a.enter(ctx)
	default:
		a.exitInvalidInput()
	}
	if a.choice("Done. 1 to continue or 2 to quit: ") == print {
		a.main(ctx)
	}
	a.exitSuccess()
}

func (a *App) input(ask string) string {
	a.screen(ask)
	a.r.Scan()
	if err := a.r.Err(); err != nil {
		a.exitInvalidInput()
	}
	return a.r.Text()
}

func (a *App) choice(ask string) int {
	a.screen(ask)
	a.r.Scan()
	if err := a.r.Err(); err != nil {
		a.exitInvalidInput()
	}
	c, err := strconv.Atoi(string(a.r.Text()))
	if err != nil {
		a.exitInvalidInput()
	}
	return c
}

func (a *App) index(ctx context.Context) {
	i, err := a.s.IndexUserData(ctx)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				if a.choice("\nNo data. 1 to continue or 2 to quit: ") == print {
					a.newLine()
					a.main(ctx)
				}
				a.exitSuccess()
			}
		}
		a.exitInternal(err)
	}
	if _, err := fmt.Fprintln(a.tw, "ID\tTYPE\tUPDATED\t"); err != nil {
		a.exitInternal(err)
	}
	for _, v := range i {
		if _, err := fmt.Fprintln(a.tw, v.ID+"\t"+v.Data+"\t"+
			v.Updated.Format(time.DateTime)+"\t"); err != nil {
			a.exitInternal(err)
		}
	}
	if err := a.tw.Flush(); err != nil {
		a.exitInternal(err)
	}
	a.newLine()
}

const (
	pass = 1
	text = 2
	file = 3
	card = 4
)

func (a *App) inputRepeat(ctx context.Context) {
	a.screen("\nID must be from the list:\n\n")
	a.print(ctx)
}

func (a *App) print(ctx context.Context) {
	a.index(ctx)
	id := a.input("ID to print: ")
	if len(id) == 0 {
		a.inputRepeat(ctx)
		return
	}
	d, err := a.s.PrintUserData(ctx, id)
	if err != nil {
		a.inputRepeat(ctx)
		return
	}
	t, c := d.Type, d.Data
	j, err := a.c.Dec(c, a.cfg.Key)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.Unauthenticated {
				a.exitTimeout()
			}
		}
		a.exitInternal(err)
	}
	a.newLine()
	switch t {
	case pass:
		a.ppass(j)
	case text:
		a.ptext(j)
	case file:
		a.pfile(j)
	case card:
		a.pcard(j)
	default:
		a.exitInvalidInput()
	}
	if err := a.tw.Flush(); err != nil {
		a.exitInternal(err)
	}
	a.newLine()
}

func unmarshal[T any](j string) (*T, error) {
	x := new(T)
	if err := json.Unmarshal([]byte(j), x); err != nil {
		return nil, err
	}
	return x, nil
}

const wait = 10 * time.Second

func (a *App) enter(ctx context.Context) {
	id := a.input("ID to enter: ")
	if len(id) == 0 {
		a.exit("ID empty", failure)
	}
	if _, err := a.s.PrintUserData(ctx, id); err == nil {
		a.index(ctx)
		a.screen("ID already exists\n")
		a.enter(ctx)
		return
	}
	a.screen("\n1. PASS\n2. TEXT\n3. FILE\n4. CARD\n\n")
	t := a.choice("Option number: ")
	var j []byte
	switch t {
	case pass:
		j = a.pass()
	case text:
		j = a.text()
	case file:
		j = a.file()
	case card:
		j = a.card()
	default:
		a.exitInvalidInput()
	}

	c, err := a.c.Enc(string(j), a.cfg.Key)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.Unauthenticated {
				a.exitTimeout()
			}
		}
		a.exitInternal(err)
	}
	if err := a.s.EnterUserData(ctx, id, t, c); err != nil {
		a.exitInternal(err)
	}
	go func() {
		if len(a.cfg.Token) > 0 {
			a.upload(ctx)
			time.Sleep(wait) // to merge data on server
			a.download(ctx)
		}
	}()
	a.newLine()
}
