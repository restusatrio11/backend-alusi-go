package postgres

import (
	"context"
	"errors"
	"fmt"

	"backend-alusi-go/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepo struct {
	pool *pgxpool.Pool
}

func NewCategoryRepo(pool *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{pool: pool}
}

func (r *CategoryRepo) List(ctx context.Context) ([]domain.Category, error) {
	query := `
	SELECT 
		c.id, c.nama, c.slug, c.deskripsi, c.urutan, c.ikon, c.created_at, c.updated_at,
		COUNT(a.id) FILTER (WHERE a.aktif = true) as app_count
	FROM categories c
	LEFT JOIN apps a ON c.id = a.category_id
	GROUP BY c.id
	ORDER BY c.urutan ASC, c.nama ASC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		err := rows.Scan(
			&c.ID, &c.Nama, &c.Slug, &c.Deskripsi, &c.Urutan, &c.Ikon, &c.CreatedAt, &c.UpdatedAt,
			&c.AppCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int) (*domain.Category, error) {
	query := `
	SELECT 
		c.id, c.nama, c.slug, c.deskripsi, c.urutan, c.ikon, c.created_at, c.updated_at,
		COUNT(a.id) FILTER (WHERE a.aktif = true) as app_count
	FROM categories c
	LEFT JOIN apps a ON c.id = a.category_id
	WHERE c.id = $1
	GROUP BY c.id
	`
	var c domain.Category
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Nama, &c.Slug, &c.Deskripsi, &c.Urutan, &c.Ikon, &c.CreatedAt, &c.UpdatedAt,
		&c.AppCount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category by id: %w", err)
	}
	return &c, nil
}

func (r *CategoryRepo) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	query := `
	SELECT 
		c.id, c.nama, c.slug, c.deskripsi, c.urutan, c.ikon, c.created_at, c.updated_at,
		COUNT(a.id) FILTER (WHERE a.aktif = true) as app_count
	FROM categories c
	LEFT JOIN apps a ON c.id = a.category_id
	WHERE c.slug = $1
	GROUP BY c.id
	`
	var c domain.Category
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&c.ID, &c.Nama, &c.Slug, &c.Deskripsi, &c.Urutan, &c.Ikon, &c.CreatedAt, &c.UpdatedAt,
		&c.AppCount,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}
	return &c, nil
}

func (r *CategoryRepo) Create(ctx context.Context, category *domain.Category) error {
	query := `
	INSERT INTO categories (nama, slug, deskripsi, urutan, ikon)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(
		ctx, query,
		category.Nama, category.Slug, category.Deskripsi, category.Urutan, category.Ikon,
	).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
}

func (r *CategoryRepo) Update(ctx context.Context, category *domain.Category) error {
	query := `
	UPDATE categories
	SET nama = $1, slug = $2, deskripsi = $3, urutan = $4, ikon = $5, updated_at = NOW()
	WHERE id = $6
	RETURNING updated_at
	`
	return r.pool.QueryRow(
		ctx, query,
		category.Nama, category.Slug, category.Deskripsi, category.Urutan, category.Ikon, category.ID,
	).Scan(&category.UpdatedAt)
}

func (r *CategoryRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
