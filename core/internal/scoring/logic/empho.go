package logic

import (
	"strings"
	"techsupport/core/internal/models"
	"techsupport/core/internal/scoring"
)

func CalculateScoreForFirstEmail(userFirstEmail string, dbFirstEmail string, weights models.Weights) float64 {

	if userFirstEmail == "" || dbFirstEmail == ""   {
        return scoring.NoMatch
    }
	
	if  userFirstEmail == "" && dbFirstEmail == "" {
		return scoring.NoMatch
	}
	
	if strings.EqualFold(userFirstEmail,dbFirstEmail) {
		return weights.FirstEmail * scoring.IdealMatch
	}
	return scoring.NoMatch
}

func CalculateScoreForFirstPhone(userFirstPhone, dbFirstPhone string, weights models.Weights) float64 {

	if userFirstPhone == "" || dbFirstPhone == "" {
		return scoring.NoMatch
	}

	if userFirstPhone == "" && dbFirstPhone == "" {
		return scoring.NoMatch
	}
    
    uPhone := CleanPhoneNumber(userFirstPhone)
    dPhone := CleanPhoneNumber(dbFirstPhone)

    if uPhone == dPhone {
        return weights.Phone * scoring.IdealMatch
    }

        if len(uPhone) < scoring.MinLen || len(dPhone) < scoring.MinLen {
            return scoring.NoMatch
        }

        checkLen := min(min(len(dPhone), len(uPhone)), 10)

        if uPhone[len(uPhone)-checkLen:] == dPhone[len(dPhone)-checkLen:] {
            return weights.Phone * scoring.PartialMatch
        }

    return scoring.NoMatch
}

func CleanPhoneNumber(phone string) string {
    return strings.Map(func(r rune) rune {
        if r >= '0' && r <= '9' {
            return r
        }
        return -1
    }, phone)
}