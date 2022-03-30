package main

import (
	"fmt"
	"math/big"
	"os"

	"html/template"

	"crypto/rand"
	"regexp"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	//"github.com/mailru/easyjson"
	"textpad.com/db"
)

type PasteBody struct {
	ID        string
	Text      string
	CSRFToken string
}

const secret string = "_`53=Aj#3tvUg`x.^2s`kk?M:un37MW7&v>Hv#*{T(=DAyEXA<C@PMQ&i*m~V&:+&`"
const letters = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Validate(urlPath string) bool {
	ok, err := regexp.MatchString("^[a-zA-Z0-9]{8}$", urlPath)
	if err != nil {
		return false
	} else {
		return ok
	}
}
func cryptoRandAndSecure(max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		fmt.Println("Can't convert")
	}
	return nBig.Int64()
}
func GenerateUID() []byte {
	buf := make([]byte, 8)
	for i := 0; i < len(buf); i++ {
		nBig := cryptoRandAndSecure(int64(len(letters)))
		buf[i] = letters[nBig]
	}
	return buf
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
	uid := GenerateUID()
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
	ok := Validate(key)
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
		uid := GenerateUID()
		uidString := db.ConvertString(uid)
		cookie := fasthttp.Cookie{}
		cookie.SetKey("csrftoken")
		cookie.SetValue(uidString)
		ctx.Response.Header.SetCookie(&cookie)
		tpl := template.Must(template.ParseFiles("public/templates/paste.html"))
		paste := &PasteBody{ID: key, Text: db.ConvertString(ans),CSRFToken: uidString}
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
	key := GenerateUID()
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
	fmt.Println(csrftokenInputField==csrftoken)
	if text == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("400 Bad request"))
		return nil
	}
	
	if csrftokenInputField != csrftoken  {
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


func main() {
	router := routing.New()
	static := router.Group("/static")

	static.Get("/css/<path>", handleStaticCSS)
	static.Get("/js/<path>", handleStaticJS)
	static.Get("/images/<path>", handleStaticImages)
	router.Get("/", handleIndex)
	router.Post("/", handlePastePost)
	router.Get("/raw/<id>", apiAsyncHandlePasteGet)
	router.Get("/<id>", handlePaste)
	router.Post("/<id>", handleEditPaste)
	fmt.Println("Start server...")
	fmt.Println("Listen tcp://localhost:8080")
	fasthttp.ListenAndServe("localhost:8080", router.HandleRequest)

}
