package main

import (
	"fmt"
	"os"

	"html/template"

	realip "github.com/Ferluci/fast-realip"
	routing "github.com/qiangxue/fasthttp-routing"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"textpad.com/db"
	"textpad.com/utils"
)

type PasteBody struct {
	ID        string
	Text      string
	CSRFToken string
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// api handle

func apiAsyncHandlePasteGet(ctx *routing.Context) error {
	key := ctx.Param("id")
	db1 := db.InitDB("db/bolt.db")
	defer db1.Close()
	ans := <-db.AsyncGetDB(db1, "Paste", key)
	if ans == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write([]byte("Not found"))
		return nil
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("text/plain; charset=utf-8")
	ctx.Write([]byte(db.ConvertString(ans)))
	return nil
}

func handleIndex(ctx *routing.Context) error {

	ctx.SetContentType("text/html")
	tpl := template.Must(template.ParseFiles("public/templates/index.html"))
	uid := utils.GenerateUID()
	cookie := fasthttp.Cookie{}
	cookieStringValue := string(uid)
	cookie.SetKey("csrftoken")
	cookie.SetValue(cookieStringValue)
	ctx.Response.Header.SetCookie(&cookie)
	paste := &PasteBody{CSRFToken: cookieStringValue}
	err := tpl.Execute(ctx, paste)
	if err != nil {
		ctx.Write([]byte("not found"))
	}
	return nil

}

func handlePaste(ctx *routing.Context) error {
	key := ctx.Param("id")
	ok := utils.Validate(key)
	if !ok {
		ctx.SetContentType("text/plain")
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Write([]byte("Not valid key"))
		return nil
	}
	db1 := db.InitDB("db/bolt.db")
	defer db1.Close()
	ans := <-db.AsyncGetDB(db1, "Paste", key)

	if ans == nil {
		ctx.SetContentType("text/plain")
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write([]byte("404 Not found"))
	} else {
		ctx.SetContentType("text/html")
		uid := utils.GenerateUID()
		uidString := db.ConvertString(uid)
		cookie := fasthttp.Cookie{}
		cookie.SetKey("csrftoken")
		cookie.SetValue(uidString)
		ctx.Response.Header.SetCookie(&cookie)
		tpl := template.Must(template.ParseFiles("public/templates/paste.html"))
		paste := &PasteBody{ID: key, Text: db.ConvertString(ans), CSRFToken: uidString}
		err := tpl.Execute(ctx, paste)
		if err != nil {
			ctx.Write([]byte("not found"))
			return nil
		} else {
			return nil
		}
	}
	return nil
}

func handlePastePost(ctx *routing.Context) error {
	text := ctx.FormValue("text")
	csrftokenInputField := db.ConvertString(ctx.FormValue("csrftoken"))
	csrftoken := db.ConvertString(ctx.Request.Header.Cookie("csrftoken"))
	if csrftokenInputField != csrftoken {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		ctx.Write([]byte("403 CSRF forbidden"))
	}
	if text == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("400 Bad request"))
	}

	db1 := db.InitDB("db/bolt.db")
	defer db1.Close()
	key := utils.GenerateUID()
	err := <-db.AsyncUpdateDB(db1, "Paste", db.ConvertString(key), db.ConvertString(text))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(fmt.Sprintf("Something went error: %v", err)))
		return nil
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		pathRedirect := "/" + db.ConvertString(key)
		ctx.Redirect(pathRedirect, fasthttp.StatusFound)
	}
	return nil
}
func handleEditPaste(ctx *routing.Context) error {
	key := ctx.Param("id")
	text := ctx.FormValue("text")
	csrftokenInputField := db.ConvertString(ctx.FormValue("csrftoken"))
	csrftoken := db.ConvertString(ctx.Request.Header.Cookie("csrftoken"))
	fmt.Println(csrftokenInputField == csrftoken)
	if text == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("400 Bad request"))
		return nil
	}

	if csrftokenInputField != csrftoken {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		ctx.Write([]byte("403 CSRF forbidden"))
		return nil
	}
	db1 := db.InitDB("db/bolt.db")
	defer db1.Close()
	err := <-db.AsyncUpdateDB(db1, "Paste", key, db.ConvertString(text))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte("500 Status Internal Server Error"))
	} else {

		ctx.Redirect("/"+key, fasthttp.StatusOK)
	}
	return nil
}

func handleStaticCSS(ctx *routing.Context) error {
	path := ctx.Param("path")
	path = "public/static/css/" + path
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write([]byte("404 Not found"))
		} else {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.Write([]byte("500 Internal server error"))
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SendFile(path)
	}
	return nil
}
func handleStaticJS(ctx *routing.Context) error {
	path := ctx.Param("path")
	path = "public/static/js/" + path
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write([]byte("404 Not found"))
		} else {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.Write([]byte("500 Internal server error"))
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SendFile(path)
	}
	return nil
}
func handleStaticImages(ctx *routing.Context) error {
	path := ctx.Param("path")
	path = "public/static/images/" + path
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write([]byte("404 Not found"))
		} else {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.Write([]byte("500 Internal server error"))
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SendFile(path)
	}
	return nil
}
func logMiddleware(ctx *routing.Context) error {
	method := string(ctx.Method())
	ip := realip.FromRequest(ctx.RequestCtx)
	status := ctx.RequestCtx.Response.StatusCode()
	uri := string(ctx.Request.URI().RequestURI())
	scheme := string(ctx.Request.URI().Scheme())
	host := string(ctx.Request.Host())
	ua := string(ctx.UserAgent())
	log.WithFields(log.Fields{
		"IP":         ip,
		"Method":     fmt.Sprintf("%s %s %s", method, uri, scheme),
		"Status":     status,
		"Host":       fmt.Sprintf("%s://%s%s", scheme, host, uri),
		"User-Agent": ua,
	}).Info()

	// fmt.Printf("%v -- %s %s://%s%s %s\n", method, t1,scheme, host, uri, ua)
	// 47.29.201.179 - - [28/Feb/2019:13:17:10 +0000] "GET /?p=1 HTTP/2.0" 200 5316 "https://domain1.com/?p=1" "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36" "2.75"
	return nil
}

func main() {
	router := routing.New()
	static := router.Group("/static")
	router.Use(logMiddleware)
	static.Get("/css/<path>", handleStaticCSS)
	static.Get("/js/<path>", handleStaticJS)
	static.Get("/images/<path>", handleStaticImages)
	router.Get("/", handleIndex)
	router.Post("/", handlePastePost)
	router.Get("/raw/<id>", apiAsyncHandlePasteGet)
	router.Get("/<id>", handlePaste)
	router.Post("/<id>", handleEditPaste)
	log.Println("Start server...")
	log.Println("Listen tcp://localhost:8080")
	fasthttp.ListenAndServe("localhost:8080", router.HandleRequest)

}
