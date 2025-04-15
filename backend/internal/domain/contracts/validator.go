package contracts

import "solution/internal/domain/dto"

type ValidatorService interface {
	ValidateRequestData(dto any) *dto.HttpErr
	ValidateRequestSlice(dtos any) *dto.HttpErr
}
