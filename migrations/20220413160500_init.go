package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

// Source: https://github.com/uptrace/bun/blob/78c14c8b92e1/example/migrate/migrations/20210505110026_foo.go
func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
