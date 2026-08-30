package aip160

import "github.com/uthereal/scheme-runtime-go/pkg/contract"

// operatorMap is a map from operators to their allowed status.
type operatorMap map[contract.Aip160Operator]bool

// mapFieldTypeToAllowedOperators defines the supported operators for each
// field type.
var mapFieldTypeToAllowedOperators = map[contract.Aip160FieldType]operatorMap{
	contract.TypeBool: {
		contract.OpEq:      true,
		contract.OpNotEq:   true,
		contract.OpDefault: true,
	},
	contract.TypeDecimal: {
		contract.OpEq:      true,
		contract.OpNotEq:   true,
		contract.OpGt:      true,
		contract.OpGtEq:    true,
		contract.OpLt:      true,
		contract.OpLtEq:    true,
		contract.OpDefault: true,
	},
	contract.TypeDuration: {
		contract.OpEq:      true,
		contract.OpNotEq:   true,
		contract.OpGt:      true,
		contract.OpGtEq:    true,
		contract.OpLt:      true,
		contract.OpLtEq:    true,
		contract.OpDefault: true,
	},
	contract.TypeEnum: {
		contract.OpEq:      true,
		contract.OpNotEq:   true,
		contract.OpDefault: true,
	},
	contract.TypeFloat: {
		contract.OpEq:      true,
		contract.OpNotEq:   true,
		contract.OpGt:      true,
		contract.OpGtEq:    true,
		contract.OpLt:      true,
		contract.OpLtEq:    true,
		contract.OpDefault: true,
	},
	contract.TypeInt: {
		contract.OpEq:      true,
		contract.OpNotEq:   true,
		contract.OpGt:      true,
		contract.OpGtEq:    true,
		contract.OpLt:      true,
		contract.OpLtEq:    true,
		contract.OpDefault: true,
	},
	contract.TypeString: {
		contract.OpEq:      true,
		contract.OpNotEq:   true,
		contract.OpHas:     true,
		contract.OpDefault: true,
	},
	contract.TypeTimestamp: {
		contract.OpEq:      true,
		contract.OpNotEq:   true,
		contract.OpGt:      true,
		contract.OpGtEq:    true,
		contract.OpLt:      true,
		contract.OpLtEq:    true,
		contract.OpDefault: true,
	},
	contract.TypeUUID: {
		contract.OpEq:      true,
		contract.OpNotEq:   true,
		contract.OpDefault: true,
	},
}

// HasOperator checks if the given FieldType supports the specified operator.
func HasOperator(t contract.Aip160FieldType, op contract.Aip160Operator) bool {
	mapOperatorToAllowed, ok := mapFieldTypeToAllowedOperators[t]
	if !ok {
		return false
	}

	allowed, ok := mapOperatorToAllowed[op]
	if !ok {
		return false
	}

	return allowed
}
