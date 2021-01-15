// Package csvvalidator provide validation of csv rows according to set up rules.
package csvvalidator

import (
	"errors"
	"fmt"
	"strings"
)

// Column represents csv column.
type Column struct {
	Number uint
	Name   string
}

// NewColumn creates new Column with passed column number,
// name of column is optional.
func NewColumn(number uint, name ...string) Column {
	var resName string
	if len(name) != 0 {
		resName = strings.Join(name, " ")
	}
	return Column{
		Number: number,
		Name:   resName,
	}

}

func (c *Column) String() string {
	if c.Name == "" {
		return fmt.Sprintf("%d", c.Number)
	}
	return fmt.Sprintf("%d:%s", c.Number, c.Name)
}

// Rule represent column row validation rule.
type Rule struct {
	MaxLength       uint // if 0 - no limit
	MinLength       uint // if 0 - no limit
	RestrictedChars []string
}

var (
	// RuleNotEmpty defines a rule for field that require it to be not empty.
	RuleNotEmpty = Rule{
		MaxLength:       0,
		MinLength:       1,
		RestrictedChars: nil,
	}
	// RuleShouldBeEmpty defines a rule for field that require it to be empty.
	RuleShouldBeEmpty = Rule{
		MaxLength:       0,
		MinLength:       0,
		RestrictedChars: nil,
	}
)

// ValidationRules represent rules for columns validation,
// counting of columns in row starts from 0.
type ValidationRules map[Column]Rule

// Row represent csv row:
// - total columns number for row
// - validation rules for each column.
type Row struct {
	ColumnsTotalNum uint
	ColumnsRules    ValidationRules
}

// ValidateRow validates passed row.
func (row Row) ValidateRow(rec []string) error {
	if row.ColumnsTotalNum == 0 {
		return errors.New("invalid columns total number specified in rules")
	}

	if uint(len(rec)) != row.ColumnsTotalNum {
		return fmt.Errorf("invalid row [row:[%v]; columns num:[%d]; should be: [%d]]",
			rec,  len(rec), row.ColumnsTotalNum)
	}

	return row.ColumnsRules.ValidateRow(rec)
}

// ValidateRow validates passed row,
// counting of columns in row starts from 0.
func (rules ValidationRules) ValidateRow(rec []string) error {
	if len(rules) == 0 {
		return errors.New("validation rules not specified")
	}
	if len(rec) == 0 {
		return errors.New("empty row")
	}

	for col, rule := range rules {
		if err := validateColumn(rec[col.Number], rule); err != nil {
			return fmt.Errorf("invalid column [row:%s; column:%s]: %w", rec,col.String(), err)
		}
	}

	return nil
}

func validateColumn(data string, rule Rule) error {
	if rule.MinLength != 0 && data == "" {
		return errors.New("should not be empty")
	}

	if rule.MinLength == rule.MaxLength {
		if uint(len(data)) != rule.MinLength {
			return fmt.Errorf("invalid length [column:[%s]; has len:[%d]; should be:[%d]]",
				data,len(data),rule.MinLength)
		}
	}

	if uint(len(data)) < rule.MinLength {
		return fmt.Errorf("invalid length [column:[%s]; has len:[%d]; should be at less[%d]]",
			data, len(data),rule.MinLength)
	}

	if uint(len(data)) > rule.MaxLength && rule.MaxLength != 0 {
		return fmt.Errorf("invalid length [column:[%s]; has len:[%d]; should be at less[%d]]",
			data, len(data),rule.MaxLength)
	}

	if rule.RestrictedChars != nil {
		if strings.ContainsAny(data, strings.Join(rule.RestrictedChars, "")) {
			return fmt.Errorf("contains restricted characters [column:[%s]]", data)
		}
	}

	return nil
}
