// Package models contains the data structures used throughout the scoring engine.
package models

import "time"

// DBRecord represents the objective, historical state of a user account 
// retrieved from the primary database (the "Truth"). This data serves 
// as the verified baseline to validate claims and calculate risk scores.
type DBRecord struct {
    // AccTag is the unique internal identifier or handle for the account.
    AccTag string `json:"acc_tag" cypher:"acc_tag"`

    // RegCountry is the official ISO country code captured at the moment of registration.
    RegCountry string `json:"reg_country" cypher:"reg_country"`

    // RegCity is the verified city name associated with the account's creation.
    RegCity string `json:"reg_city" cypher:"reg_city"`

    // FirstEmail is the original primary email address used to register the account.
    FirstEmail string `json:"first_email" cypher:"first_email"`

    // Phone is the current verified phone number linked to the user's profile.
    Phone string `json:"phone" cypher:"phone"`

    // FirstDevice is the unique hardware identifier of the device used during registration.
    FirstDevice string `json:"first_device" cypher:"first_device"`

    // IsDonator indicates if the user has premium/supporter status, which may adjust scoring weights.
    IsDonator bool `json:"is_donator" cypher:"is_donator"`

    // RegDate is the RFC3339 formatted timestamp of when the account was officially created.
    RegDate time.Time `json:"reg_date" cypher:"reg_date"`
}