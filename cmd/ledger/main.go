package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/howeyc/cli-ledger/pkg/ledger"
)

func main() {
	var startDate, endDate time.Time
	startDate = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
	endDate = time.Now()
	var startString, endString string
	var columnWidth, transactionDepth int
	var showEmptyAccounts bool
	var columnWide bool

	var ledgerFileName string

	ledger.TransactionDateFormat = "2006/01/02"
	TransactionDateFormat := ledger.TransactionDateFormat

	flag.StringVar(&ledgerFileName, "f", "", "Ledger file name (*Required).")
	flag.StringVar(&startString, "s", startDate.Format(TransactionDateFormat), "Start date of transaction processing.")
	flag.StringVar(&endString, "e", endDate.Format(TransactionDateFormat), "End date of transaction processing.")
	flag.BoolVar(&showEmptyAccounts, "empty", false, "Show empty (zero balance) accounts.")
	flag.IntVar(&transactionDepth, "depth", -1, "Depth of transaction output (balance).")
	flag.IntVar(&columnWidth, "columns", 80, "Set a column width for output.")
	flag.BoolVar(&columnWide, "wide", false, "Wide output (same as --columns=132).")
	flag.Parse()

	if columnWidth == 80 && columnWide {
		columnWidth = 132
	}

	if len(ledgerFileName) == 0 {
		flag.Usage()
		return
	}

	// TODO(chris): Handle parse of arg errors
	startDate, _ = time.Parse(TransactionDateFormat, startString)
	endDate, _ = time.Parse(TransactionDateFormat, endString)

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Specify a command.")
		return
	}

	ledgerFileReader, err := os.Open(ledgerFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ledgerFileReader.Close()

	generalLedger, parseError := ledger.ParseLedger(ledgerFileReader)
	if parseError != nil {
		fmt.Println(parseError)
		return
	}

	containsFilterArray := args[1:]
	switch strings.ToLower(args[0]) {
	case "balance", "bal":
		ledger.PrintBalances(ledger.GetBalances(generalLedger, containsFilterArray), showEmptyAccounts, transactionDepth, columnWidth)
	case "print":
		ledger.PrintLedger(os.Stdout, generalLedger, columnWidth)
	case "register", "reg":
		ledger.PrintRegister(generalLedger, containsFilterArray, columnWidth)
	}
}