package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/codewithtoucans/goweb/context"
	"github.com/codewithtoucans/goweb/models"
	"github.com/go-chi/chi/v5"
)

type GalleryController struct {
	Template struct {
		New   Template
		Edit  Template
		Index Template
		Show  Template
	}
	GalleryService *models.GalleryService
}

func (g GalleryController) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Template.New.Execute(w, r, data)
}

func (g GalleryController) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID int
		Title  string
	}
	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")
	log.Println(data)

	gallery, err := g.GalleryService.Create(data.Title, data.UserID)
	if err != nil {
		g.Template.New.Execute(w, r, data, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/gallery/%d/edit", gallery.ID), http.StatusFound)
}

func (g GalleryController) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "gallery edit error", http.StatusInternalServerError)
		return
	}
	gallery, err := g.GalleryService.GetGalleryByID(id)
	if err != nil {
		http.Error(w, "gallery edit error", http.StatusInternalServerError)
		return
	}
	if gallery.UserID != context.User(r.Context()).ID {
		http.Error(w, "gallery edit permission denied", http.StatusForbidden)
		return
	}
	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string
	}
	var data struct {
		ID     int
		Title  string
		Images []Image
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		http.Error(w, "gallery edit error", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       gallery.ID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}
	g.Template.Edit.Execute(w, r, data)
}

func (g GalleryController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "gallery edit error", http.StatusInternalServerError)
		return
	}
	gallery, err := g.GalleryService.GetGalleryByID(id)
	if err != nil {
		http.Error(w, "gallery edit error", http.StatusInternalServerError)
		return
	}
	if gallery.UserID != context.User(r.Context()).ID {
		http.Error(w, "gallery edit permission denied", http.StatusForbidden)
		return
	}
	gallery.Title = r.FormValue("title")
	err = g.GalleryService.UpdateGallery(*gallery)
	if err != nil {
		http.Error(w, "gallery edit error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/gallery/%d/edit", gallery.ID), http.StatusFound)
}

func (g GalleryController) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID    int
		Title string
	}
	var data struct {
		Galleries []Gallery
	}
	galleries, err := g.GalleryService.GetGalleriesByUserID(context.User(r.Context()).ID)
	if err != nil {
		http.Error(w, "gallery index error", http.StatusInternalServerError)
		return
	}
	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		})
	}
	g.Template.Index.Execute(w, r, data)
}

func (g GalleryController) Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "show error", http.StatusInternalServerError)
		return
	}
	gallery, err := g.GalleryService.GetGalleryByID(id)
	if err != nil {
		http.Error(w, "show error", http.StatusInternalServerError)
		return
	}
	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string
	}
	var data struct {
		ID     int
		Title  string
		Images []Image
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		http.Error(w, "show error", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       image.GalleryID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}
	// for i := 0; i < 20; i++ {
	// 	width, height := rand.Intn(500)+200, rand.Intn(500)+200
	// 	data.Images = append(data.Images, fmt.Sprintf("https://placekitten.com/%d/%d", width, height))
	// }
	g.Template.Show.Execute(w, r, data)
}

func (g GalleryController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "delete error", http.StatusInternalServerError)
		return
	}
	gallery, err := g.GalleryService.GetGalleryByID(id)
	if err != nil {
		http.Error(w, "delete error", http.StatusInternalServerError)
		return
	}
	if gallery.UserID != context.User(r.Context()).ID {
		http.Error(w, "delete permission denied", http.StatusForbidden)
		return
	}
	err = g.GalleryService.DeleteGalleryByID(gallery.ID)
	if err != nil {
		http.Error(w, "delete error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g GalleryController) Image(w http.ResponseWriter, r *http.Request) {
	filename := g.filename(w, r)
	galleryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "image error", http.StatusInternalServerError)
		return
	}
	image, err := g.GalleryService.Image(galleryID, filename)
	if err != nil {
		log.Println(err)
		http.Error(w, "image error", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, image.Path)
}

func (g GalleryController) DeleteImage(w http.ResponseWriter, r *http.Request) {
	filename := g.filename(w, r)
	galleryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "image error", http.StatusInternalServerError)
		return
	}
	gallery, err := g.GalleryService.GetGalleryByID(galleryID)
	if err != nil {
		http.Error(w, "delete error", http.StatusInternalServerError)
		return
	}
	if gallery.UserID != context.User(r.Context()).ID {
		http.Error(w, "delete permission denied", http.StatusForbidden)
		return
	}
	err = g.GalleryService.DeleteImage(galleryID, filename)
	if err != nil {
		http.Error(w, "delete error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/galleries/%d/edit", galleryID), http.StatusFound)
}

func (g GalleryController) UploadImage(w http.ResponseWriter, r *http.Request) {
	galleryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "upload image error", http.StatusInternalServerError)
		return
	}
	gallery, err := g.GalleryService.GetGalleryByID(galleryID)
	if err != nil {
		http.Error(w, "upload image error", http.StatusInternalServerError)
		return
	}
	if gallery.UserID != context.User(r.Context()).ID {
		http.Error(w, "upload image permission denied", http.StatusForbidden)
		return
	}
	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, "upload image error", http.StatusInternalServerError)
		return
	}
	fileHeaders := r.MultipartForm.File["images"]
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "upload image error", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		err = g.GalleryService.CreateImage(gallery.ID, fileHeader.Filename, file)
		if err != nil {
			http.Error(w, "upload image error", http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/galleries/%d/edit", galleryID), http.StatusFound)
}

func (g GalleryController) filename(w http.ResponseWriter, r *http.Request) string {
	filename := chi.URLParam(r, "filename")
	return filepath.Base(filename)
}
