package logic

import (
	"techsupport/core/internal/models"
	"techsupport/core/internal/constants"
)

func GetWeights(isdonator bool) models.Weights{
	if isdonator {

		return models.Weights{
			RegDate:          constants.Solvent_Weight_RegDate,
			RegCountry:       constants.Solvent_Weight_RegCountry,
			RegCity:          constants.Solvent_Weight_RegCity,
			FirstEmail:       constants.Solvent_Weight_FirstEmail,
			Phone:            constants.Solvent_Weight_Phone,
			FirstDevice:      constants.Solvent_Weight_FirstDevice,
			Devices:          constants.Solvent_Weight_Devices,
			FirstTransaction: constants.Solvent_Weight_FirstTransaction,
		}
	}

	return models.Weights{
		RegDate:     constants.Insolvent_Weight_RegDate,
		RegCountry:  constants.Insolvent_Weight_RegCountry,
		RegCity:     constants.Insolvent_Weight_RegCity,
		FirstEmail:  constants.Insolvent_Weight_FirstEmail,
		Phone:       constants.Insolvent_Weight_Phone,
		FirstDevice: constants.Insolvent_Weight_FirstDevice,
		Devices:     constants.Insolvent_Weight_Devices,
		FirstTransaction: constants.Insolvent_Weight_FirstTransaction,
	}

}