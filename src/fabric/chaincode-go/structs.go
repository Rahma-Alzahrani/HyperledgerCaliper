package main

import (
	"time"
	//"github/gogo "
)

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	Record    *DataOffer `json:"record"`
	TxId      string     `json:"txId"`
	Timestamp time.Time  `json:"timestamp"`
	IsDelete  bool       `json:"isDelete"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *DataOffer
}

type HistoricalQuery struct {
	bool
	Key     string `json:"Key"`
	Records *HistoricalDataOffer
}

type JourneySchedule struct {
	DocType    string `json:"docType"`
	UID        string `json:"uid"`
	Type       string `json:"type"`
	Cat        string `json:"cat"`
	Journey    string `json:"journey"`
	Valid_from string `json:"valid_from"`
	Valid_to   string `json:"valid_to"`
	Days       string `json:"days"`
	Operator   string `json:"operator"`
}

type JourneyDataHash struct {
	DocType      string     `json:"docType"`
	ID           string     `json:"id"`
	DataHashes   []DataHash `json:"data_hashes"`
	JourneyUID   string     `json:"journey_uid"`
	DataProvider string     `json:"data_provider"`
}

type HistoricalDataAgreement struct {
	DocType           string   `json:"docType"`
	ID                string   `json:"id"`
	OfferID           string   `json:"offer_id"`
	DataProvider      string   `json:"data_provider"`
	DataConsumer      string   `json:"data_consumer"`
	EscrowID          string   `json:"escrow_id"`
	StartDate         string   `json:"startdate"`
	EndDate           string   `json:"enddate"`
	Price             float64  `json:"price"`
	OfferRequestID    string   `json:"offerrequest_id"`
	JourneyDataHashID []string `json:"journeydatahash_id"`
	State             bool     `json:"state"`
}

type HistoricalDataOffer struct {
	DocType         string   `json:"docType"`
	ID              string   `json:"id"`
	Validity        bool     `json:"validity"`
	DataOwner       string   `json:"data_owner"`
	Equipment       string   `json:"equipment"`
	MoniteredAsset  string   `json:"moniteredAsset"`
	ProcessingLevel string   `json:"processing_level"`
	Price           float64  `json:"price"`
	Deposit         float64  `json:"deposit"`
	OwnerOrg        string   `json:"owner_org"`
	Creator         string   `json:"creator"`
	Operator        string   `json:"operator"`
	JourneyUID      string   `json:"journey_uid"`
	StartDate       string   `json:"start_date"`
	EndDate         string   `json:"end_date"`
}

// DataOffer type
type DataOffer struct {
	DocType         string  `json:"docType"`
	ID              string  `json:"id"`
	Validity        bool    `json:"validity"`
	DataOwner       string  `json:"data_owner"`
	Equipment       string  `json:"equipment"`
	MoniteredAsset  string  `json:"monitered_asset"`
	ProcessingLevel string  `json:"processing_level"`
	Price           float64 `json:"price"`
	Deposit         float64 `json:"deposit"`
	OwnerOrg        string  `json:"owner_org"`
	Creator         string  `json:"creator"`
	Operator        string  `json:"operator"`
	Journey_UID     string  `json:"journey_uid"`
	Depart_time     string  `json:"depart_time"`
	Arrival_time    string  `json:"arrival_time"`
	IsActive		bool    `json:"is_active"`
}

// Escrow type
type Escrow struct {
	Consumer        string  `json:"consumer"`
	Provider        string  `json:"provider"`
	DocType         string  `json:"docType"`
	ID              string  `json:"id"`
	ProviderDeposit float64 `json:"providerDeposit"`
	ConsumerDeposit float64 `json:"consumerDeposit"`
	ConsumerPayment float64 `json:"consumerPayment"`
	Released        bool    `json:"released"`
	OfferRequestID  string  `json:"offer_request_id"`
	OfferID         string  `json:"offer_id"`
	Status          string  `json:"status"`
	StartDate       string  `json:"startDate"`
	EndDate         string  `json:"endDate"`
	AgreementID     string  `json:"agreement_id"`
}

// DataHash type
type DataHash struct {
	DocType   string `json:"docType"`
	ID        string `json:"id"`
	Hash      string `json:"hash"`
	FileName  string `json:"file_name"`
	EntryDate string `json:"entryDate"`
}

// DataAgreement type
type DataAgreement struct {
	Price           float64  `json:"price"`
	DocType         string   `json:"docType"`
	ID              string   `json:"id"`
	DataProvider    string   `json:"dataProvider"`
	DataConsumer    string   `json:"dataConsumer"`
	EscrowID        string   `json:"escrow_id"`
	State           bool     `json:"state"`
	OfferRequestID  string   `json:"offer_request_id"`
	OfferID         string   `json:"offer_id"`
	StartDate       string   `json:"start_date"`
	EndDate         string   `json:"end_date"`
	OfferDataHashID []string `json:"offer_data_hash_id"`
	ProviderDeposit float64  `json:"providerDeposit"`
	ConsumerDeposit float64  `json:"consumerDeposit"`
}

type OfferDataHash struct {
	ID           string     `json:"id"`
	DataHashes   []DataHash `json:"data_hashes"`
	OfferID      string     `json:"offer_id"`
	DocType      string     `json:"docType"`
	DataProvider string     `json:"data_provider"`
}

// OfferRequest type
type OfferRequest struct {
	DocType        string    `json:"docType"`
	OfferRequestID string    `json:"offer_request_id"`
	OfferID        string    `json:"offer_id"`
	DataConsumer   string    `json:"dataConsumer"`
	DataProvider   string    `json:"dataProvider"`
	Price          float64   `json:"price"`
	PDeposit       float64   `json:"pDeposit"`
	CDeposit       float64   `json:"cDeposit"`
	CPayment       float64   `json:"cPayment"`
	StartDate      string    `json:"startDate"`
	EndDate        string    `json:"endDate"`
	Status         string    `json:"status"`
	OwnerOrg       string    `json:"owner_org"`
	EscrowID       string    `json:"escrow_id"`
	OfferDetails   DataOffer `json:"offer_details"`
	AgreementID    string    `json:"agreement_id"`
}

// historicalOfferRequest type
type historicalOfferRequest struct {
	DocType                string              `json:"docType"`
	OfferRequestID         string              `json:"offer_request_id"`
	OfferID                string              `json:"offer_id"`
	DataConsumer           string              `json:"dataConsumer"`
	DataProvider           string              `json:"dataProvider"`
	Price                  float64             `json:"price"`
	PDeposit               float64             `json:"pDeposit"`
	CDeposit               float64             `json:"cDeposit"`
	CPayment               float64             `json:"cPayment"`
	StartDate              string              `json:"startDate"`
	EndDate                string              `json:"endDate"`
	Status                 string              `json:"status"`
	OwnerOrg               string              `json:"owner_org"`
	EscrowID               string              `json:"escrow_id"`
	HistoricalOfferDetails HistoricalDataOffer `json:"historical_details"`
	AgreementID            string              `json:"agreement_id"`
}

// Costs type
type Costs struct {
	FalsifyCount          int     `json:"falsifyCount"`
	LatencyCount          int     `json:"latencyCount"`
	CreatedAt             string  `json:"createdAt"`
	DocType               string  `json:"docType"`
	ID                    string  `json:"id"`
	Agreement             string  `json:"agreement"`
	ProviderReimbursement float64 `json:"providerReimbursement"`
	ConsumerRefund        float64 `json:"consumerRefund"`
	EscrowID              string  `json:"escrow_id"`
	DataConsumer          string  `json:"dataConsumer"`
	DataProvider          string  `json:"dataProvider"`
	OfferRequestID        string  `json:"offer_request_id"`
}

// Fees type
type Fees struct {
	DocType string  `json:"docType"`
	ID      string  `json:"id"`
	EP      float64 `json:"ep"`
	EC      float64 `json:"ec"`
}

const (
	DOC_TYPE          = "docType"
	JOURNEY           = "journey"
	DATA_OFFER        = "dataoffer"
	OFFER_REQUEST     = "offerRequest"
	ACTIVE            = "ACTIVE"
	CREATED           = "CREATED"
	REJECTED          = "REJECTED"
	ESCROW            = "escrow"
	DATA_HASH         = "offer_data_hash"
	DATA_AGREEMENT    = "data_agreement"
	DATA_HASH_VALUE   = "data_hash_value"
	REVOKED           = "REVOKED"
	COST              = "cost"
	CODE_A            = "101"
	CODE_B            = "102"
	CODE_C            = "103"
	LATENCY_TIME      = 0.15
	EXPIRED           = "EXPIRED"
	HISTORICALOFFER   = "historicalOffer"
	HISTORICALREQUEST = "historicalRequest"
)



type AgreementHash struct {
	Hashes    []*OfferDataHash `json:"hashes"`
	Agreement *DataAgreement   `json:"agreement"`
}
