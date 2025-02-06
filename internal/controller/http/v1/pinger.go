package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/spanwalla/docker-monitoring-backend/internal/service"
	"net/http"
)

// pingerRoutes -.
type pingerRoutes struct {
	pingerService service.Pinger
}

// newPingerRoutes -.
func newPingerRoutes(g *echo.Group, pingerService service.Pinger) {
	r := &pingerRoutes{pingerService}

	g.POST("/login", r.login)
	g.POST("/register", r.register)
}

// loginInput -.
type loginInput struct {
	Name     string `json:"name" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

// @Summary Login
// @Description Login
// @Tags pinger
// @Accept json
// @Produce json
// @Param input body loginInput true "input"
// @Success 200 {object} v1.pingerRoutes.login.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/pingers/login [post]
func (r *pingerRoutes) login(c echo.Context) error {
	var input loginInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	token, expiresAt, err := r.pingerService.GenerateToken(c.Request().Context(), service.PingerGenerateTokenInput{
		Name:     input.Name,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrPingerNotFound) {
			newErrorResponse(c, http.StatusBadRequest, "invalid name or password")
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}

	return c.JSON(http.StatusOK, response{
		Token:     token,
		ExpiresAt: expiresAt.Unix(),
	})
}

// registerInput -.
type registerInput struct {
	Name     string `json:"name" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

// @Summary Register
// @Description Register pinger
// @Tags pinger
// @Accept json
// @Produce json
// @Param input body registerInput true "input"
// @Success 201 {object} v1.pingerRoutes.register.response
// @Failure 400 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/pingers/register [post]
func (r *pingerRoutes) register(c echo.Context) error {
	var input registerInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.pingerService.CreatePinger(c.Request().Context(), service.PingerCreateInput{
		Name:     input.Name,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrPingerAlreadyExists) {
			newErrorResponse(c, http.StatusConflict, "pinger already exists")
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Id int `json:"id"`
	}

	return c.JSON(http.StatusCreated, response{
		Id: id,
	})
}
