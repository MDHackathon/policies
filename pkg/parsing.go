package policies

import (
	"encoding/json"
	// TODO use real log library
	"errors"
	"io/ioutil"
	"net/url"
	"strings"
)

// TODO Develop struct to reduce the string everywhere
// policies define the input s3 storare. the output s3 storage
// and the associatate credentials, with also the rules and the operationn
// to do
type policies struct {
	// s3 intput storage
	InEndpoint string `json:"in_endpoint"`
	// prefix
	InPath string `json:"in_path"`
	// bucket source
	InBucket string `json:"in_bucket"`
	// access key for the s3 input
	InKey string `json:"in_key"`
	// secret key for the s3 input
	InSecret string `json:"in_secret"`

	// s3 output storage
	OutEndpoint string `json:"out_endpoint"`
	// prefix
	OutPath string `json:"out_path"`
	// bucket destination in the
	OutBucket string `json:"out_bucket"`
	// access key for the s3 output
	OutKey string `json:"out_key"`
	// secret key for the s3 output
	OutSecret string `json:"out_secret"`

	// match or date
	RuleType string `json:"rule_type"`
	// > >= = < <= !=
	RuleCmpOp string `json:"rule_compare_op"`
	// use to match or as a date
	RuleValue string `json:"rule_value"`

	Parallelisation bool `json:"parallelisation"`
	// move copy
	Operation          string `json:"operations"`
	CustumFunctionName string `json:"custom_function_name"`
}

// Ensure that endpoint are http url
func validateEndpoint(endpoint string) (err error) {
	_, err = url.ParseRequestURI(endpoint)
	return
}

// Ensure that the rule type is implemented
func validateRuleType(rule string) (err error) {
	if strings.Compare(rule, "match") != 0 {
		return
	}
	if strings.Compare(rule, "date") != 0 {
		return
	}
	err = errors.New("Unknow rule: use match or date")
	return
}

// Checn values
func validateJSONInput(p policies) (err error) {
	if err = validateEndpoint(p.InEndpoint); err != nil {
		return
	}
	if err = validateEndpoint(p.OutEndpoint); err != nil {
		return
	}
	if err = validateRuleType(p.RuleType); err != nil {
		return
	}
	return
}

// Loading policies from
func LoadPolicieFromPath(path string) (p policies, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		panic(err)
	}
	p, err = LoadPolicieFromByte(data)
	return
}

func LoadPolicieFromByte(data []byte) (p policies, err error) {
	err = json.Unmarshal(data, &p)
	return
}
