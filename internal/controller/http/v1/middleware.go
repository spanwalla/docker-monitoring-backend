package v1

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spanwalla/docker-monitoring-backend/internal/service"
	"net/http"
	"strings"
)

const pingerIdCtx = "pingerId"

// AuthMiddleware -.
type AuthMiddleware struct {
	pingerService service.Pinger
}

// TODO: Можно использовать встроенный middleware

// PingerIdentity -.
func (h *AuthMiddleware) PingerIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := bearerToken(c.Request())
		if !ok {
			log.Errorf("AuthMiddleware.PingerIdentity - bearerToken: %v", ErrInvalidAuthHeader)
			newErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
			return nil
		}

		pingerId, err := h.pingerService.ParseToken(token)
		if err != nil {
			log.Errorf("AuthMiddleware.PingerIdentity - ParseToken: %v", err)
			newErrorResponse(c, http.StatusUnauthorized, ErrCannotParseToken.Error())
			return err
		}

		c.Set(pingerIdCtx, pingerId)

		return next(c)
	}
}

// bearerToken -.
func bearerToken(req *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := req.Header.Get(echo.HeaderAuthorization)
	if len(header) == 0 {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
