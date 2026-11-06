package crud

import (
	"backend/utils/queries"
	"context"
)

type Core[T any] struct {
	res  *Resource[T]
	repo *Repository[T]
}

func NewCore[T any](res *Resource[T]) *Core[T] {
	return &Core[T]{
		res:  res,
		repo: NewRepository(res),
	}
}

func (c *Core[T]) Search(ctx context.Context, sc queries.SearchCriteria) ([]T, error) {
	results, err := c.repo.Search(ctx, sc)
	if err != nil {
		return nil, err
	}
	for i := range results {
		if err := c.check(ctx, &results[i]); err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (c *Core[T]) Get(ctx context.Context, id any) (*T, error) {
	result, err := c.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := c.check(ctx, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Core[T]) check(ctx context.Context, x *T) error {
	if c.res.Check == nil {
		return nil
	}
	return c.res.Check(ctx, x)
}
