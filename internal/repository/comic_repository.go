package repository

import (
	"context"
	"errors"

	"github.com/cuongtranba/mynoti/internal/domain"
	"github.com/cuongtranba/mynoti/internal/repository/sqlc/comic"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type userRepository struct {
	query *comic.Queries
}

func (u *userRepository) Delete(ctx context.Context, id int32) error {
	return u.query.DeleteComicTracking(ctx, id)
}

func isNotFoundError(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func IgnoreNotFoundError2Params[T any](result T, err error) (*T, error) {
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func IgnoreNotFoundError(err error) error {
	if isNotFoundError(err) {
		return nil
	}
	return err
}

func (u *userRepository) Get(ctx context.Context, id int32) (*domain.Comic, error) {
	result, err := IgnoreNotFoundError2Params(u.query.GetComicTrackingByID(ctx, id))
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	domainComic := toDomainComic(*result)
	return &domainComic, nil
}

func toDomainComic(c comic.ComicTracking) domain.Comic {
	return domain.Comic{
		ID:          c.ID,
		Url:         c.Url,
		Name:        c.Name.String,
		Description: c.Description.String,
		Html:        c.Html.String,
	}
}

func (u *userRepository) List(ctx context.Context) ([]domain.Comic, error) {
	result, err := IgnoreNotFoundError2Params(u.query.GetAllComicTrackings(ctx))
	if err != nil {
		return nil, err
	}
	if result == nil || len(*result) == 0 {
		return nil, nil
	}
	var comics []domain.Comic
	for _, r := range *result {
		comics = append(comics, toDomainComic(r))
	}
	return comics, nil
}

func (u *userRepository) Save(ctx context.Context, req *domain.Comic) error {
	err := IgnoreNotFoundError(u.query.CreateComicTracking(ctx, comic.CreateComicTrackingParams{
		Url: req.Url,
		Name: pgtype.Text{
			String: req.Name,
		},
		Description: pgtype.Text{
			String: req.Description,
		},
		Html: pgtype.Text{
			String: req.Html,
		},
	}))
	return err
}

func NewComicRepository(query *comic.Queries) domain.ComicRepository {
	return &userRepository{
		query: query,
	}
}
