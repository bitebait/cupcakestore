package messages

import (
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type Message struct {
	Error   string
	Success string
}

const (
	ErrorMessageKey   = "error_message"
	SuccessMessageKey = "success_message"
)

func SetErrorMessage(ctx fiber.Ctx, message string) {
	setSessionMessage(ctx, ErrorMessageKey, message)
}

func SetSuccessMessage(ctx fiber.Ctx, message string) {
	setSessionMessage(ctx, SuccessMessageKey, message)
}

func setSessionMessage(ctx fiber.Ctx, key, message string) {
	current, _ := ctx.Locals("Messages").(Message)
	if key == ErrorMessageKey {
		current.Error = message
	} else {
		current.Success = message
	}
	ctx.Locals("Messages", current)
	sess := session.FromContext(ctx)
	if sess == nil {
		log.Error("middleware de sessão indisponível para mensagens")
		return
	}
	sess.Set(key, message)
}

func LoadMessages(ctx fiber.Ctx) Message {
	return Message{
		Error:   clearSessionMessage(ctx, ErrorMessageKey),
		Success: clearSessionMessage(ctx, SuccessMessageKey),
	}
}

func clearSessionMessage(ctx fiber.Ctx, key string) string {
	sess := session.FromContext(ctx)
	if sess == nil {
		return ""
	}
	message, _ := sess.Get(key).(string)
	sess.Delete(key)
	return message
}
