// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUser(ctx context.Context, username string) (User, error)
	GetUserById(ctx context.Context, id int32) (User, error)
}

var _ Querier = (*Queries)(nil)
