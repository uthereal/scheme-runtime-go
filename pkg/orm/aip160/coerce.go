package aip160

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"go.chromium.org/luci/common/data/aip160"
)

// CoerceArgToComparable extracts the Comparable node safely.
func CoerceArgToComparable(arg *aip160.Arg) (*aip160.Comparable, error) {
	if arg.Composite != nil {
		return nil, errors.New(
			"composite expressions in arguments not supported yet",
		)
	}
	if arg.Comparable == nil {
		return nil, errors.New("missing comparable in argument")
	}
	return arg.Comparable, nil
}

// CoerceArgToStringConstant attempts to extract a string constant from an Arg.
func CoerceArgToStringConstant(arg *aip160.Arg) (string, error) {
	cmp, err := CoerceArgToComparable(arg)
	if err != nil {
		return "", err
	}
	if cmp.Member == nil {
		return "", errors.New("invalid string comparable -> member missing")
	}
	if cmp.Member.Value == nil {
		return "", errors.New("invalid string comparable -> value missing")
	}
	return cmp.Member.Value.Value, nil
}

// CoerceArgToFloatConstant attempts to extract a float constant from
// an Arg AST node. The AST node must be an unquoted float value.
func CoerceArgToFloatConstant(arg *aip160.Arg) (float64, error) {
	cmp, err := CoerceArgToComparable(arg)
	if err != nil {
		return 0, fmt.Errorf("failed to coerce arg -> %w", err)
	}

	if cmp.Member == nil {
		return 0, errors.New("invalid float comparable -> member missing")
	}

	if cmp.Member.Value == nil {
		return 0, errors.New("invalid float comparable -> value missing")
	}

	if cmp.Member.Value.Quoted {
		return 0, fmt.Errorf(
			"expected an unquoted float literal but found string %q",
			cmp.Member.Value.Value,
		)
	}

	rawVal := cmp.Member.Input()
	if rawVal == "" {
		return 0, errors.New("empty float literal")
	}

	val, err := strconv.ParseFloat(rawVal, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to parse float literal %q -> %w",
			rawVal,
			err,
		)
	}

	return val, nil
}

// CoerceArgToIntConstant attempts to extract an integer constant from an Arg.
func CoerceArgToIntConstant(arg *aip160.Arg) (int64, error) {
	cmp, err := CoerceArgToComparable(arg)
	if err != nil {
		return 0, fmt.Errorf("failed to coerce arg -> %w", err)
	}

	if cmp.Member == nil {
		return 0, errors.New("invalid int comparable -> member missing")
	}

	if cmp.Member.Value == nil {
		return 0, errors.New("invalid int comparable -> value missing")
	}

	if cmp.Member.Value.Quoted {
		return 0, fmt.Errorf(
			"expected an unquoted int literal but found string %q",
			cmp.Member.Value.Value,
		)
	}

	rawVal := cmp.Member.Input()
	if rawVal == "" {
		return 0, errors.New("empty int literal")
	}

	val, err := strconv.ParseInt(rawVal, 10, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to parse int literal %q -> %w",
			rawVal,
			err,
		)
	}

	return val, nil
}

// CoerceArgToBoolConstant attempts to extract a boolean constant from an Arg.
func CoerceArgToBoolConstant(arg *aip160.Arg) (bool, error) {
	cmp, err := CoerceArgToComparable(arg)
	if err != nil {
		return false, fmt.Errorf("failed to coerce arg -> %w", err)
	}

	if cmp.Member == nil {
		return false, errors.New("invalid bool comparable -> member missing")
	}

	rawVal := cmp.Member.Input()
	if rawVal == "true" {
		return true, nil
	}
	if rawVal == "false" {
		return false, nil
	}
	return false, fmt.Errorf("invalid boolean value %q", rawVal)
}

// CoerceArgToTimestampConstant attempts to extract a time.Time constant from
// an Arg AST node. The AST node must be a double-quoted RFC3339 string.
func CoerceArgToTimestampConstant(arg *aip160.Arg) (time.Time, error) {
	cmp, err := CoerceArgToComparable(arg)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to coerce arg -> %w", err)
	}

	if cmp.Member == nil {
		return time.Time{}, errors.New(
			"invalid timestamp comparable -> member missing",
		)
	}

	if cmp.Member.Value == nil {
		return time.Time{}, errors.New(
			"invalid timestamp comparable -> value missing",
		)
	}

	if !cmp.Member.Value.Quoted {
		return time.Time{}, fmt.Errorf(
			"expected double-quoted timestamp but found unquoted value %q",
			cmp.Member.Value.Value,
		)
	}

	rawVal := cmp.Member.Value.Value
	if rawVal == "" {
		return time.Time{}, errors.New("empty timestamp literal")
	}

	t, err := time.Parse(time.RFC3339, rawVal)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"invalid timestamp format %q -> %w",
			rawVal,
			err,
		)
	}

	return t.Truncate(time.Microsecond), nil
}

// CoerceArgToDurationConstant attempts to extract a time.Duration constant
// from an Arg. It supports either a double-quoted duration string (e.g.,
// "10s", "1m") or an unquoted float/int representing seconds.
func CoerceArgToDurationConstant(arg *aip160.Arg) (time.Duration, error) {
	cmp, err := CoerceArgToComparable(arg)
	if err != nil {
		return 0, err
	}
	if cmp.Member == nil || cmp.Member.Value == nil {
		return 0, errors.New("invalid duration comparable")
	}

	if cmp.Member.Value.Quoted {
		rawVal := cmp.Member.Value.Value
		d, err := time.ParseDuration(rawVal)
		if err != nil {
			return 0, fmt.Errorf(
				"failed to parse duration string %q -> %w",
				rawVal,
				err,
			)
		}
		return d, nil
	}

	// Try integer seconds first, then float seconds
	rawVal := cmp.Member.Input()
	valInt, err := strconv.ParseInt(rawVal, 10, 64)
	if err == nil {
		return time.Duration(valInt) * time.Second, nil
	}
	valFloat, err := strconv.ParseFloat(rawVal, 64)
	if err == nil {
		return time.Duration(valFloat * float64(time.Second)), nil
	}

	return 0, fmt.Errorf("invalid duration format %q", rawVal)
}

// CoerceArgToNumericConstant attempts to extract an exact decimal constant
// (pgtype.Numeric) from an Arg.
func CoerceArgToNumericConstant(arg *aip160.Arg) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	cmp, err := CoerceArgToComparable(arg)
	if err != nil {
		return n, fmt.Errorf("failed to coerce arg -> %w", err)
	}

	if cmp.Member == nil {
		return n, errors.New("invalid decimal comparable -> member missing")
	}

	if cmp.Member.Value == nil {
		return n, errors.New("invalid decimal comparable -> value missing")
	}

	if cmp.Member.Value.Quoted {
		return n, fmt.Errorf(
			"expected an unquoted numeric literal but found string %q",
			cmp.Member.Value.Value,
		)
	}

	rawVal := cmp.Member.Input()
	if rawVal == "" {
		return n, errors.New("empty decimal literal")
	}

	err = n.Scan(rawVal)
	if err != nil {
		return n, fmt.Errorf(
			"failed to parse decimal literal %q -> %w",
			rawVal,
			err,
		)
	}

	return n, nil
}
