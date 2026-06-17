package response

import "github.com/gofiber/fiber/v3"

type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
}

func Success(data interface{}) fiber.Map {
	return fiber.Map{"success": true, "data": data}
}

func Created(data interface{}) fiber.Map {
	return fiber.Map{"success": true, "data": data}
}

func Error(status int, message string) error {
	return fiber.NewError(status, message)
}

func Paginated(data interface{}, page, pageSize int, total int64) fiber.Map {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return fiber.Map{
		"success": true,
		"data":    data,
		"meta": Meta{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}
}
