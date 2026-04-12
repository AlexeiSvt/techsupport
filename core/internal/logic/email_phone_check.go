package logic

import (
	"strings"
	"techsupport/core/internal/models"
	"techsupport/core/internal/constants"
)

func CalculateScoreForFirstEmail(userFirstEmail string, dbFirstEmail string, weights models.Weights) float64 {

	if userFirstEmail == "" || dbFirstEmail == ""   {
        return constants.NoMatch
    }
	
	if  userFirstEmail == "" && dbFirstEmail == "" {
		return constants.NoMatch
	}
	
	if strings.EqualFold(userFirstEmail,dbFirstEmail) {
		return weights.FirstEmail * constants.IdealMatch
	}
	return constants.NoMatch
}

func CalculateScoreForFirstPhone(userFirstPhone, dbFirstPhone string, weights models.Weights) float64 {

	if userFirstPhone == "" || dbFirstPhone == "" {
		return constants.NoMatch
	}

	if userFirstPhone == "" && dbFirstPhone == "" {
		return constants.NoMatch
	}
    
    uPhone := CleanPhoneNumber(userFirstPhone)
    dPhone := CleanPhoneNumber(dbFirstPhone)

    if uPhone == dPhone {
        return weights.Phone * constants.IdealMatch
    }

        if len(uPhone) < constants.MinLen || len(dPhone) < constants.MinLen {
            return constants.NoMatch
        }

        checkLen := min(min(len(dPhone), len(uPhone)), 10)

        if uPhone[len(uPhone)-checkLen:] == dPhone[len(dPhone)-checkLen:] {
            return weights.Phone * constants.PartialMatch
        }

    return constants.NoMatch
}

func CleanPhoneNumber(phone string) string {
    return strings.Map(func(r rune) rune {
        if r >= '0' && r <= '9' {
            return r
        }
        return -1
    }, phone)
}