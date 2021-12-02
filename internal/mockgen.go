package mocks

//go:generate mockgen -destination=./mocks/db_mock.go -package=mocks github.com/Format-C-eft/middleware/internal/database/cache Client
