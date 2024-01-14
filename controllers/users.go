package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/codewithtoucans/goweb/context"
	"github.com/codewithtoucans/goweb/errors"
	"github.com/codewithtoucans/goweb/models"
)

type Users struct {
	Template struct {
		New            Template
		SignIn         Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword  Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")
	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		log.Println(err)
		http.Error(w, "process reset password was error", http.StatusInternalServerError)
		return
	}
	err = u.UserService.UpdatePassword(user.ID, data.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "process reset password was error", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "process reset password was error", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Template.ResetPassword.Execute(w, r, data)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	exist, err := u.UserService.CheckUserExist(data.Email)
	if err != nil || !exist {
		log.Println(err)
		log.Println("user not exist")
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		// TODO: if the user was not exist
		log.Println(err)
		http.Error(w, "something went wrong in process forgot password", http.StatusInternalServerError)
		return
	}
	values := url.Values{}
	values.Set("token", pwReset.Token)
	resetURL := fmt.Sprintf("http://localhost:3000/reset-pw?%s", values.Encode())
	err = u.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		log.Println(err)
		http.Error(w, "something went wrong in process forgot password", http.StatusInternalServerError)
		return
	}
	u.Template.CheckYourEmail.Execute(w, r, data)
}

func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Template.ForgotPassword.Execute(w, r, data)
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Template.New.Execute(w, r, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Template.SignIn.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		log.Printf("something was error in create %s", err.Error())
		u.Template.New.Execute(w, r, data, err)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			u.Template.New.Execute(w, r, data, err)
			return
		}
		log.Printf("something was wrong in process sign in %s", err)
		u.Template.New.Execute(w, r, data, errors.New("something was wrong in process sign in"))
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	_, _ = fmt.Fprintf(w, "user email %v\n", user)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := ReadCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(cookie)
	if err != nil {
		log.Println("delete session was error")
		http.Error(w, "delete session was error", http.StatusInternalServerError)
		return
	}
	DeleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}
