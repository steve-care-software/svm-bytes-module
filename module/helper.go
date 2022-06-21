package module

import (
	selector_application "github.com/steve-care-software/selector/applications"
	"github.com/steve-care-software/selector/domain/selectors"
	validator_application "github.com/steve-care-software/validator/applications"
	"github.com/steve-care-software/validator/domain/grammars"
)

func extract(
	content string,
	validatorApp validator_application.Application,
	validatorGrammar grammars.Grammar,
	selectorApp selector_application.Application,
	selector selectors.Selector,
) ([]byte, error) {

	result, err := validatorApp.Execute(validatorGrammar, []byte(content), false)
	if err != nil {
		return nil, err
	}

	return selectorApp.Execute(selector, result)
}
