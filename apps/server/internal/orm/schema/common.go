package schema

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math/big"
	"reflect"
)

type Money big.Int

func (j *Money) ToNum() *big.Int {
	return (*big.Int)(j)
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *Money) Scan(value interface{}) error {
	var result *big.Int

	switch v := value.(type) {
	case int64:
		result = big.NewInt(v)
	case []byte:
		var ok bool
		if result, ok = big.NewInt(0).SetString(string(v), 10); !ok {
			return errors.New("Invalid number value: " + string(v))
		}
		*j = (Money)(*result)
	case string:
		var ok bool
		if result, ok = big.NewInt(0).SetString(v, 10); !ok {
			return errors.New("Invalid number value: " + v)
		}
	default:
		return errors.New(fmt.Sprint("Failed to scan value to big.Int: ", value, " type ", reflect.TypeOf(value)))
	}

	*j = (Money)(*result)
	return nil
}

// Value return json value, implement driver.Valuer interface
func (j Money) Value() (driver.Value, error) {
	return j.ToNum().String(), nil
}

func NewMoney(value string) (Money, bool) {
	result, ok := big.NewInt(0).SetString(value, 10)
	return Money(*result), ok
}
