package wrapper

import (
	"solution/internal/domain/dto"
)

func ServiceErr(httpCode int, message, description string) *dto.HttpErr {
	return &dto.HttpErr{HttpCode: httpCode, Message: message, Description: description}
}

func NotFoundErr(description string) *dto.HttpErr {
	return ServiceErr(404, "Not Found", description)
}

func InternalServerErr(description string) *dto.HttpErr {
	return ServiceErr(500, "Internal Server Error", description)
}

func ValidationErr(description string) *dto.HttpErr {
	return ServiceErr(400, "Request Data field value is invalid", description)
}

func BadRequestErr(description string) *dto.HttpErr {
	return ServiceErr(400, "Request Data structure is invalid", description)
}

func AccessForbiddenErr(description string) *dto.HttpErr {
	return ServiceErr(403, "Access Prohibited", description)
}

func ConflictErr(description string) *dto.HttpErr {
	return ServiceErr(409, "Conflict", description)
}

func OkResponse() *dto.OkResponse {
	return &dto.OkResponse{Message: "ok"}
}
