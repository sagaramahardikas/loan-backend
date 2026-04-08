package main

import (
	"example.com/loan/cmd/cli/command"
	"github.com/alecthomas/kong"
)

type CLI struct {
	CreateUser            command.CreateUserCmd            `cmd:"" help:"Create user command, script to create user"`
	CreateAccount         command.CreateAccountCmd         `cmd:"" help:"Create account command, script to create account"`
	CreateLoan            command.CreateLoanCmd            `cmd:"" help:"Create loan command, script to create loan"`
	ForceDisburseLoan     command.ForceDisburseLoanCmd     `cmd:"" help:"Force disburse loan command, script to force disburse loan"`
	OverdueBillingChecker command.OverdueBillingCheckerCmd `cmd:"" help:"Overdue billing checker command, should be run by daily cron"`
}

func main() {
	cli := CLI{}
	cmd := kong.Parse(&cli)

	err := cmd.Run()
	cmd.FatalIfErrorf(err)
}
