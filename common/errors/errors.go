package commonerrors

import "errors"

// commonly shared errors across services
var (
	ErrNoItems     = errors.New("Items must have at least one item.")
	ErrNoItemFound = errors.New("Item with the provided ID was not found in the database.")
)
