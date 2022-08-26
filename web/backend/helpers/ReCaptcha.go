package helpers

import (
	"gopkg.in/ezzarghili/recaptcha-go.v4"
	"log"
	"time"
)

func RecaptchaCheck(key string) bool { // Check recaptcha
	{
		captcha, _ := recaptcha.NewReCAPTCHA("6Lcxl8caAAAAAKNxpAnH4n96vtciU7pq_h20gwgK", recaptcha.V2, 10*time.Second) // for v2 API get your secret from https://www.google.com/recaptcha/admin
		err := captcha.Verify(key)
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
}
