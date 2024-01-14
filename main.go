package main

import (
	"log"
	"net/http"

	"github.com/codewithtoucans/goweb/controllers"
	"github.com/codewithtoucans/goweb/middleware"
	"github.com/codewithtoucans/goweb/migrations"
	"github.com/codewithtoucans/goweb/models"
	"github.com/codewithtoucans/goweb/templates"
	"github.com/codewithtoucans/goweb/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

const _port = ":3000"

func main() {
	faqHandler := controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml")))
	homeHandler := controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml")))
	contactHandler := controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml")))
	notFoundHandler := controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "notfound.gohtml", "tailwind.gohtml")))

	signupTemplate := views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	signinTemplate := views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	forgotTemplate := views.Must(views.ParseFS(templates.FS, "forgot-pw.gohtml", "tailwind.gohtml"))
	checkYourEmailTemplate := views.Must(views.ParseFS(templates.FS, "check-your-email.gohtml", "tailwind.gohtml"))
	resetPasswordTemplate := views.Must(views.ParseFS(templates.FS, "reset-pw.gohtml", "tailwind.gohtml"))

	galleryNewTemplate := views.Must(views.ParseFS(templates.FS, "gallery-new.gohtml", "tailwind.gohtml"))
	galleryEditTemplate := views.Must(views.ParseFS(templates.FS, "gallery-edit.gohtml", "tailwind.gohtml"))
	galleryIndexTemplate := views.Must(views.ParseFS(templates.FS, "gallery-index.gohtml", "tailwind.gohtml"))
	galleryImageTemplate := views.Must(views.ParseFS(templates.FS, "gallery-show.gohtml", "tailwind.gohtml"))

	config := models.DefaultPostgresConfig()
	conn, err := models.Open(config)
	if err != nil {
		panic("connect database error")
	}
	defer conn.Close()
	if err != nil {
		panic("Unable to connect to database")
	}

	err = models.MigrateFS(migrations.FS, ".")
	if err != nil {
		log.Println(err)
		panic("migrations error")
	}

	userService := &models.UserService{DB: conn}
	sessionService := &models.SessionService{DB: conn}
	emailService := models.NewEmailService()
	passwordResetService := &models.PasswordResetService{DB: conn}
	galleryService := &models.GalleryService{DB: conn}

	userMiddleWare := middleware.UserMiddleWare{SessionService: sessionService}
	setUserMiddleWare := userMiddleWare.SetUser
	requireUserMiddleWare := userMiddleWare.RequireUser

	userController := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		EmailService:         emailService,
		PasswordResetService: passwordResetService,
	}
	userController.Template.New = signupTemplate
	userController.Template.SignIn = signinTemplate
	userController.Template.ForgotPassword = forgotTemplate
	userController.Template.CheckYourEmail = checkYourEmailTemplate
	userController.Template.ResetPassword = resetPasswordTemplate

	galleryController := controllers.GalleryController{
		GalleryService: galleryService,
	}
	galleryController.Template.New = galleryNewTemplate
	galleryController.Template.Edit = galleryEditTemplate
	galleryController.Template.Index = galleryIndexTemplate
	galleryController.Template.Show = galleryImageTemplate

	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/faq", faqHandler)
	r.Get("/contact", contactHandler)
	r.Get("/signup", userController.New)
	r.Get("/signin", userController.SignIn)
	r.Get("/forgot-pw", userController.ForgotPassword)
	r.Get("/reset-pw", userController.ResetPassword)
	r.Post("/users", userController.Create)
	r.Post("/signin", userController.ProcessSignIn)
	r.Post("/signout", userController.ProcessSignOut)
	r.Post("/reset-pw", userController.ProcessResetPassword)
	r.Post("/forgot-pw", userController.ProcessForgotPassword)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(requireUserMiddleWare)
		r.Get("/", userController.CurrentUser)
	})

	r.Route("/gallery", func(r chi.Router) {
		r.Get("/{id}", galleryController.Show)
		r.Get("/{id}/images/{filename}", galleryController.Image)
		r.Group(func(r chi.Router) {
			r.Use(requireUserMiddleWare)
			r.Get("/", galleryController.Index)
			r.Get("/new", galleryController.New)
			r.Post("/", galleryController.Create)
			r.Get("/{id}/edit", galleryController.Edit)
			r.Post("/{id}", galleryController.Update)
			r.Post("/{id}/delete", galleryController.Delete)
			r.Post("/{id}/images/{filename}/delete", galleryController.DeleteImage)
		})
	})
	r.NotFound(notFoundHandler)

	h := http.FileServer(http.Dir("asserts"))
	r.Get("/asserts/*", http.StripPrefix("/asserts/", h).ServeHTTP)

	CSRF := csrf.Protect([]byte("l4Tu4L5CPtXH8xD9g8tRS9tvknyCdoNf"), csrf.Secure(false))

	log.Printf("start server on %s\n", _port)
	err = http.ListenAndServe(_port, CSRF(setUserMiddleWare(r)))
	if err != nil {
		log.Fatalln("start server failed")
		return
	}
}
