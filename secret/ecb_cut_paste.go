package main

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func getFlag() string {
	dat, err := ioutil.ReadFile("flag1.txt")
	if err != nil {
		return "Can't read flag"
	}
	return string(dat)
}

var Key []byte = make([]byte, 16)

func profile_for(email string) string {
	filter := func(r rune) rune {
		switch {
		case r == '=' || r == '&':
			return -1
		}
		return r
	}
	sanitizedemail := strings.Map(filter, email)
	return "email=" + sanitizedemail + "&uid=10&role=user"
}

func EncryptProfile(email string) []byte {
	return ECBEncrypt(Key, PKCS7([]byte(profile_for(email)), 16))
}

func DecryptProfile(encprofile []byte) string {
	decryptedprofile, err := PKCS7UnPad(ECBDecrypt(Key, encprofile))
	if err != nil {
		panic(err)
	}
	return string(decryptedprofile)
}

func main() {
	k, err := GenerateRandomBytes(16)
	if err != nil {
		panic(err)
	}
	Key = k

	http.HandleFunc("/encrypt", func(w http.ResponseWriter, r *http.Request) {
		emails, ok := r.URL.Query()["email"]
		if !ok || len(emails[0]) < 1 {
			fmt.Fprintf(w, "URL Parameter 'email' is missing\n")
			return
		}
		email, err := b64.StdEncoding.DecodeString(emails[0])
		if err != nil {
			fmt.Fprintf(w, "Something went wrong")
		}
		expire := time.Now().AddDate(0, 0, 1)
		cookie := http.Cookie{Name: "profile", Value: b64.StdEncoding.EncodeToString(EncryptProfile(string(email))), Path: "/", Expires: expire, MaxAge: 86400}
		http.SetCookie(w, &cookie)

	})

	http.HandleFunc("/decrypt", func(w http.ResponseWriter, r *http.Request) {
		var cookie, err = r.Cookie("profile")
		if err == nil {
			var cookievalue = cookie.Value
			encryptedProfile, err := b64.StdEncoding.DecodeString(cookievalue)
			if err != nil {
				fmt.Fprintf(w, "Bad cookie")
			}
			profile := DecryptProfile(encryptedProfile)
			w.Header().Set("Content-Type", "text/plain; charset=us-ascii")
			fmt.Fprintf(w, "Your cookie decodes to: %s\n", profile)
			m, err := url.ParseQuery(profile)
			if err != nil {
				log.Fatal(err)
			}
			if m["role"][0] == "admin" {
				fmt.Fprintf(w, getFlag())
			}
		}
	})
	http.ListenAndServe(":8088", nil)
}
