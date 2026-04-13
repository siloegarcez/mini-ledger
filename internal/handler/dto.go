package handler

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type (
	accountCreateRequest struct {
		Body AccountCreateRequest
	}
	AccountCreateRequest struct {
		DocumentNumber string `example:"12345678900" json:"document_number" maxLength:"15" minLength:"1"`
	}
	accountCreateResponse struct {
		Body AccountCreateResponse
	}
	AccountCreateResponse struct {
		AccountID      int64  `json:"account_id"     minimum:"1"`
		DocumentNumber string `example:"12345678900" json:"document_number" maxLength:"15"`
	}
)

type (
	accountGetByIDRequest struct {
		AccountID int64 `minimum:"1" path:"account_id"`
	}
	accountGetByIDResponse struct {
		Body AccountGetByIDResponse
	}
	AccountGetByIDResponse struct {
		AccountID      int64  `json:"account_id"     minimum:"1"`
		DocumentNumber string `example:"12345678900" json:"document_number" maxLength:"15" minLength:"1"`
	}
)

type (
	transactionCreateRequest struct {
		Body TransactionCreateRequest
	}
	TransactionCreateRequest struct {
		AccountID       int64  `json:"account_id"        minimum:"1"`
		OperationTypeID int64  `json:"operation_type_id" minimum:"1"`
		Amount          Number `example:"123.45"         json:"amount" minimum:"1"`
	}
	transactionCreateResponse struct {
		Body TransactionCreateResponse
	}
	TransactionCreateResponse struct {
		TransactionID   int64  `json:"transaction_id"    minimum:"1"`
		AccountID       int64  `json:"account_id"        minimum:"1"`
		OperationTypeID int64  `json:"operation_type_id" minimum:"1"`
		Amount          Number `example:"123.45"         json:"amount" minimum:"1"`
	}
	Number    string
	NumberOut json.RawMessage
)

// This overrides the default schema generation for the Number type, which would otherwise be treated as a string. Instead, we want it to be treated as a float64 in the schema but still be able to unmarshal from a JSON string. This allows us to have the benefits of using json.Number for precise decimal handling while still generating the correct schema for documentation and validation purposes.
func (o *Number) Schema(r huma.Registry) *huma.Schema {
	return huma.SchemaFromType(r, reflect.TypeFor[float64]())
}
func (o *NumberOut) Schema(r huma.Registry) *huma.Schema {
	return huma.SchemaFromType(r, reflect.TypeFor[float64]())
}

func (n *Number) UnmarshalJSON(data []byte) error {
	var num json.Number
	err := json.Unmarshal(data, &num)
	if err != nil {
		fmt.Println(err)
		return err
	}
	*n = Number(num)
	return nil
}

func (n Number) MarshalJSON() ([]byte, error) {
	num := json.Number(n)
	return json.Marshal(num)
}
