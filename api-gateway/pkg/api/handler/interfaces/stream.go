package interfaces

import "github.com/labstack/echo/v4"

type StreamHandler interface {
	Upload(ctx echo.Context) error
}
