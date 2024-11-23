package repository

import (
	"context"

	"emperror.dev/errors"
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
	return toDomainComic(result)
}

func toDomainComic[T comic.GetAllComicTrackingsRow | *comic.GetComicTrackingByIDRow | *comic.ComicTracking](c T) (*domain.Comic, error) {
	extractFields := func(id int32, url, name, description, html, cronSpec string) *domain.Comic {
		return &domain.Comic{
			ID:          id,
			Url:         url,
			Name:        name,
			Description: description,
			Html:        html,
			CronSpec:    cronSpec,
		}
	}
	switch v := any(c).(type) {
	case *comic.GetComicTrackingByIDRow:
		if v == nil {
			return nil, nil
		}
		return extractFields(v.ID, v.Url, v.Name.String, v.Description.String, v.Html.String, v.CronSpec.String), nil
	case comic.GetAllComicTrackingsRow:
		return extractFields(v.ID, v.Url, v.Name.String, v.Description.String, v.Html.String, v.CronSpec.String), nil
	case *comic.ComicTracking:
		if v == nil {
			return nil, nil
		}
		return extractFields(v.ID, v.Url, v.Name.String, v.Description.String, v.Html.String, v.CronSpec.String), nil
	default:
		return nil, errors.New("not support type")
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
		comic, err := toDomainComic(r)
		if err != nil {
			return nil, err
		}
		comics = append(comics, *comic)
	}
	return comics, nil
}

func (u *userRepository) Save(ctx context.Context, req *domain.Comic) (*domain.Comic, error) {
	result, err := IgnoreNotFoundError2Params(u.query.CreateComicTracking(ctx, comic.CreateComicTrackingParams{
		Url: req.Url,
		Name: pgtype.Text{
			String: req.Name,
			Valid:  true,
		},
		Description: pgtype.Text{
			String: req.Description,
			Valid:  true,
		},
		Html: pgtype.Text{
			String: req.Html,
			Valid:  true,
		},
		CronSpec: pgtype.Text{
			String: req.CronSpec,
			Valid:  true,
		},
	}))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to save comic")
	}
	return toDomainComic(result)
}

func NewComicRepository(query *comic.Queries) domain.ComicRepository {
	return &userRepository{
		query: query,
	}
}
