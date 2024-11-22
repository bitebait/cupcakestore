package messages

import (
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
)

type Message struct {
	Error   string
	Success string
}

const (
	ErrorMessageKey   = "error_message"
	SuccessMessageKey = "success_message"
)

func SetErrorMessage(ctx *fiber.Ctx, message string) {
	setSessionMessage(ctx, ErrorMessageKey, message)
}

func SetSuccessMessage(ctx *fiber.Ctx, message string) {
	setSessionMessage(ctx, SuccessMessageKey, message)
}

func setSessionMessage(ctx *fiber.Ctx, key, message string) {
	sess, _ := session.Store.Get(ctx)
	sess.Set(key, message)
	sess.Save()
}

func LoadMessages(ctx *fiber.Ctx) Message {
	msg := Message{}
	msg.Error = clearSessionMessage(ctx, ErrorMessageKey)
	msg.Success = clearSessionMessage(ctx, SuccessMessageKey)
	return msg
}

func clearSessionMessage(ctx *fiber.Ctx, key string) string {
	sess, err := session.Store.Get(ctx)
	if err != nil {
		return ""
	}
	message := ""
	if msg := sess.Get(key); msg != nil {
		message = msg.(string)
		sess.Delete(key)
	}
	sess.Save()
	return message
}
