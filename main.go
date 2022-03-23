package main

import (
	"encoding/json"
	"fmt"

	"html/template"

	"github.com/google/uuid"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/mailru/easyjson"
	"textpad.com/db"
)



type PasteBody struct {
	ID string
	Text string
}
func handleIndex(ctx *routing.Context) error {
	ctx.SetContentType("text/html")
	tpl := template.Must(template.ParseFiles("public/templates/index.html"))

	err := tpl.Execute(ctx, "index.html")
	if err != nil {
		ctx.Write([]byte("not found"))
		return nil
	} else {
		return nil

	}
}

func Validate(urlPath string) (uuid.UUID, error) {
	uid, err := uuid.Parse(urlPath)
	if err != nil {
		return uid, err
	}
	return uid, nil
}
func handlePaste(ctx *routing.Context) error {
	uid, err := Validate(ctx.Param("id"))
	if err != nil {
		ctx.SetContentType("text/plain")
		return err
	} else {
		ctx.SetContentType("text/plain")
		ctx.SetStatusCode(fasthttp.StatusOK)
		resp := fmt.Sprintf("Handle past %v", uid)
		ctx.SetBody([]byte(resp))
		return nil
	}

}
func apiHandlePastGet(ctx *routing.Context) error {
	ctx.SetContentType("text/plain")
	key := ctx.Param("id")
	db1 := db.InitDB("db/bolt.db")
	defer db1.Close()
	ans, err := db.GetDB(db1,"Paste",key)
	if err != nil {
		ctx.Write([]byte("Not found"))
	} else {
		ctx.Write([]byte(db.ConvertString(ans)))
	}

	return nil
}
func apiHandlePastPost(ctx *routing.Context) error {
	body := ctx.Request.Body()
	paste := &PasteBody{}
	return nil
}

func main() {
	router := routing.New()
	api := router.Group("/api")
	api.Post("/paste",apiHandlePastPost)
	api.Get("/paste/<id>",apiHandlePastGet)
	router.Get("/", handleIndex)
	// router.Get("/static/<path>",handleStatic)
	router.Get("/<id>", handlePaste)

	fmt.Println("Start server...")
	fmt.Println("Listen tcp://localhost:8080")
	fasthttp.ListenAndServe("localhost:8080", router.HandleRequest)
}
