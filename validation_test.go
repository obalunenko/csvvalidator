package csvvalidator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateColumn(t *testing.T) {
	type args struct {
		data string
		rule Rule
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "all rules passed",
			args: args{
				data: "1234567890",
				rule: Rule{
					MaxLength: 10,
					MinLength: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "max length exceed",
			args: args{
				data: "12345678901",
				rule: Rule{
					MaxLength: 10,
					MinLength: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "empty when not allowed",
			args: args{
				data: "",
				rule: Rule{
					MaxLength: 10,
					MinLength: 2,
				},
			},
			wantErr: true,
		},
		{
			name: "min length exceed",
			args: args{
				data: "1",
				rule: Rule{
					MaxLength: 10,
					MinLength: 2,
				},
			},
			wantErr: true,
		},
		{
			name: "resrticted chars not contain",
			args: args{
				data: "2 data valid",
				rule: Rule{
					MaxLength: 100,
					MinLength: 2,
					RestrictedChars: []string{",", "`", "~", "!", "@", "#", "$", "%", "^", "&", "*", "_", "{", "}",
						"<", ">", "[", "]", "=", `\`, ";"},
				},
			},
			wantErr: false,
		},
		{
			name: "resrticted chars  contain",
			args: args{
				data: "2 data not~valid",
				rule: Rule{
					MaxLength: 100,
					MinLength: 2,
					RestrictedChars: []string{",", "`", "~", "!", "@", "#", "$", "%", "^", "&", "*", "_", "{", "}",
						"<", ">", "[", "]", "=", `\`, ";"},
				},
			},
			wantErr: true,
		},
		{
			name: "max length not set up",
			args: args{
				data: "2 data valid",
				rule: Rule{
					MaxLength: 0,
					MinLength: 2,
					RestrictedChars: []string{",", "`", "~", "!", "@", "#", "$", "%", "^", "&", "*", "_", "{", "}",
						"<", ">", "[", "]", "=", `\`, ";"},
				},
			},
			wantErr: false,
		},
		{
			name: "max and min lengths are 0",
			args: args{
				data: "2 data valid",
				rule: Rule{
					MaxLength: 0,
					MinLength: 0,
					RestrictedChars: []string{",", "`", "~", "!", "@", "#", "$", "%", "^", "&", "*", "_", "{", "}",
						"<", ">", "[", "]", "=", `\`, ";"},
				},
			},
			wantErr: true,
		},
		{
			name: "check defined rule RuleShouldBeEmpty - error case",
			args: args{
				data: "2 data valid",
				rule: RuleShouldBeEmpty,
			},
			wantErr: true,
		},
		{
			name: "check defined rule RuleShouldBeEmpty - ok",
			args: args{
				data: "",
				rule: RuleShouldBeEmpty,
			},
			wantErr: false,
		},
		{
			name: "check defined rule RuleNotEmpty - error case",
			args: args{
				data: "",
				rule: RuleNotEmpty,
			},
			wantErr: true,
		},
		{
			name: "check defined rule RuleNotEmpty - ok",
			args: args{
				data: "2 data valid",
				rule: RuleNotEmpty,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateColumn(tt.args.data, tt.args.rule)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestValidationRules_ValidateRow(t *testing.T) {

	rules := ValidationRules{
		Column{Number: 0}: {
			MaxLength: 10,
			MinLength: 3,
		},
		Column{Number: 1}: {
			MaxLength: 8,
			MinLength: 0,
		},
		Column{Number: 2}: {
			MaxLength: 10,
			MinLength: 3,
		},
		Column{Number: 3}: {
			MaxLength: 10,
			MinLength: 2,
		},
	}
	type args struct {
		rec []string
	}
	tests := []struct {
		name    string
		rules   ValidationRules
		args    args
		wantErr bool
	}{
		// columns count start from 0
		{
			name:  "all columns correct",
			rules: rules,
			args: args{
				rec: []string{"1234567890", "", "123", "12", "", ""},
			},
			wantErr: false,
		},
		{
			name:  "column 1 exceed max length",
			rules: rules,
			args: args{
				rec: []string{"1234567890", "123456789", "123", "", "", ""},
			},
			wantErr: true,
		},
		{
			name:  "column 2 exceed min length",
			rules: rules,
			args: args{
				rec: []string{"1234567890", "123456", "12", "", "", ""},
			},
			wantErr: true,
		},
		{
			name:  "nil rules",
			rules: nil,
			args: args{
				rec: []string{"1234567890", "123456", "123", "", "", ""},
			},
			wantErr: true,
		},
		{
			name:  "col 1 empty and min is 0",
			rules: nil,
			args: args{
				rec: []string{"1234567890", "", "123", "", "", ""},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rules.ValidateRow(tt.args.rec)
			switch tt.wantErr {
			case false:
				assert.NoError(t, err)
			case true:
				assert.Error(t, err)
			}
		})
	}
}

func TestRow_ValidateRow(t *testing.T) {
	rules := ValidationRules{
		Column{Number: 0}: {
			MaxLength: 10,
			MinLength: 3,
		},
		Column{Number: 1}: {
			MaxLength: 8,
			MinLength: 0,
		},
		Column{Number: 2}: {
			MaxLength: 10,
			MinLength: 3,
		},
		Column{Number: 3}: {
			MaxLength: 10,
			MinLength: 2,
		},
	}
	type fields struct {
		ColumnsTotalNum uint
		ColumnsRules    ValidationRules
	}
	type args struct {
		rec []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "all rules pass",
			fields: fields{
				ColumnsTotalNum: 6,
				ColumnsRules:    rules,
			},
			args: args{
				rec: []string{"1234567890", "", "123", "12", "", ""},
			},
			wantErr: false,
		},
		{
			name: "error when columns number invalid",
			fields: fields{
				ColumnsTotalNum: 6,
				ColumnsRules:    rules,
			},
			args: args{
				rec: []string{"1234567890", "", "123", "12"},
			},
			wantErr: true,
		},
		{
			name: "error when columns number not specified in rule",
			fields: fields{
				ColumnsRules: rules,
			},
			args: args{
				rec: []string{"1234567890", "", "123", "12"},
			},
			wantErr: true,
		},
		{
			name: "error when columns rule not passed",
			fields: fields{
				ColumnsTotalNum: 6,
				ColumnsRules:    rules,
			},
			args: args{
				rec: []string{"1234567890", "123456789", "123", "", "", ""},
			},
			wantErr: true,
		},
		{
			name: "error when column rules not specified",
			fields: fields{
				ColumnsTotalNum: 6,
				ColumnsRules:    nil,
			},
			args: args{
				rec: []string{"1234567890", "123456789", "123", "", "", ""},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := Row{
				ColumnsTotalNum: tt.fields.ColumnsTotalNum,
				ColumnsRules:    tt.fields.ColumnsRules,
			}
			if err := row.ValidateRow(tt.args.rec); (err != nil) != tt.wantErr {
				t.Errorf("Row.ValidateRow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewColumn(t *testing.T) {
	type args struct {
		number uint
		name   []string
	}
	tests := []struct {
		name string
		args args
		want Column
	}{
		{
			name: "construct column with empty name",
			args: args{
				number: 12,
				name:   nil,
			},
			want: Column{
				Number: 12,
				Name:   "",
			},
		},
		{
			name: "construct column with multiple name parameters",
			args: args{
				number: 13,
				name:   []string{"Column", "TEST", "this"},
			},
			want: Column{
				Number: 13,
				Name:   "Column TEST this",
			},
		},
		{
			name: "construct column with 1 name parameter",
			args: args{
				number: 13,
				name:   []string{"Column"},
			},
			want: Column{
				Number: 13,
				Name:   "Column",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewColumn(tt.args.number, tt.args.name...)
			assert.Equal(t, tt.want, got)
		})
	}
}
