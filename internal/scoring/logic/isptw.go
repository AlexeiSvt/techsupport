package logic

import (
	"techsupport/internal/models"
	"techsupport/internal/scoring"
)

func GetWeights(isdonator bool) models.Weights{
	if isdonator {

		return models.Weights{
			RegDate:          scoring.P2W_Weight_RegDate,
			RegCountry:       scoring.P2W_Weight_RegCountry,
			RegCity:          scoring.P2W_Weight_RegCity,
			FirstEmail:       scoring.P2W_Weight_FirstEmail,
			Phone:            scoring.P2W_Weight_Phone,
			FirstDevice:      scoring.P2W_Weight_FirstDevice,
			Devices:          scoring.P2W_Weight_Devices,
			FirstTransaction: scoring.P2W_Weight_FirstTransaction,
		}
	}

	return models.Weights{
		RegDate:     scoring.F2P_Weight_RegDate,
		RegCountry:  scoring.F2P_Weight_RegCountry,
		RegCity:     scoring.F2P_Weight_RegCity,
		FirstEmail:  scoring.F2P_Weight_FirstEmail,
		Phone:       scoring.F2P_Weight_Phone,
		FirstDevice: scoring.F2P_Weight_FirstDevice,
		Devices:     scoring.F2P_Weight_Devices,
		FirstTransaction: scoring.F2P_Weight_FirstTransaction,
	}

}