package util

const (
	USD = "USD"
	INR = "INR"
	AED = "AED"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, INR, AED:
		return true
	}
	return false
}
