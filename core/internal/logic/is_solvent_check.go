// Package logic provides calculation engines for scoring user data.
package logic

import (
	"techsupport/core/internal/constants"
	"techsupport/core/pkg/models"
)

// GetWeights returns a pre-defined set of scoring weights based on the user's donation status.
// It branches the scoring logic into two main categories: Solvent (donators) and Insolvent (non-donators).
// This allows the engine to apply stricter or more relaxed validation rules depending on user history.
func GetWeights(isDonator bool) models.Weights {
	if isDonator {
		// Return weights optimized for users with a history of financial transactions (Solvent).
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

	// Return weights optimized for users without transaction history (Insolvent).
	return models.Weights{
		RegDate:          constants.Insolvent_Weight_RegDate,
		RegCountry:       constants.Insolvent_Weight_RegCountry,
		RegCity:          constants.Insolvent_Weight_RegCity,
		FirstEmail:       constants.Insolvent_Weight_FirstEmail,
		Phone:            constants.Insolvent_Weight_Phone,
		FirstDevice:      constants.Insolvent_Weight_FirstDevice,
		Devices:          constants.Insolvent_Weight_Devices,
		FirstTransaction: constants.Insolvent_Weight_FirstTransaction,
	}
}