package logic

import (
	"strings"
	"techsupport/internal/models"
	"techsupport/internal/scoring"
)

func CalculateScoreForFirstEmail(userFirstEmail string, dbFirstEmail string, weights models.Weights) float64 {

	if userFirstEmail == "" || dbFirstEmail == "" {
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

    uPhone := CleanPhoneNumber(userFirstPhone)
    dPhone := CleanPhoneNumber(dbFirstPhone)

    if uPhone == dPhone {
        return weights.Phone * scoring.IdealMatch
    }

    const minLen = 7
    if len(uPhone) >= minLen && len(dPhone) >= minLen {

        checkLen := min(len(dPhone), len(uPhone))
        if checkLen > 10 {
            checkLen = 10
        }

        if uPhone[len(uPhone)-checkLen:] == dPhone[len(dPhone)-checkLen:] {
            return weights.Phone * scoring.PartialMatch
        }
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