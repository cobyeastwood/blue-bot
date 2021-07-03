package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

// Still in Development
func OAuth(mc *MainConfig, PORT string) {
	go StartTLSServer(mc, PORT)
	go StartOAuth(mc, PORT)
}

func StartOAuth(mc *MainConfig, PORT string) {

	p := url.Values{}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cookieJar, _ := cookiejar.New(nil)

	cli := &http.Client{Transport: tr,
		Jar: cookieJar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	p.Add("client_id", os.Get("YAHOO_KEY"))
	p.Add("redirect_uri", fmt.Sprintf("https://localhost:%v/oauth2/yahoo/receive", PORT))
	p.Add("response_type", "code")
	p.Add("state", "0000")

	params := strings.NewReader(p.Encode())

	req, _ := http.NewRequest("POST", mc.y.Endpoint.AuthURL, params)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	re, _ := cli.Do(req)

	bs, _ := ioutil.ReadAll(re.Body)

	re.Body.Close()
}

func StartTLSServer(mc *MainConfig, PORT string) {

	wg := new(sync.WaitGroup)
	wg.Add(2)

	cert, errCert := tls.LoadX509KeyPair("localhost.crt", "localhost.key")

	if errCert != nil {
		log.Fatal(errCert)
	}

	s := &http.Server{
		Addr:    ":" + fmt.Sprintf("%v", PORT),
		Handler: nil,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	s.RegisterOnShutdown(func() {
		wg.Done()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "test!")
	})

	http.HandleFunc("/oauth2/yahoo", func(w http.ResponseWriter, r *http.Request) {
		redirect := mc.y.AuthCodeURL("0000")
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	})

	http.HandleFunc("/oauth2/yahoo/receive", func(w http.ResponseWriter, r *http.Request) {

		code := r.FormValue("code")
		state := r.FormValue("state")

		if state != "0000" {
			http.Error(w, "state is incorrect", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		tkn, err := mc.y.Exchange(ctx, code)

		if err != nil {
			http.Error(w, "authorization is incorrect", http.StatusInternalServerError)
			return
		}

		tknSrc := mc.y.TokenSource(ctx, tkn)

		mc.yc = oauth2.NewClient(ctx, tknSrc)

		re, err := mc.yc.Get(url)

		if err != nil {
			log.Fatal(err)
		}

		defer re.Body.Close()

		bs, _ := ioutil.ReadAll(re.Body)

		var ok bool

		if mc.yd, ok = ctx.Deadline(); !ok {
			log.Fatal(err)
		}

	})

	now := time.Now()

	if now.After(mc.yd) {
		s.Shutdown(context.Background())
		wg.Done()
	}

	s.ListenAndServeTLS("", "")
	wg.Wait()

}
