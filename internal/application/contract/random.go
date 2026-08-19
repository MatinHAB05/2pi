package service_contract

type RandomService interface {
	GeneratePositiveIntRandomNumberInRange(min, max int64) int64
	GeneratePositiveIntRandomNumber(digits int) int64
}
