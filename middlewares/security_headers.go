package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/helmet"
)

// Local scripts/styles need no inline exceptions. Image previews use blob URLs;
// Pix is the only external destination allowed after a form submission.
func SecurityHeaders() fiber.Handler {
	headers := helmet.New(helmet.Config{
		ContentSecurityPolicy: "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data: blob:; font-src 'self' data:; connect-src 'self'; object-src 'none'; base-uri 'none'; frame-ancestors 'none'; form-action 'self'",
		XFrameOptions:         "DENY",
		ReferrerPolicy:        "no-referrer",
		PermissionPolicy:      "camera=(), microphone=(), geolocation=()",
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: true,
	})
	return func(c fiber.Ctx) error {
		err := headers(c)
		contentType := c.GetRespHeader(fiber.HeaderContentType)
		if strings.HasPrefix(contentType, fiber.MIMETextHTML) || strings.HasPrefix(contentType, fiber.MIMEApplicationJSON) {
			// HTML contains session-bound CSRF tokens and JSON can contain account data.
			c.Set(fiber.HeaderCacheControl, "private, no-store")
		}
		return err
	}
}
