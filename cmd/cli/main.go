package main

import (
	"example.com/loan/cmd/cli/command"
	"github.com/alecthomas/kong"
)

type CLI struct {
	CreateUser        command.CreateUserCmd        `cmd:"" help:"Create user command"`
	CreateAccount     command.CreateAccountCmd     `cmd:"" help:"Create account command"`
	CreateLoan        command.CreateLoanCmd        `cmd:"" help:"Create loan command"`
	ForceDisburseLoan command.ForceDisburseLoanCmd `cmd:"" help:"Force disburse loan command"`
}

func main() {
	cli := CLI{}
	cmd := kong.Parse(&cli)

	err := cmd.Run()
	cmd.FatalIfErrorf(err)
}
