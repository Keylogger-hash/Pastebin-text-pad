package main

import (
	"encoding/json"
	"fmt"
	"math/big"

	"html/template"

	"crypto/rand"

	"regexp"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	//"github.com/mailru/easyjson"
	"textpad.com/db"
)

type PasteBody struct {
	ID   string
	Text string
}

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
	ctx.SetContentType("text/plain")
	key := ctx.Param("id")
	db1 := db.InitDB("db/bolt.db")
	defer db1.Close()
	ans := <-db.AsyncGetDB(db1, "Paste", key)
	if ans == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write([]byte("Not found"))
		return nil
	}
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.Write([]byte(db.ConvertString(ans)))
	return nil
}
func apiAsyncHandlePastePost(ctx *routing.Context) error{
	body := ctx.PostBody()
	paste := &PasteBody{}
	err := json.Unmarshal(body, paste)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Not valid json"))
	} else {
		db1 := db.InitDB("db/bolt.db")
		defer db1.Close()
		err := <-db.AsyncUpdateDB(db1, "Paste", paste.ID, paste.Text)
		if err != nil {
			ctx.SetContentType("text/plain")
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.Write([]byte("Internal server error"))
		} else {
			ctx.SetContentType("text/plain")
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.Write([]byte("success:true"))
		}
	}
	return nil
}

func handleIndex(ctx *routing.Context) error {

	ctx.SetContentType("text/html")
	tpl := template.Must(template.ParseFiles("public/templates/index.html"))

	err := tpl.Execute(ctx, "index.html")
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
		tpl := template.Must(template.ParseFiles("public/templates/paste.html"))
		paste := &PasteBody{ID:key,Text: db.ConvertString(ans)}
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
	db1 := db.InitDB("db/bolt.db")
	defer db1.Close()
	key := GenerateUID()
	err := <-db.AsyncUpdateDB(db1, "Paste", db.ConvertString(key), db.ConvertString(text))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(fmt.Sprintf("Something went error: %v", err)))
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		pathRedirect := "/" + db.ConvertString(key)
		ctx.Redirect(pathRedirect, fasthttp.StatusFound)

	}
	return nil
}
func main() {
	router := routing.New()
	api := router.Group("/api")
	api.Post("/paste", apiAsyncHandlePastePost)
	api.Get("/raw/<id>", apiAsyncHandlePasteGet)
	router.Get("/", handleIndex)
	router.Post("/", handlePastePost)
	router.Get("/raw/<id>",apiAsyncHandlePasteGet)
	router.Get("/<id>", handlePaste)

	fmt.Println("Start server...")
	fmt.Println("Listen tcp://localhost:8080")
	fasthttp.ListenAndServe("localhost:8080", router.HandleRequest)

}
