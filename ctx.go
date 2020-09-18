package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"runtime"
)

type Ctx struct {
	w     http.ResponseWriter
	r     *http.Request
	users Users
}

func (ctx Ctx) Abort() {
	panic(AbortPanic{})
}

func (ctx Ctx) Return(content string, statusCode int) {
	ctx.w.WriteHeader(statusCode)
	ctx.w.Write([]byte(content))
	ctx.Abort()
}

func (ctx Ctx) ReturnBadRequest(message string) {
	ctx.Return(message, 400)
}

func (ctx Ctx) ReturnJSON(v interface{}, statusCode int) {
	jsonString, err := json.Marshal(v)
	ctx.CatchError(err)
	ctx.Return(string(jsonString), statusCode)
}

func (ctx Ctx) ReturnInternalError(err error) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("%s:%d %s: %v\n", file, line, ctx.r.URL, err)
	ctx.Return(err.Error(), 500)
}

func (ctx Ctx) CatchError(err error) {
	if err != nil {
		ctx.ReturnInternalError(err)
	}
}

func (ctx Ctx) ParseUTCOffset(key string) int {

	min := func(x, y int) int {
		if x < y {
			return x
		}
		return y
	}

	max := func(x, y int) int {
		if x > y {
			return x
		}
		return y
	}

	utcOffset, err := strconv.Atoi(ctx.r.FormValue(key))
	if err != nil {
		utcOffset = 0
	}
	return max(min(utcOffset, 14), -12)
}

func (ctx Ctx) SetSessionUser(userId string) {
	session, _ := store.Get(ctx.r, "swa")
	session.Values["user"] = userId
	session.Save(ctx.r, ctx.w)
}

func (ctx Ctx) ForceUserId() string {
	session, _ := store.Get(ctx.r, "swa")
	userId, ok := session.Values["user"].(string)
	if !ok {
		ctx.Return("Forbidden", 403)
	}
	return userId
}

func (ctx Ctx) ReturnUserData(userId string) {
	user := ctx.users.New(userId)
	defer user.Close()

	userData, err := user.GetData(ctx.ParseUTCOffset("utcoffset"))
	ctx.CatchError(err)
	ctx.ReturnJSON(userData, 200)
}

//func (ctx Ctx) Authenticate() {
//	ctx.User = ctx.users.New(ctx.ForceUserId())
//
//}