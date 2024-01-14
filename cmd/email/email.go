package main

import (
	"github.com/codewithtoucans/goweb/models"
)

func main() {
	emailService := models.NewEmailService()
	err := emailService.ForgotPassword("gzy52013@qq.com", "https://www.bing.com")
	if err != nil {
		println(err.Error())
		println("send email failed")
	}
}
