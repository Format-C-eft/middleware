package api

import uuid "github.com/satori/go.uuid"

// UUIDIsValid - check uuid for correctness
func UUIDIsValid(text string) bool {

	_, err := uuid.FromString(text)

	return err == nil

}
