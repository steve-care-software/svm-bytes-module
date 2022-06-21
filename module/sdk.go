package module

import (
	"errors"
	"fmt"

	selector_application "github.com/steve-care-software/selector/applications"
	"github.com/steve-care-software/svm/domain/interpreters"
	validator_application "github.com/steve-care-software/validator/applications"
)

// NewModule creates a new module instance
func NewModule() interpreters.Module {
	// create the validator:
	validator := `
		%rootToken;
		-space;
		-endOfLine;

		rootToken: .bytes
			     | .byte
				 ;

		bytes: .openSquareBracket .byteWithSemiColon[1,] .closeSquareBracket;
		byteWithSemiColon: .byte .semiColon;
		byte: .dollar .number[1,3];

		number: .zero
			  | .one
			  | .two
			  | .three
			  | .four
			  | .five
			  | .six
			  | .seven
			  | .height
			  | .nine
			  ;

		openSquareBracket: $91;
		closeSquareBracket: $93;
		semiColon: $59;
		dollar: $36;
		zero: $48;
		one: $49;
		two: $50;
		three: $51;
		four: $52;
		five: $53;
		six: $54;
		seven: $55;
		height: $56;
		nine: $57;
		space: $32;
		endOfLine: $10;
	`

	validatorApplication := validator_application.NewApplication()
	validatorGrammar, err := validatorApplication.Compile(validator)
	if err != nil {
		panic(err)
	}

	// create the byte selector:
	selectorApp := selector_application.NewApplication()
	byteSelector, err := selectorApp.Compile("+ @rootToken @byte .number")
	if err != nil {
		panic(err)
	}

	// create the bytes selector:
	bytesSelector, err := selectorApp.Compile("+ @rootToken @bytes @byteWithSemiColon @byte .number")
	if err != nil {
		panic(err)
	}

	name := "bytes"
	executeFn := func(input map[string]string, application string) (string, error) {
		if amount, ok := input["amount"]; ok {
			if data, ok := input["data"]; ok {

				amountBytes, err := extract(amount, validatorApplication, validatorGrammar, selectorApp, byteSelector)
				if err != nil {
					return "", err
				}

				dataBytes, err := extract(data, validatorApplication, validatorGrammar, selectorApp, bytesSelector)
				if err != nil {
					return "", err
				}

				fmt.Printf("\n++++%v\n", amountBytes)
				fmt.Printf("\n++++%v\n", dataBytes)

				return "", nil
			}

			str := fmt.Sprintf("the 'data' variable is undeclared while executing: %s.%s", name, application)
			return "", errors.New(str)
		}

		str := fmt.Sprintf("the 'amount' variable is undeclared while executing: %s.%s", name, application)
		return "", errors.New(str)
	}

	ins, err := interpreters.NewModuleBuilder().Create().WithName(name).WithEvent(executeFn).Now()
	if err != nil {
		panic(err)
	}

	return ins
}
