package models

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Image struct {
	GalleryID int
	Path      string
	Filename  string
}

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *pgxpool.Pool

	ImageDir string
}

func (g *GalleryService) Create(title string, userID int) (*Gallery, error) {
	gallery := Gallery{Title: title, UserID: userID}
	row := g.DB.QueryRow(context.Background(), `INSERT INTO galleries (title, user_id) values ($1, $2) RETURNING id`, title, userID)
	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}

func (g *GalleryService) GetGalleryByID(id int) (*Gallery, error) {
	gallery := Gallery{}
	row := g.DB.QueryRow(context.Background(), `SELECT id, user_id, title FROM galleries WHERE id = $1`, id)
	err := row.Scan(&gallery.ID, &gallery.UserID, &gallery.Title)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrGalleryNotFound
		}
		return nil, fmt.Errorf("get gallery by id: %w", err)
	}
	return &gallery, nil
}

// GetGalleriesByUserID returns all galleries for a given user
func (g *GalleryService) GetGalleriesByUserID(userID int) ([]Gallery, error) {
	galleries := []Gallery{}
	rows, err := g.DB.Query(context.Background(), `SELECT id, title FROM galleries WHERE user_id = $1`, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrGalleryNotFound
		}
		return nil, fmt.Errorf("get galleries by user id: %w", err)
	}
	for rows.Next() {
		gallery := Gallery{UserID: userID}
		err := rows.Scan(&gallery.ID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("get galleries by user id: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("get galleries by user id: %w", err)
	}
	return galleries, nil
}

func (g *GalleryService) UpdateGallery(gallery Gallery) error {
	_, err := g.DB.Exec(context.Background(), `UPDATE galleries SET title = $1 WHERE id = $2`, gallery.Title, gallery.ID)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}

func (g *GalleryService) DeleteGalleryByID(id int) error {
	_, err := g.DB.Exec(context.Background(), `DELETE FROM galleries WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete gallery by id: %w", err)
	}
	err = os.RemoveAll(g.galleryDir(id))
	if err != nil {
		return fmt.Errorf("delete gallery by id: %w", err)
	}
	return nil
}

func (g *GalleryService) Images(galleryID int) ([]Image, error) {
	globPattern := filepath.Join(g.galleryDir(galleryID), "*")
	matchFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("images: %w", err)
	}
	images := make([]Image, 0, len(matchFiles))
	for _, f := range matchFiles {
		if hasExtension(f, extensions()) {
			images = append(images, Image{
				GalleryID: galleryID,
				Path:      f,
				Filename:  filepath.Base(f),
			})
		}
	}
	return images, nil
}

func (g *GalleryService) Image(galleryID int, filename string) (Image, error) {
	return Image{
		GalleryID: galleryID,
		Path:      filepath.Join(g.galleryDir(galleryID), filename),
		Filename:  filename,
	}, nil
}

func (g *GalleryService) DeleteImage(galleryID int, filename string) error {
	image, err := g.Image(galleryID, filename)
	if err != nil {
		return fmt.Errorf("delete image: %w", err)
	}
	err = os.Remove(image.Path)
	if err != nil {
		return fmt.Errorf("delete image: %w", err)
	}
	return nil
}

func extensions() []string {
	return []string{"jpg", "jpeg", "png", "gif"}
}

func (g *GalleryService) CreateImage(galleryID int, filename string, contents io.Reader) error {
	galleryDir := g.galleryDir(galleryID)
	err := os.MkdirAll(galleryDir, 0755)
	if err != nil {
		return fmt.Errorf("create image: %w", err)
	}
	f, err := os.Create(filepath.Join(galleryDir, filename))
	if err != nil {
		return fmt.Errorf("create image: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, contents)
	if err != nil {
		return fmt.Errorf("create image: %w", err)
	}
	return nil
}

func hasExtension(file string, exts []string) bool {
	ext := filepath.Ext(file)
	for _, e := range exts {
		if ext == "."+e {
			return true
		}
	}
	return false
}

func (g *GalleryService) galleryDir(id int) string {
	imageDir := g.ImageDir
	if imageDir == "" {
		imageDir = "images"
	}
	return filepath.Join(imageDir, fmt.Sprintf("gallery-%d", id))
}
