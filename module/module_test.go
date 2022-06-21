package module

import (
	"fmt"
	"testing"

	svm_application "github.com/steve-care-software/svm/applications"
	"github.com/steve-care-software/svm/domain/interpreters"
)

func TestModule_isSuccess(t *testing.T) {
	module := NewModule()
	modules, err := interpreters.NewModulesBuilder().WithList([]interpreters.Module{
		module,
	}).Now()

	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	application, err := svm_application.NewBuilder(createScreenWriterForTests()).Create().WithModules(modules).Now()
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	script := `
        // declare the module and its types;;
        module bytes;;
        type data bytes.byte;;
        type data bytes.bytes;;
        type application bytes.leftShift;;

        // declare the parameters;;
        -> bytes.bytes $input;;
        <- bytes.bytes $output;;

        // declare the application;;
        bytes.leftShift $leftShiftApp;;

        // declare the amount to shift;;
        bytes.byte $amount = $1;;

        // atachthe data;;
        attach input:data @ leftShiftApp;;

        // attach the amount;;
        attach amount:amount @ leftShiftApp;;

        // execute;;
        $output = execute leftShiftApp;;
	`

	program, remaining, err := application.Compile(script)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	fmt.Printf("\n%s\n", remaining)

	_, err = application.Execute(map[string]string{
		"input": "[$255; $255; $12;]",
	}, program)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

}
