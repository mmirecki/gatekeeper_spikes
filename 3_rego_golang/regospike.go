package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/open-policy-agent/opa/rego"
)

func readFile(file string) []byte {
	b, err := ioutil.ReadFile(file) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	return b

}

func main() {
	ctx := context.TODO()
	regoModule := string(readFile("namespace_rule.rego"))
	inputJson := readFile("namespace_input")

	input := &map[string]interface{}{}
	json.Unmarshal(inputJson, input)

	query, err := rego.New(
		rego.Query("x = data.regospike.ns_mutated(input)"),
		rego.Module("regospike", regoModule),
	).PrepareForEval(ctx)

	if err != nil {
		fmt.Println(err)
		return
	}
	results, err := query.Eval(ctx, rego.EvalInput(input))

	fmt.Println(results)
}
