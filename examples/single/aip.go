package single

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/uthereal/scheme-runtime-go/pkg/contract"
	"github.com/uthereal/scheme-runtime-go/pkg/orm/aip160/fields"
)

// GlobalPostAip160Fields defines global search fields for Post models.
var GlobalPostAip160Fields = []contract.Aip160Field[contract.Where]{
	fields.NewStringField[
		string,
		contract.Where,
	](Schema.Public.Post.Title),
	fields.NewStringField[
		string,
		contract.Where,
	](Schema.Public.Post.Content),
}

// MapPostColumnToAip160Field maps Post column names to AIP-160 fields.
var MapPostColumnToAip160Field = map[string]contract.Aip160Field[
	contract.Where,
]{
	"id": fields.NewIntegerField[
		int64,
		contract.Where,
	](Schema.Public.Post.ID),
	"title": fields.NewStringField[
		string,
		contract.Where,
	](Schema.Public.Post.Title),
	"content": fields.NewStringField[
		string,
		contract.Where,
	](Schema.Public.Post.Content),
	"rating": fields.NewDecimalField[
		pgtype.Numeric,
		contract.Where,
	](Schema.Public.Post.Rating),
}

// GlobalUserAip160Fields defines global search fields for User models.
var GlobalUserAip160Fields = []contract.Aip160Field[contract.Where]{
	fields.NewStringField[
		string,
		contract.Where,
	](Schema.Public.User.Email),
}

// MapUserColumnToAip160Field maps User column names to AIP-160 fields.
var MapUserColumnToAip160Field = map[string]contract.Aip160Field[
	contract.Where,
]{
	"id": fields.NewIntegerField[
		int64,
		contract.Where,
	](Schema.Public.User.ID),
	"email": fields.NewStringField[
		string,
		contract.Where,
	](Schema.Public.User.Email),
	"age": fields.NewIntegerField[
		int,
		contract.Where,
	](Schema.Public.User.Age),
}
