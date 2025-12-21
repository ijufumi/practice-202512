package util

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type contextKey string

const (
	dbContextKey contextKey = "db"
)

// SetDB sets gorm.DB instance to context
func SetDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbContextKey, db)
}

// GetDB retrieves gorm.DB instance from context
func GetDB(ctx context.Context) (*gorm.DB, error) {
	db, ok := ctx.Value(dbContextKey).(*gorm.DB)
	if !ok || db == nil {
		return nil, errors.New("database connection not found in context")
	}
	return db, nil
}

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, "userID", userID)
}

func GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}
