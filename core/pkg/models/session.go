package models

import "time"

// Session represents a continuous authenticated period of user activity.
// In a graph, comparing Session nodes against Transaction nodes helps 
// identify proxy usage or account takeovers (ATO).
type Session struct {
    // SessionID is the unique token identifying this specific authenticated period.
    SessionID string `json:"session_id" cypher:"id"`

    // SessionIP is the network address captured when the session was established.
    SessionIP string `json:"session_ip" cypher:"ip"`

    // DeviceID is the hardware identifier. Used to create a relationship with a Device node.
    DeviceID string `json:"device_id" cypher:"device_id"`

    // ASN is the Autonomous System Number (ISP) associated with the session's origin.
    ASN string `json:"asn" cypher:"asn"`

    // Country is the geographic origin of the session establishment.
    Country string `json:"country" cypher:"country"`

    // City is the specific city where the session originated.
    City string `json:"city" cypher:"city"`

    // StartTime is the RFC3339 timestamp of the login or session creation.
    StartTime time.Time `json:"start_time" cypher:"start_time"`

    // EndTime is the RFC3339 timestamp of logout or session expiration.
    EndTime time.Time `json:"end_time" cypher:"end_time"`
}