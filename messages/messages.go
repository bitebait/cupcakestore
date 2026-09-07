package messages

import (
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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
	current, _ := ctx.Locals("Messages").(Message)
	if key == ErrorMessageKey {
		current.Error = message
	} else {
		current.Success = message
	}
	ctx.Locals("Messages", current)
	sess, err := session.Store.Get(ctx)
	if err != nil {
		log.Error("falha ao carregar sessão de mensagens")
		return
	}
	sess.Set(key, message)
	if err := sess.Save(); err != nil {
		log.Error("falha ao salvar sessão de mensagens")
	}
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
		message, _ = msg.(string)
		sess.Delete(key)
	}
	if err := sess.Save(); err != nil {
		log.Error("falha ao salvar sessão de mensagens")
	}
	return message
}
