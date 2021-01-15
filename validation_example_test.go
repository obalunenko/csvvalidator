package csvvalidator_test

import (
	"fmt"

	"github.com/obalunenko/csvvalidator"
)

func ExampleValidationRules_ValidateRow() {
	// initialise columns
	var (
		clientNameCol        = csvvalidator.NewColumn(0, "Client Name")
		transactionDateCol   = csvvalidator.NewColumn(1, "Transaction Date")
		transactionAmountCol = csvvalidator.NewColumn(2, "Transaction Amount")
		debitCreditCodeCol   = csvvalidator.NewColumn(3, "Debit Credit Code")
		customerIDCol        = csvvalidator.NewColumn(4, "Customer ID")

		// set total columns number
		totalColNum = uint(5)
	)

	// set up validation rules
	var rowRules = csvvalidator.ValidationRules{
		clientNameCol: csvvalidator.Rule{
			MaxLength:       4,
			MinLength:       4,
			RestrictedChars: nil,
		},
		transactionDateCol: csvvalidator.Rule{
			MaxLength:       10,
			MinLength:       10,
			RestrictedChars: nil,
		},
		transactionAmountCol: csvvalidator.Rule{
			MaxLength:       14,
			MinLength:       4,
			RestrictedChars: nil,
		},
		debitCreditCodeCol: csvvalidator.Rule{
			MaxLength:       1,
			MinLength:       1,
			RestrictedChars: nil,
		},
		customerIDCol: csvvalidator.Rule{
			MaxLength:       10,
			MinLength:       1,
			RestrictedChars: nil,
		},
	}

	// create csv row ([]string)
	var testDataRow = make([]string, totalColNum)
	testDataRow[clientNameCol.Number] = "Test Client Name too long" // exceed maxlength
	testDataRow[transactionDateCol.Number] = "2019/05/30"
	testDataRow[transactionAmountCol.Number] = "100.45"
	testDataRow[debitCreditCodeCol.Number] = "C"
	testDataRow[customerIDCol.Number] = "ID"

	// validate row
	err := rowRules.ValidateRow(testDataRow)
	fmt.Print(err)
	// Output:
	// invalid column [row:[Test Client Name too long 2019/05/30 100.45 C ID]; column:0:Client Name]: invalid length [column:[Test Client Name too long]; has len:[25]; should be:[4]]
}

func ExampleNewColumn() {
	clientNameCol := csvvalidator.NewColumn(0, "Client Name")
	fmt.Print(clientNameCol.String())
	// Output:
	// 0:Client Name
}
