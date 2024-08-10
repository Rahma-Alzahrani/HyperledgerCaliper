package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/flogging"
)

var logger = flogging.MustGetLogger("datamanagement_cc")

// dataManagement the chaincode interface implementation to manage
type dataManagement struct {
	contractapi.Contract
}

// Init initialize chaincode with mapping between actions and real methods
func (pm *dataManagement) InitLedger(ctx contractapi.TransactionContextInterface) error {

	var fee Fees

	fee.DocType = "fees"
	fee.ID = "feeAll"
	fee.EC, _ = strconv.ParseFloat("50.00", 64)
	fee.EP, _ = strconv.ParseFloat("50.00", 64)
	feeAsBytes, _ := json.Marshal(fee)
	err := ctx.GetStub().PutState(fee.ID, feeAsBytes)
	if err != nil {
		return fmt.Errorf("Failed to add Fee catch: %s", fee.ID)
	}
	return nil
}

// insertDataOffer inserts DataOffer into ledger
func (pm *dataManagement) JourneySchedule(ctx contractapi.TransactionContextInterface, journeySchedule string) (map[string]interface{}, error) {
	var details JourneySchedule
	logger.Info(journeySchedule)
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	err := json.Unmarshal([]byte(journeySchedule), &details)
	if err != nil {
		response["message"] = "Error occurred while unmarshalling"
		return response, fmt.Errorf("Failed while unmarshalling JourneySchedule: %s", err.Error())
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return response, fmt.Errorf("Failed to get verified OrgID: %v", err)
	}

	// Set the necessary details
	details.DocType = JOURNEY
	details.Type = clientOrgID

	// Store the data in the ledger
	detailsAsBytes, _ := json.Marshal(details)
	err = ctx.GetStub().PutState(details.UID, detailsAsBytes)
	if err != nil {
		response["message"] = "Error occurred while storing data"
		return response, fmt.Errorf("Failed to add JourneySchedule: %s", details.UID)
	}

	response["message"] = "JourneySchedule added successfully"
	return response, nil // Return success
}

// UpdateJourneySchedule Updates JourneySchedule into ledger
func (pm *dataManagement) UpdateJourneySchedule(ctx contractapi.TransactionContextInterface, journeySchedule string) (map[string]interface{}, error) {

	var details JourneySchedule
	logger.Info(journeySchedule)
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	err := json.Unmarshal([]byte(journeySchedule), &details)
	if err != nil {
		response["message"] = "Error occured while Unmarshal"
		return response, fmt.Errorf("Failed while unmarshling Order. %s", err.Error())
	}
	logger.Info(details)
	exists, err := pm.ObjectExists(ctx, details.UID)
	if err != nil {
		response["message"] = "failed to get journey"
		return response, fmt.Errorf("failed to get journey: %v", err)
	}
	if !exists {
		response["message"] = "No Record found"
		return response, fmt.Errorf("Journey does not exists")
	}

	var _journey JourneySchedule
	_journeyAsBytes, _ := ctx.GetStub().GetState(details.UID)
	err = json.Unmarshal(_journeyAsBytes, &_journey)
	if err != nil {
		return response, err
	}
	details.DocType = JOURNEY
	detailsAsBytes, _ := json.Marshal(details)
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return response, fmt.Errorf("failed to get verified OrgID: %v", err)
	}
	logger.Info("clientOrgID : %s", clientOrgID)
	//logger.Info("details.OwnerOrg : %s", details.OwnerOrg)

	err = ctx.GetStub().PutState(details.UID, detailsAsBytes)
	if err != nil {
		response["message"] = "Error occured while updating data"
		return response, fmt.Errorf("Failed to update journey catch: %s", details.UID)
	}

	return response, nil
}

// GetAllJourney to retrive all journey
func (pm *dataManagement) GetAllJourney(ctx contractapi.TransactionContextInterface) ([]*JourneySchedule, error) {

	queryString := fmt.Sprintf(`{"selector":{"docType":"%s"}}`, JOURNEY)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var journeySchedule []*JourneySchedule
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var journey JourneySchedule
		err = json.Unmarshal(queryResponse.Value, &journey)
		if err != nil {
			return nil, err
		}
		logger.Info(journey.UID)

		journeySchedule = append(journeySchedule, &journey)
	}

	return journeySchedule, nil

}

//GetJourneyByUID to retrive by uid

func (pm *dataManagement) GetJourneyByUID(ctx contractapi.TransactionContextInterface, key string) (*JourneySchedule, error) {

	assetJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", key)
	}

	var asset JourneySchedule
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (pm *dataManagement) InsertTestHistoricalDataOffer(ctx contractapi.TransactionContextInterface, dataOffer string) (map[string]interface{}, error) {
	var details HistoricalDataOffer
	logger.Info(dataOffer)
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	err := json.Unmarshal([]byte(dataOffer), &details)
	if err != nil {
		response["message"] = "Error occurred while Unmarshal"
		return response, fmt.Errorf("Failed while unmarshaling Order %s", err.Error())
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return response, fmt.Errorf("Failed to get verified OrgID: %v", err)
	}

	details.DocType = HISTORICALOFFER
	details.OwnerOrg = clientOrgID
	detailsAsBytes, _ := json.Marshal(details)

	// Store the offer, overwriting if it already exists
	err = ctx.GetStub().PutState(details.ID, detailsAsBytes)
	if err != nil {
		response["message"] = "Error occurred while storing data"
		return response, fmt.Errorf("Failed to add DataOffer catch: %s", details.ID)
	}

	return response, nil
}

// insertDataOffer inserts DataOffer into ledger
func (pm *dataManagement) InsertDataOffer(ctx contractapi.TransactionContextInterface, dataOffer string) (map[string]interface{}, error) {

	var details DataOffer
	logger.Info(dataOffer)
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	err := json.Unmarshal([]byte(dataOffer), &details)
	if err != nil {
		response["message"] = "Error occured while Unmarshal"
		return response, fmt.Errorf("Failed while unmarshling Order %s", err.Error())
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return response, fmt.Errorf("failed to get verified OrgID: %v", err)
	}
	details.DocType = DATA_OFFER
	details.OwnerOrg = clientOrgID
	detailsAsBytes, _ := json.Marshal(details)
	err = ctx.GetStub().PutState(details.ID, detailsAsBytes)
	if err != nil {
		response["message"] = "Error occured while storing data"
		return response, fmt.Errorf("Failed to add DataOffer catch: %s", details.ID)
	}

	return response, nil

}

// UpdateDataOffer Updates DataOffer into ledger
func (pm *dataManagement) UpdateDataOffer(ctx contractapi.TransactionContextInterface, dataOffer string) (map[string]interface{}, error) {

	var details DataOffer
	logger.Info(dataOffer)
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	err := json.Unmarshal([]byte(dataOffer), &details)
	if err != nil {
		response["message"] = "Error occured while Unmarshal"
		return response, fmt.Errorf("Failed while unmarshling Order. %s", err.Error())
	}
	logger.Info(details)
	exists, err := pm.ObjectExists(ctx, details.ID)
	if err != nil {
		response["message"] = "failed to get Offer"
		return response, fmt.Errorf("failed to get Offer: %v", err)
	}
	if !exists {
		response["message"] = "No Record found"
		return response, fmt.Errorf("Data Offer does not exists")
	}

	var _offer DataOffer
	_offerAsBytes, _ := ctx.GetStub().GetState(details.ID)
	err = json.Unmarshal(_offerAsBytes, &_offer)
	if err != nil {
		return response, err
	}
	details.DocType = DATA_OFFER
	details.OwnerOrg = _offer.OwnerOrg
	details.Creator = _offer.Creator

	detailsAsBytes, _ := json.Marshal(details)
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return response, fmt.Errorf("failed to get verified OrgID: %v", err)
	}
	logger.Info("clientOrgID : %s", clientOrgID)
	logger.Info("details.OwnerOrg : %s", details.OwnerOrg)

	err = ctx.GetStub().PutState(details.ID, detailsAsBytes)
	if err != nil {
		response["message"] = "Error occured while updating data"
		return response, fmt.Errorf("Failed to update Data offer catch: %s", details.ID)
	}

	return response, nil
}

// GetAllOffers retrieves all offers
func (pm *dataManagement) GetAllOffers(ctx contractapi.TransactionContextInterface, creator string) ([]QueryResult, error) {

	var queryString string
	logger.Info(creator)
	if len(creator) == 0 {
		queryString = fmt.Sprintf(`{"selector":{"docType":"%s"}}`, DATA_OFFER)
	} else {
		queryString = fmt.Sprintf(`{"selector":{"docType":"%s","creator":"%s"}}`, DATA_OFFER, creator)
	}

	logger.Info(queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var results []QueryResult

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var dataOffer DataOffer
		err = json.Unmarshal(queryResponse.Value, &dataOffer)
		if err != nil {
			return nil, err
		}

		queryResult := QueryResult{Key: queryResponse.Key, Record: &dataOffer}
		results = append(results, queryResult)
	}

	return results, nil
}

// GetOffer retrieves a single offer
func (pm *dataManagement) GetOffer(ctx contractapi.TransactionContextInterface, offerID string) (*DataOffer, error) {
	var offer DataOffer
	offerAsBytes, err := ctx.GetStub().GetState(offerID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if offerAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", offerID)
	}
	err = json.Unmarshal(offerAsBytes, &offer)
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

// CreateOfferRequest Creating OfferRequest and escrow
func (pm *dataManagement) CreateOfferRequest(ctx contractapi.TransactionContextInterface, offerRequest string) (map[string]interface{}, error) {
	var _offerRequest OfferRequest
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	err := json.Unmarshal([]byte(offerRequest), &_offerRequest)
	if err != nil {
		response["message"] = "Failed to unmarshal offer request"
		return response, fmt.Errorf("FAILED TO UNMARSHAL OFFER REQUEST %s", err.Error())
	}

	exist, err := pm.ObjectExists(ctx, _offerRequest.OfferID)
	if err != nil {
		response["message"] = "Failed to get offer"
		return response, fmt.Errorf("FAILED TO GET OFFER %s", err.Error())
	}
	if !exist {
		response["message"] = fmt.Sprintf("No such offer exist with offer id %s", _offerRequest.OfferID)
		return response, fmt.Errorf("NO SUCH OFFER EXIST WITH OFFER ID %s: %s", _offerRequest.OfferID, err.Error())
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		response["message"] = "Failed to get offer ID"
		return response, fmt.Errorf("FAILED TO GET OFFER ID %s", err.Error())
	}
	var offer DataOffer
	offerAsBytes, _ := ctx.GetStub().GetState(_offerRequest.OfferID)

	err = json.Unmarshal(offerAsBytes, &offer)
	if err != nil {
		response["message"] = "Failed to unmarshal"
		return response, fmt.Errorf("FAILED TO UNMARSHAL %s", err.Error())
	}

	escrowID := rotateLeft(_offerRequest.OfferRequestID, 5)
	_offerRequest.DocType = OFFER_REQUEST
	_offerRequest.OwnerOrg = clientOrgID
	_offerRequest.DataProvider = offer.Creator
	_offerRequest.Status = CREATED
	_offerRequest.EscrowID = escrowID
	_offerRequest.PDeposit = offer.Deposit
	_offerRequest.OfferDetails = offer

	offerRequestAsBytes, _ := json.Marshal(_offerRequest)

	err = ctx.GetStub().PutState(_offerRequest.OfferRequestID, offerRequestAsBytes)
	if err != nil {
		response["message"] = "Failed to do put state for offer request"
		return response, err
	}

	// create escrow
	escrow := Escrow{
		Status:          CREATED,
		Consumer:        _offerRequest.DataConsumer,
		Provider:        _offerRequest.DataProvider,
		DocType:         ESCROW,
		ID:              escrowID,
		ProviderDeposit: offer.Deposit,
		ConsumerDeposit: _offerRequest.CDeposit,
		ConsumerPayment: _offerRequest.Price,
		Released:        false,
		OfferRequestID:  _offerRequest.OfferRequestID,
		OfferID:         _offerRequest.OfferID,
	}
	escrowAsBytes, _ := json.Marshal(escrow)
	err = ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
	if err != nil {
		response["message"] = "Failed to do put state for escrow"
		return response, err
	}
	response["message"] = fmt.Sprintf("Offer Request ID: %s, Escrow ID: %s", _offerRequest.OfferRequestID, escrow.ID)
	return response, nil
}

func (pm *dataManagement) CreateOfferRequest2(ctx contractapi.TransactionContextInterface, offerRequest string) (map[string]interface{}, error) {

	var _offerRequest OfferRequest
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	err := json.Unmarshal([]byte(offerRequest), &_offerRequest)
	if err != nil {
		response["message"] = "Failed to unmarshal offer request"
		return response, fmt.Errorf("FAILED TO UNMARSHAL OFFER REQUEST %s", err.Error())
	}

	exist, err := pm.ObjectExists(ctx, _offerRequest.OfferID)
	if err != nil {
		response["message"] = "Failed to get offer"
		return response, fmt.Errorf("FAILED TO GET OFFER %s", err.Error())
	}
	if !exist {
		response["message"] = fmt.Sprintf("No such offer exist with offer id %s", _offerRequest.OfferID)
		return response, fmt.Errorf("NO SUCH OFFER EXIST WITH OFFER ID %s: %s", _offerRequest.OfferID, err.Error())
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		response["message"] = "Failed to get offer ID"
		return response, fmt.Errorf("FAILED TO GET OFFER ID %s", err.Error())
	}
	var offer DataOffer
	offerAsBytes, _ := ctx.GetStub().GetState(_offerRequest.OfferID)
	err = json.Unmarshal(offerAsBytes, &offer)
	if err != nil {
		response["message"] = "Failed to unmarshal"
		return response, fmt.Errorf("FAILED TO UNMARSHAL %s", err.Error())
	}

	// Update the validity field to false
	offer.IsActive = false

	// Marshal the updated offer back to JSON including other fields
	updatedOfferAsBytes, err := json.Marshal(offer)
	if err != nil {
		response["message"] = "Failed to marshal updated offer"
		return response, fmt.Errorf("FAILED TO MARSHAL UPDATED OFFER: %s", err.Error())
	}

	// Put the updated offer back into the state
	err = ctx.GetStub().PutState(_offerRequest.OfferID, updatedOfferAsBytes)
	if err != nil {
		response["message"] = "Failed to put state for updated offer"
		return response, err
	}

	err = json.Unmarshal(offerAsBytes, &offer)
	if err != nil {
		response["message"] = "Failed to unmarshal"
		return response, fmt.Errorf("FAILED TO UNMARSHAL %s", err.Error())
	}

	escrowID := rotateLeft(_offerRequest.OfferRequestID, 5)
	_offerRequest.DocType = OFFER_REQUEST
	_offerRequest.OwnerOrg = clientOrgID
	_offerRequest.DataProvider = offer.Creator
	_offerRequest.Status = CREATED
	_offerRequest.EscrowID = escrowID
	_offerRequest.PDeposit = offer.Deposit
	_offerRequest.OfferDetails = offer
	_offerRequest.EndDate = _offerRequest.OfferDetails.Arrival_time
	_offerRequest.StartDate = _offerRequest.OfferDetails.Depart_time

	offerRequestAsBytes, _ := json.Marshal(_offerRequest)

	err = ctx.GetStub().PutState(_offerRequest.OfferRequestID, offerRequestAsBytes)
	if err != nil {
		response["message"] = "Failed to do put state for offer request"
		return response, err
	}

	// create escrow
	escrow := Escrow{
		Status:          CREATED,
		Consumer:        _offerRequest.DataConsumer,
		Provider:        _offerRequest.DataProvider,
		DocType:         ESCROW,
		ID:              escrowID,
		ProviderDeposit: offer.Deposit,
		ConsumerDeposit: _offerRequest.CDeposit,
		ConsumerPayment: _offerRequest.Price,
		Released:        false,
		OfferRequestID:  _offerRequest.OfferRequestID,
		OfferID:         _offerRequest.OfferID,
	}
	escrowAsBytes, _ := json.Marshal(escrow)
	err = ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
	if err != nil {
		response["message"] = "Failed to do put state for escrow"
		return response, err
	}
	response["message"] = fmt.Sprintf("Offer Request ID: %s, Escrow ID: %s", _offerRequest.OfferRequestID, escrow.ID)
	return response, nil
}

// GetAllOfferRequest to retrive all requests
func (pm *dataManagement) GetAllOfferRequest(ctx contractapi.TransactionContextInterface) ([]*OfferRequest, error) {

	queryString := fmt.Sprintf(`{"selector":{"docType":"%s"}}`, OFFER_REQUEST)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var offerRequest []*OfferRequest
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var request OfferRequest
		err = json.Unmarshal(queryResponse.Value, &request)
		if err != nil {
			return nil, err
		}
		logger.Info(request.OfferID)

		offerRequest = append(offerRequest, &request)
	}

	return offerRequest, nil

}

// CreateHistoricalOfferRequest Creating OfferRequest and escrow
func (pm *dataManagement) CreateHistoricalOfferRequest(ctx contractapi.TransactionContextInterface, historicalRequest string) (map[string]interface{}, error) {
	var _historicalRequest historicalOfferRequest
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	err := json.Unmarshal([]byte(historicalRequest), &_historicalRequest)
	if err != nil {
		response["message"] = "Failed to unmarshal offer request"
		return response, fmt.Errorf("FAILED TO UNMARSHAL OFFER REQUEST %s", err.Error())
	}

	exist, err := pm.ObjectExists(ctx, _historicalRequest.OfferID)
	if err != nil {
		response["message"] = "Failed to get offer"
		return response, fmt.Errorf("FAILED TO GET OFFER %s", err.Error())
	}
	if !exist {
		response["message"] = fmt.Sprintf("No such offer exist with offer id %s", _historicalRequest.OfferID)
		return response, fmt.Errorf("NO SUCH OFFER EXIST WITH OFFER ID %s: %s", _historicalRequest.OfferID, err.Error())
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		response["message"] = "Failed to get offer ID"
		return response, fmt.Errorf("FAILED TO GET OFFER ID %s", err.Error())
	}
	var offer HistoricalDataOffer
	offerAsBytes, _ := ctx.GetStub().GetState(_historicalRequest.OfferID)

	err = json.Unmarshal(offerAsBytes, &offer)
	if err != nil {
		response["message"] = "Failed to unmarshal"
		return response, fmt.Errorf("FAILED TO UNMARSHAL %s", err.Error())
	}

	escrowID := rotateLeft(_historicalRequest.OfferRequestID, 5)
	_historicalRequest.DocType = HISTORICALREQUEST
	_historicalRequest.OwnerOrg = clientOrgID
	_historicalRequest.DataProvider = offer.Creator
	_historicalRequest.Status = CREATED
	_historicalRequest.EscrowID = escrowID
	_historicalRequest.PDeposit = offer.Deposit
	_historicalRequest.HistoricalOfferDetails = offer

	// Iterate through the original slice and add unique offer IDs to the map

	// Replace the original slice with the unique offer IDs

	offerRequestAsBytes, _ := json.Marshal(_historicalRequest)

	err = ctx.GetStub().PutState(_historicalRequest.OfferRequestID, offerRequestAsBytes)
	if err != nil {
		response["message"] = "Failed to do put state for offer request"
		return response, err
	}

	// create escrow
	escrow := Escrow{
		Status:          CREATED,
		Consumer:        _historicalRequest.DataConsumer,
		Provider:        _historicalRequest.DataProvider,
		DocType:         ESCROW,
		ID:              escrowID,
		ProviderDeposit: offer.Deposit,
		ConsumerDeposit: _historicalRequest.CDeposit,
		ConsumerPayment: _historicalRequest.Price,
		Released:        false,
		OfferRequestID:  _historicalRequest.OfferRequestID,
		OfferID:         _historicalRequest.OfferID,
	}
	escrowAsBytes, _ := json.Marshal(escrow)
	err = ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
	if err != nil {
		response["message"] = "Failed to do put state for escrow"
		return response, err
	}
	response["message"] = fmt.Sprintf("Offer Request ID: %s, Escrow ID: %s", _historicalRequest.OfferRequestID, escrow.ID)
	return response, nil
}

func (pm *dataManagement) CreateHistoricalOfferRequestTest(ctx contractapi.TransactionContextInterface, historicalRequest string) (map[string]interface{}, error) {
	var _historicalRequest historicalOfferRequest
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	err := json.Unmarshal([]byte(historicalRequest), &_historicalRequest)
	if err != nil {
		response["message"] = "Failed to unmarshal offer request"
		return response, fmt.Errorf("FAILED TO UNMARSHAL OFFER REQUEST %s", err.Error())
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		response["message"] = "Failed to get offer ID"
		return response, fmt.Errorf("FAILED TO GET OFFER ID %s", err.Error())
	}
	var offer HistoricalDataOffer
	offerAsBytes, _ := ctx.GetStub().GetState(_historicalRequest.OfferID)

	err = json.Unmarshal(offerAsBytes, &offer)
	if err != nil {
		// Create a new offer if the offer ID doesn't exist
		offer = HistoricalDataOffer{
			// Populate the fields accordingly
		}

		// Marshal and save the newly created offer
		offerAsBytes, err = json.Marshal(offer)
		if err != nil {
			response["message"] = "Failed to marshal offer"
			return response, fmt.Errorf("FAILED TO MARSHAL OFFER %s", err.Error())
		}
		err = ctx.GetStub().PutState(_historicalRequest.OfferID, offerAsBytes)
		if err != nil {
			response["message"] = "Failed to put state for offer"
			return response, fmt.Errorf("FAILED TO PUT STATE FOR OFFER %s", err.Error())
		}
	}

	escrowID := rotateLeft(_historicalRequest.OfferRequestID, 5)
	_historicalRequest.DocType = HISTORICALREQUEST
	_historicalRequest.OwnerOrg = clientOrgID
	_historicalRequest.DataProvider = offer.Creator
	_historicalRequest.Status = CREATED
	_historicalRequest.EscrowID = escrowID
	_historicalRequest.PDeposit = offer.Deposit
	_historicalRequest.HistoricalOfferDetails = offer

	// Any other necessary operations for the historical request

	// Create escrow
	escrow := Escrow{
		Status:          CREATED,
		Consumer:        _historicalRequest.DataConsumer,
		Provider:        _historicalRequest.DataProvider,
		DocType:         ESCROW,
		ID:              escrowID,
		ProviderDeposit: offer.Deposit,
		ConsumerDeposit: _historicalRequest.CDeposit,
		ConsumerPayment: _historicalRequest.Price,
		Released:        false,
		OfferRequestID:  _historicalRequest.OfferRequestID,
		OfferID:         _historicalRequest.OfferID,
	}
	escrowAsBytes, _ := json.Marshal(escrow)
	err = ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
	if err != nil {
		response["message"] = "Failed to put state for escrow"
		return response, fmt.Errorf("FAILED TO PUT STATE FOR ESCROW %s", err.Error())
	}
	response["message"] = fmt.Sprintf("Offer Request ID: %s, Escrow ID: %s", _historicalRequest.OfferRequestID, escrow.ID)
	return response, nil
}

// func (pm *dataManagement) CreateHistoricalOfferRequest(ctx contractapi.TransactionContextInterface, historicalRequest string) (map[string]interface{}, error) {
// 	var _historicalRequest historicalOfferRequest
// 	response := make(map[string]interface{})
// 	response["txId"] = ctx.GetStub().GetTxID()
// 	err := json.Unmarshal([]byte(historicalRequest), &_historicalRequest)
// 	if err != nil {
// 		response["message"] = "Failed to unmarshal offer request"
// 		return response, fmt.Errorf("FAILED TO UNMARSHAL OFFER REQUEST %s", err.Error())
// 	}

// 	for _, offerID := range _historicalRequest.OfferID {
// 		exist, err := pm.ObjectExists(ctx, offerID)
// 		if err != nil {
// 			response["message"] = "Failed to get offer"
// 			return response, fmt.Errorf("FAILED TO GET OFFER %s", err.Error())
// 		}
// 		if !exist {
// 			response["message"] = fmt.Sprintf("No such offer exists with offer ID %s", offerID)
// 			return response, fmt.Errorf("NO SUCH OFFER EXISTS WITH OFFER ID %s: %s", offerID, err.Error())
// 		}

// 		clientOrgID, err := getClientOrgID(ctx, false)
// 		if err != nil {
// 			response["message"] = "Failed to get offer ID"
// 			return response, fmt.Errorf("FAILED TO GET OFFER ID %s", err.Error())
// 		}

// 		var offer HistoricalDataOffer
// 		offerAsBytes, _ := ctx.GetStub().GetState(offerID)

// 		err = json.Unmarshal(offerAsBytes, &offer)
// 		if err != nil {
// 			response["message"] = "Failed to unmarshal"
// 			return response, fmt.Errorf("FAILED TO UNMARSHAL %s", err.Error())
// 		}

// 		escrowID := rotateLeft(_historicalRequest.OfferRequestID, 5)
// 		_historicalRequest.DocType = HISTORICALREQUEST
// 		_historicalRequest.OwnerOrg = clientOrgID
// 		_historicalRequest.DataProvider = offer.Creator
// 		_historicalRequest.Status = CREATED
// 		_historicalRequest.EscrowID = escrowID
// 		_historicalRequest.PDeposit = offer.Deposit
// 		_historicalRequest.HistoricalOfferDetails = offer

// 		offerRequestAsBytes, _ := json.Marshal(_historicalRequest)

// 		err = ctx.GetStub().PutState(_historicalRequest.OfferRequestID, offerRequestAsBytes)
// 		if err != nil {
// 			response["message"] = "Failed to do put state for offer request"
// 			return response, err
// 		}

// 		// create escrow
// 		escrow := Escrow{
// 			Status:          CREATED,
// 			Consumer:        _historicalRequest.DataConsumer,
// 			Provider:        _historicalRequest.DataProvider,
// 			DocType:         ESCROW,
// 			ID:              escrowID,
// 			ProviderDeposit: offer.Deposit,
// 			ConsumerDeposit: _historicalRequest.CDeposit,
// 			ConsumerPayment: _historicalRequest.Price,
// 			Released:        false,
// 			OfferRequestID:  _historicalRequest.OfferRequestID,
// 			OfferID:         offerID,
// 		}
// 		escrowAsBytes, _ := json.Marshal(escrow)
// 		err = ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
// 		if err != nil {
// 			response["message"] = "Failed to do put state for escrow"
// 			return response, err
// 		}
// 		response["message"] = fmt.Sprintf("Offer Request ID: %s, Escrow ID: %s", _historicalRequest.OfferRequestID, escrow.ID)
// 	}

// 	return response, nil
// }

// getHistoricalOfferRequest to retrive offer request
func (pm *dataManagement) GetHistoricalOfferRequest(ctx contractapi.TransactionContextInterface, offerRequestID string) (*historicalOfferRequest, error) {
	offerRequestAsBytes, err := ctx.GetStub().GetState(offerRequestID)
	if err != nil {
		return nil, fmt.Errorf("FAILED TO GET OFFER REQUEST %s", err.Error())
	}
	var offerRequest historicalOfferRequest
	err = json.Unmarshal(offerRequestAsBytes, &offerRequest)
	if err != nil {
		return nil, fmt.Errorf("FAILED TO UNMARSHAL OFFER REQUEST %s", err.Error())
	}
	return &offerRequest, nil
}

// GetAllOfferRequest to retrive all requests
func (pm *dataManagement) GetAllHistoricalOfferRequest(ctx contractapi.TransactionContextInterface) ([]*historicalOfferRequest, error) {

	queryString := fmt.Sprintf(`{"selector":{"docType":"%s"}}`, HISTORICALREQUEST)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var offerRequest []*historicalOfferRequest
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var request historicalOfferRequest
		err = json.Unmarshal(queryResponse.Value, &request)
		if err != nil {
			return nil, err
		}
		logger.Info(request.OfferID)

		offerRequest = append(offerRequest, &request)
	}

	return offerRequest, nil

}

func (pm *dataManagement) GetHistoricalOfferRequestByID(ctx contractapi.TransactionContextInterface, key string) (*historicalOfferRequest, error) {

	assetJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", key)
	}

	var asset historicalOfferRequest
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (pm *dataManagement) CreateAndAcceptHistoricalOfferRequest(ctx contractapi.TransactionContextInterface, historicalRequest string) (map[string]interface{}, error) {
	var _historicalRequest historicalOfferRequest
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")
	err := json.Unmarshal([]byte(historicalRequest), &_historicalRequest)
	if err != nil {
		response["message"] = "Failed to unmarshal offer request"
		return response, fmt.Errorf("FAILED TO UNMARSHAL OFFER REQUEST %s", err.Error())
	}

	exist, err := pm.ObjectExists(ctx, _historicalRequest.OfferID)
	if err != nil {
		response["message"] = "Failed to get offer"
		return response, fmt.Errorf("FAILED TO GET OFFER %s", err.Error())
	}
	if !exist {
		response["message"] = fmt.Sprintf("No such offer exists with offer id %s", _historicalRequest.OfferID)
		return response, fmt.Errorf("NO SUCH OFFER EXISTS WITH OFFER ID %s: %s", _historicalRequest.OfferID, err.Error())
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		response["message"] = "Failed to get offer ID"
		return response, fmt.Errorf("FAILED TO GET OFFER ID %s", err.Error())
	}

	var offer HistoricalDataOffer
	offerAsBytes, _ := ctx.GetStub().GetState(_historicalRequest.OfferID)

	err = json.Unmarshal(offerAsBytes, &offer)
	if err != nil {
		response["message"] = "Failed to unmarshal"
		return response, fmt.Errorf("FAILED TO UNMARSHAL %s", err.Error())
	}

	escrowID := rotateLeft(_historicalRequest.OfferRequestID, 5)
	_historicalRequest.DocType = HISTORICALREQUEST
	_historicalRequest.OwnerOrg = clientOrgID
	_historicalRequest.DataProvider = offer.Creator
	_historicalRequest.Status = CREATED
	_historicalRequest.EscrowID = escrowID
	_historicalRequest.PDeposit = offer.Deposit
	_historicalRequest.HistoricalOfferDetails = offer

	// Create a new DataAgreement
	agreementID := rotateLeft(_historicalRequest.OfferRequestID, 10)
	_historicalRequest.AgreementID = agreementID
	_historicalRequest.Status = ACTIVE

	escrow := Escrow{
		Status:          ACTIVE,
		Consumer:        _historicalRequest.DataConsumer,
		Provider:        _historicalRequest.DataProvider,
		DocType:         ESCROW,
		ID:              escrowID,
		ProviderDeposit: offer.Deposit,
		ConsumerDeposit: _historicalRequest.CDeposit,
		ConsumerPayment: _historicalRequest.Price,
		Released:        false,
		OfferRequestID:  _historicalRequest.OfferRequestID,
		OfferID:         _historicalRequest.OfferID,
		AgreementID:     agreementID,
		StartDate:       today,
		EndDate:         today,
	}

	dataAgreement := DataAgreement{
		Price:           offer.Price,
		DocType:         DATA_AGREEMENT,
		ID:              agreementID,
		DataProvider:    _historicalRequest.DataProvider,
		DataConsumer:    _historicalRequest.DataConsumer,
		EscrowID:        escrowID,
		State:           true, // Agreement accepted
		OfferRequestID:  _historicalRequest.OfferRequestID,
		OfferID:         _historicalRequest.OfferID,
		StartDate:       today,
		EndDate:         today,
		OfferDataHashID: []string{},
		ProviderDeposit: _historicalRequest.PDeposit,
		ConsumerDeposit: _historicalRequest.CDeposit,
	}

	dataAgreementAsBytes, err := json.Marshal(dataAgreement)
	if err != nil {
		response["message"] = "Failed to marshal Agreement"
		return response, err
	}

	err = ctx.GetStub().PutState(agreementID, dataAgreementAsBytes)
	if err != nil {
		response["message"] = "Failed to Putstate Agreement"
		return response, err
	}

	offerRequestAsBytes, _ := json.Marshal(_historicalRequest)

	err = ctx.GetStub().PutState(_historicalRequest.OfferRequestID, offerRequestAsBytes)
	if err != nil {
		response["message"] = "Failed to do put state for offer request"
		return response, err
	}

	escrowAsBytes, _ := json.Marshal(escrow)
	err = ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
	if err != nil {
		response["message"] = "Failed to do put state for escrow"
		return response, err
	}

	response["message"] = fmt.Sprintf("Offer Request ID: %s, Escrow ID: %s", _historicalRequest.OfferRequestID, escrow.ID)
	return response, nil
}

// AcceptOfferRequest to respond to data request and bild the agreement in case of acceptance
func (pm *dataManagement) AcceptHistoricalOfferRequest(ctx contractapi.TransactionContextInterface, offerID, offerRequestID string, isAccepted bool) (map[string]interface{}, error) {
	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")
	response := make(map[string]interface{})

	// Get the offer details from offerID
	var offer HistoricalDataOffer
	offerAsBytes, err := ctx.GetStub().GetState(offerID)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(offerAsBytes, &offer)
	if err != nil {
		return response, err
	}

	// Get the offer Request from the offerRequestID
	var historicalRequest historicalOfferRequest
	offerRequestAsBytes, err := ctx.GetStub().GetState(offerRequestID)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(offerRequestAsBytes, &historicalRequest)
	if err != nil {
		return response, err
	}
	// Get the escrow details
	escrowID := historicalRequest.EscrowID
	var escrow Escrow
	escrowAsBytes, err := ctx.GetStub().GetState(escrowID)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(escrowAsBytes, &escrow)
	if err != nil {
		return response, err
	}

	//newEndDate, _ := time.Parse("2006-01-02 15:04", today)
	//if hours > 0 {
	//	newEndDate = newEndDate.Add(time.Duration(hours) * time.Hour)
	//}
	//if minutes > 0 {
	//	newEndDate = newEndDate.Add(time.Duration(minutes) * time.Minute)
	//}

	// Business Logic Here
	if isAccepted {
		// agreementId := generateUUID()
		agreementId := rotateLeft(offerRequestID, 10)
		historicalRequest.AgreementID = agreementId
		historicalRequest.Status = ACTIVE
		escrow.Status = ACTIVE
		escrow.AgreementID = agreementId
		escrow.StartDate = today
		escrow.EndDate = today
		escrow.ProviderDeposit = offer.Deposit
		historicalRequest.HistoricalOfferDetails.StartDate = today
		historicalRequest.HistoricalOfferDetails.EndDate = today
		dataAgreement := DataAgreement{
			Price:           offer.Price,
			DocType:         DATA_AGREEMENT,
			ID:              agreementId,
			DataProvider:    historicalRequest.DataProvider,
			DataConsumer:    historicalRequest.DataConsumer,
			EscrowID:        escrowID,
			State:           isAccepted,
			OfferRequestID:  offerRequestID,
			OfferID:         offerID,
			StartDate:       today,
			EndDate:         today,
			OfferDataHashID: []string{},
			ProviderDeposit: historicalRequest.PDeposit,
			ConsumerDeposit: historicalRequest.CDeposit,
		}
		dataAgreementAsBytes, err := json.Marshal(dataAgreement)
		if err != nil {
			response["message"] = "Failed to marshal Agreement"
			return response, err
		}
		err = ctx.GetStub().PutState(agreementId, dataAgreementAsBytes)
		if err != nil {
			response["message"] = "Failed to Putstate Agreement"
			return response, err
		}
	} else {
		costId := rotateLeft(offerRequestID, 10)
		escrow.Status = REJECTED
		historicalRequest.Status = REJECTED
		escrow.Released = true
		cost := Costs{
			DocType:               COST,
			ID:                    costId,
			ProviderReimbursement: 0,
			ConsumerRefund:        (escrow.ConsumerDeposit + escrow.ConsumerPayment),
			EscrowID:              escrowID,
			DataConsumer:          escrow.Consumer,
			DataProvider:          escrow.Provider,
			OfferRequestID:        escrow.OfferRequestID,
		}
		costAsBytes, err := json.Marshal(cost)
		if err != nil {
			response["message"] = "Failed to marshal Cost"
			return response, err
		}
		err = ctx.GetStub().PutState(costId, costAsBytes)
		if err != nil {
			response["message"] = "Failed to Putstate Cost"
			return response, err
		}
	}

	offerRequestAsBytes, err = json.Marshal(historicalRequest)
	err = ctx.GetStub().PutState(historicalRequest.OfferRequestID, offerRequestAsBytes)
	if err != nil {
		return response, err
	}

	escrowAsBytes, err = json.Marshal(escrow)
	err = ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
	if err != nil {
		return response, err
	}

	return response, nil
}
func (pm *dataManagement) CreateAndAcceptOfferRequest(ctx contractapi.TransactionContextInterface, offerRequest string) (map[string]interface{}, error) {
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()

	var _offerRequest OfferRequest
	err := json.Unmarshal([]byte(offerRequest), &_offerRequest)
	if err != nil {
		response["message"] = "Failed to unmarshal offer request"
		return response, fmt.Errorf("FAILED TO UNMARSHAL OFFER REQUEST %s", err.Error())
	}

	exist, err := pm.ObjectExists(ctx, _offerRequest.OfferID)
	if err != nil || !exist {
		response["message"] = fmt.Sprintf("No such offer exists with offer id %s", _offerRequest.OfferID)
		return response, fmt.Errorf("NO SUCH OFFER EXISTS WITH OFFER ID %s: %s", _offerRequest.OfferID, err.Error())
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		response["message"] = "Failed to get offer ID"
		return response, fmt.Errorf("FAILED TO GET OFFER ID %s", err.Error())
	}

	var offer DataOffer
	offerAsBytes, _ := ctx.GetStub().GetState(_offerRequest.OfferID)
	err = json.Unmarshal(offerAsBytes, &offer)
	if err != nil {
		response["message"] = "Failed to unmarshal"
		return response, fmt.Errorf("FAILED TO UNMARSHAL %s", err.Error())
	}

	escrowID := rotateLeft(_offerRequest.OfferRequestID, 5)
	_offerRequest.DocType = OFFER_REQUEST
	_offerRequest.OwnerOrg = clientOrgID
	_offerRequest.DataProvider = offer.Creator
	_offerRequest.Status = ACTIVE // Automatically accept the offer
	_offerRequest.EscrowID = escrowID
	_offerRequest.PDeposit = offer.Deposit
	_offerRequest.OfferDetails = offer

	offerRequestAsBytes, _ := json.Marshal(_offerRequest)
	err = ctx.GetStub().PutState(_offerRequest.OfferRequestID, offerRequestAsBytes)
	if err != nil {
		response["message"] = "Failed to do put state for offer request"
		return response, err
	}

	// Create escrow directly and update its state
	escrow := Escrow{
		Status:          ACTIVE,
		Consumer:        _offerRequest.DataConsumer,
		Provider:        _offerRequest.DataProvider,
		DocType:         ESCROW,
		ID:              escrowID,
		ProviderDeposit: offer.Deposit,
		ConsumerDeposit: _offerRequest.CDeposit,
		ConsumerPayment: _offerRequest.Price,
		Released:        false,
		OfferRequestID:  _offerRequest.OfferRequestID,
		OfferID:         _offerRequest.OfferID,
	}

	escrowAsBytes, _ := json.Marshal(escrow)
	err = ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
	if err != nil {
		response["message"] = "Failed to do put state for escrow"
		return response, err
	}

	// Create agreement directly and update its state
	agreementId := rotateLeft(_offerRequest.OfferRequestID, 10)
	dataAgreement := DataAgreement{
		Price:           offer.Price,
		DocType:         DATA_AGREEMENT,
		ID:              agreementId,
		DataProvider:    _offerRequest.DataProvider,
		DataConsumer:    _offerRequest.DataConsumer,
		EscrowID:        escrowID,
		State:           true, // Set to true as the offer is accepted
		OfferRequestID:  _offerRequest.OfferRequestID,
		OfferID:         _offerRequest.OfferID,
		StartDate:       _offerRequest.StartDate,
		EndDate:         _offerRequest.EndDate,
		OfferDataHashID: []string{},
		ProviderDeposit: _offerRequest.PDeposit,
		ConsumerDeposit: _offerRequest.CDeposit,
	}

	dataAgreementAsBytes, err := json.Marshal(dataAgreement)
	if err != nil {
		response["message"] = "Failed to marshal Agreement"
		return response, err
	}

	err = ctx.GetStub().PutState(agreementId, dataAgreementAsBytes)
	if err != nil {
		response["message"] = "Failed to Putstate Agreement"
		return response, err
	}

	response["message"] = fmt.Sprintf("Offer Request ID: %s, Escrow ID: %s, Agreement ID: %s", _offerRequest.OfferRequestID, escrow.ID, agreementId)
	return response, nil
}

// AcceptOfferRequest to respond to data request and bild the agreement in case of acceptance
func (pm *dataManagement) AcceptOfferRequest(ctx contractapi.TransactionContextInterface, offerID, offerRequestID string, isAccepted bool) (map[string]interface{}, error) {

	response := make(map[string]interface{})

	// Get the offer details from offerID
	var offer DataOffer
	offerAsBytes, err := ctx.GetStub().GetState(offerID)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(offerAsBytes, &offer)
	if err != nil {
		return response, err
	}

	// Get the offer Request from the offerRequestID
	var offerRequest OfferRequest
	offerRequestAsBytes, err := ctx.GetStub().GetState(offerRequestID)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(offerRequestAsBytes, &offerRequest)
	if err != nil {
		return response, err
	}
	// Get the escrow details
	escrowID := offerRequest.EscrowID
	var escrow Escrow
	escrowAsBytes, err := ctx.GetStub().GetState(escrowID)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(escrowAsBytes, &escrow)
	if err != nil {
		return response, err
	}

	// Business Logic Here
	if isAccepted {
		// agreementId := generateUUID()
		agreementId := rotateLeft(offerRequestID, 10)
		offerRequest.AgreementID = agreementId
		offerRequest.Status = ACTIVE
		escrow.Status = ACTIVE
		escrow.AgreementID = agreementId
		escrow.StartDate = offerRequest.StartDate
		escrow.EndDate = offerRequest.EndDate
		escrow.ProviderDeposit = offer.Deposit
		dataAgreement := DataAgreement{
			Price:           offer.Price,
			DocType:         DATA_AGREEMENT,
			ID:              agreementId,
			DataProvider:    offerRequest.DataProvider,
			DataConsumer:    offerRequest.DataConsumer,
			EscrowID:        escrowID,
			State:           isAccepted,
			OfferRequestID:  offerRequestID,
			OfferID:         offerID,
			StartDate:       offerRequest.StartDate,
			EndDate:         offerRequest.EndDate,
			OfferDataHashID: []string{},
			ProviderDeposit: offerRequest.PDeposit,
			ConsumerDeposit: offerRequest.CDeposit,
		}
		dataAgreementAsBytes, err := json.Marshal(dataAgreement)
		if err != nil {
			response["message"] = "Failed to marshal Agreement"
			return response, err
		}
		err = ctx.GetStub().PutState(agreementId, dataAgreementAsBytes)
		if err != nil {
			response["message"] = "Failed to Putstate Agreement"
			return response, err
		}
	} else {
		costId := rotateLeft(offerRequestID, 10)
		escrow.Status = REJECTED
		offerRequest.Status = REJECTED
		escrow.Released = true
		cost := Costs{
			DocType:               COST,
			ID:                    costId,
			ProviderReimbursement: 0,
			ConsumerRefund:        (escrow.ConsumerDeposit + escrow.ConsumerPayment),
			EscrowID:              escrowID,
			DataConsumer:          escrow.Consumer,
			DataProvider:          escrow.Provider,
			OfferRequestID:        escrow.OfferRequestID,
		}
		costAsBytes, err := json.Marshal(cost)
		if err != nil {
			response["message"] = "Failed to marshal Cost"
			return response, err
		}
		err = ctx.GetStub().PutState(costId, costAsBytes)
		if err != nil {
			response["message"] = "Failed to Putstate Cost"
			return response, err
		}
	}

	offerRequestAsBytes, err = json.Marshal(offerRequest)
	err = ctx.GetStub().PutState(offerRequest.OfferRequestID, offerRequestAsBytes)
	if err != nil {
		return response, err
	}

	escrowAsBytes, err = json.Marshal(escrow)
	err = ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (pm *dataManagement) GetOfferRequestByID(ctx contractapi.TransactionContextInterface, key string) (*OfferRequest, error) {

	assetJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", key)
	}

	var asset OfferRequest
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// GetOfferRequestByOfferID to retrive offer request by offer id
func (pm *dataManagement) GetOfferRequestByOfferID(ctx contractapi.TransactionContextInterface, offerID string) ([]*OfferRequest, error) {

	queryString := fmt.Sprintf(`{"selector":{"docType":"%s","offer_id":"%s"}}`, OFFER_REQUEST, offerID)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)

	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	var offerRequest []*OfferRequest
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var request OfferRequest
		err = json.Unmarshal(queryResponse.Value, &request)
		if err != nil {
			return nil, err
		}
		logger.Info(request.OfferID)

		offerRequest = append(offerRequest, &request)
	}

	return offerRequest, nil

}

// GetAllEscrow to retrive escrows
func (pm *dataManagement) GetAllEscrow(ctx contractapi.TransactionContextInterface) ([]*Escrow, error) {

	queryString := fmt.Sprintf(`{"selector":{"%s": "%s"}}`, DOC_TYPE, ESCROW)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var offerRequest []*Escrow
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var request Escrow
		err = json.Unmarshal(queryResponse.Value, &request)
		if err != nil {
			return nil, err
		}
		offerRequest = append(offerRequest, &request)
	}

	return offerRequest, nil

}

func (pm *dataManagement) GetEscrowByID(ctx contractapi.TransactionContextInterface, dataProvider, dataConsumer, operator string) ([]*Escrow, error) {
	if len(operator) == 0 {
		operator = "$or"
		logger.Info("No Operator found.... Defaulting to $OR operator")
	}

	queryString := fmt.Sprintf(`{"selector":{"%s":"%s","%s":[{"consumer":"%s"},{"producer":"%s"}]}}`, DOC_TYPE, ESCROW, operator, dataConsumer, dataProvider)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var offerRequest []*Escrow
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var request Escrow
		err = json.Unmarshal(queryResponse.Value, &request)
		if err != nil {
			return nil, err
		}
		offerRequest = append(offerRequest, &request)
	}

	return offerRequest, nil

}

func (pm *dataManagement) OldInsertTestDataHash(ctx contractapi.TransactionContextInterface, offerID, hashID, dataHash, filename, entrydate string, offerDataHashID string) error {
	var offer DataOffer

	offerAsBytes, err := ctx.GetStub().GetState(offerID)
	if err != nil || len(offerAsBytes) == 0 {
		return fmt.Errorf("No Such Offer Exists")
	}

	err = json.Unmarshal(offerAsBytes, &offer)
	if err != nil {
		return err
	}

	data := DataHash{
		DocType:   DATA_HASH_VALUE,
		ID:        hashID,
		Hash:      dataHash,
		FileName:  filename,
		EntryDate: entrydate,
	}

	// Create or fetch existing offer data hash
	var offerDataHash OfferDataHash
	offerDataHashAsBytes, err := ctx.GetStub().GetState(offerDataHashID)
	if err != nil || len(offerDataHashAsBytes) == 0 {
		// Create a new offer data hash structure
		offerDataHash = OfferDataHash{
			ID:           offerDataHashID,
			DataHashes:   []DataHash{data},
			OfferID:      offerID,
			DocType:      DATA_HASH,
			DataProvider: offer.Creator,
		}
	} else {
		// Fetch existing offer data hash and append the new data
		err = json.Unmarshal(offerDataHashAsBytes, &offerDataHash)
		if err != nil {
			return err
		}
		offerDataHash.DataHashes = append(offerDataHash.DataHashes, data)
	}

	// Update offer data hash on the ledger
	err = ctx.GetStub().PutState(offerDataHash.ID, offerDataHashAsBytes)
	if err != nil {
		return err
	}

	// Fetch related active data agreements and link the new data hash
	query := fmt.Sprintf(`{"selector":{"%s":"%s","offer_id":"%s","state": %t}}`, DOC_TYPE, DATA_AGREEMENT, offerID, true)
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return err
		}

		var agreement DataAgreement
		err = json.Unmarshal(queryResponse.Value, &agreement)
		if err != nil {
			return err
		}

		// Link the new data hash with the agreement
		agreement.OfferDataHashID = append(agreement.OfferDataHashID, hashID)

		agreementAsBytes, err := json.Marshal(agreement)
		if err != nil {
			return err
		}

		// Update the agreement on the ledger
		err = ctx.GetStub().PutState(agreement.ID, agreementAsBytes)
		if err != nil {
			return err
		}
	}

	return nil
}
func (pm *dataManagement) InsertTestDataHash(ctx contractapi.TransactionContextInterface, offerID, hashID, dataHash, filename, entrydate string, offerDataHashID string) error {
	var offer HistoricalDataOffer

	offerAsBytes, _ := ctx.GetStub().GetState(offerID)
	_ = json.Unmarshal(offerAsBytes, &offer)

	data := DataHash{
		DocType:   DATA_HASH,
		ID:        hashID,
		Hash:      dataHash,
		FileName:  filename,
		EntryDate: entrydate,
	}

	var offerDataHash OfferDataHash
	offerDataHashAsBytes, err := ctx.GetStub().GetState(offerDataHashID)
	if err != nil || len(offerDataHashAsBytes) == 0 {
		offerDataHash = OfferDataHash{
			ID:           offerDataHashID,
			DataHashes:   []DataHash{data},
			OfferID:      offerID,
			DocType:      DATA_HASH,
			DataProvider: offer.Creator,
		}
	} else {
		_ = json.Unmarshal(offerDataHashAsBytes, &offerDataHash)
		offerDataHash.DataHashes = append(offerDataHash.DataHashes, data)
	}

	offerDataHashAsBytes, _ = json.Marshal(offerDataHash)
	err = ctx.GetStub().PutState(offerDataHash.ID, offerDataHashAsBytes)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`{"selector":{"%s":"%s","offer_id":"%s","state": %t}}`, DOC_TYPE, DATA_AGREEMENT, offerID, true)
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return err
		}

		var agreement DataAgreement
		_ = json.Unmarshal(queryResponse.Value, &agreement)

		agreement.OfferDataHashID = append(agreement.OfferDataHashID, hashID)

		agreementAsBytes, _ := json.Marshal(agreement)
		err = ctx.GetStub().PutState(agreement.ID, agreementAsBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

// InsertDataHash to upload data hash to the blockchain

func (pm *dataManagement) InsertDataHash(ctx contractapi.TransactionContextInterface, offerID, hashID, dataHash, filename, entrydate string, offerDataHashID string) error {

	offerAsBytes, err := ctx.GetStub().GetState(offerID)
	if err != nil {
		return err
	}
	if len(offerAsBytes) == 0 {
		return fmt.Errorf("No Such Offer Exists")
	}

	query := fmt.Sprintf(`{"selector":{"%s":"%s","offer_id":"%s"}}`, DOC_TYPE, DATA_HASH, offerID)

	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")
	logger.Info(fmt.Sprintf("Today: %s", today))
	query2 := fmt.Sprintf(`{"selector":{"%s":"%s","offer_id":"%s","state": %t}}`, DOC_TYPE, DATA_AGREEMENT, offerID, true)

	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return err
	}
	defer resultsIterator.Close()
	logger.Info(fmt.Sprintf("Query: %s", query))

	var offerDataHashes []*OfferDataHash
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return err
		}

		var hash OfferDataHash
		err = json.Unmarshal(queryResponse.Value, &hash)
		if err != nil {
			return err
		}
		offerDataHashes = append(offerDataHashes, &hash)
	}
	logger.Info(query)
	logger.Info(fmt.Sprintf("OfferDataHashes: %v", offerDataHashes))
	logger.Info(fmt.Sprintf("OfferDataHashes Length: %d", len(offerDataHashes)))
	logger.Info(offerDataHashes)
	resultsIterator2, err := ctx.GetStub().GetQueryResult(query2)
	if err != nil {
		return err
	}
	defer resultsIterator2.Close()

	var dataAgreement []*DataAgreement
	for resultsIterator2.HasNext() {
		queryResponse, err := resultsIterator2.Next()
		if err != nil {
			return err
		}

		var agreement DataAgreement
		err = json.Unmarshal(queryResponse.Value, &agreement)
		if err != nil {
			return err
		}
		dataAgreement = append(dataAgreement, &agreement)
	}
	logger.Info(query2)
	logger.Info(fmt.Sprintf("DataAgreement: %v", dataAgreement))
	logger.Info(fmt.Sprintf("DataAgreement Length: %d", len(dataAgreement)))
	logger.Info(dataAgreement)
	data := DataHash{
		DocType:   DATA_HASH_VALUE,
		ID:        hashID,
		Hash:      dataHash,
		FileName:  filename,
		EntryDate: entrydate,
	}
	if len(offerDataHashes) == 0 {
		// TODO Create new Offer datahash
		var offer DataOffer
		err := json.Unmarshal(offerAsBytes, &offer)
		if err != nil {
			return err
		}
		offerDataHash := OfferDataHash{
			ID:           offerDataHashID,
			DataHashes:   []DataHash{},
			OfferID:      offerID,
			DocType:      DATA_HASH,
			DataProvider: offer.Creator,
		}

		offerDataHash.DataHashes = append(offerDataHash.DataHashes, data)
		offerDataHashAsBytes, err := json.Marshal(offerDataHash)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(offerDataHash.ID, offerDataHashAsBytes)
		if err != nil {
			return err
		}
	} else {

		offerDataHashes[0].DataHashes = append(offerDataHashes[0].DataHashes, data)
		logger.Info(fmt.Sprintf("OfferDataHashes: %v", offerDataHashes))

		offerDataHashAsBytes, err := json.Marshal(offerDataHashes[0])
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(offerDataHashes[0].ID, offerDataHashAsBytes)
		if err != nil {
			return err
		}
	}

	for _, _agreement := range dataAgreement {
		logger.Info(fmt.Sprintf("Agreement: %v", _agreement))
		logger.Info(fmt.Sprintf("Agreement Length: %d", len(dataAgreement)))
		_agreement.OfferDataHashID = append(_agreement.OfferDataHashID, hashID)
		_aggrementAsBytes, err := json.Marshal(_agreement)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(_agreement.ID, _aggrementAsBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pm *dataManagement) InsertTestHistoricalDataHash(ctx contractapi.TransactionContextInterface, offerID, hashID, dataHash, filename, entrydate string, offerDataHashID string) error {
	offerAsBytes, _ := ctx.GetStub().GetState(offerID)

	var offer HistoricalDataOffer
	_ = json.Unmarshal(offerAsBytes, &offer)

	data := DataHash{
		DocType:   DATA_HASH_VALUE,
		ID:        hashID,
		Hash:      dataHash,
		FileName:  filename,
		EntryDate: entrydate,
	}

	var offerDataHash OfferDataHash
	offerDataHashAsBytes, err := ctx.GetStub().GetState(offerDataHashID)
	if err != nil || len(offerDataHashAsBytes) == 0 {
		offerDataHash = OfferDataHash{
			ID:           offerDataHashID,
			DataHashes:   []DataHash{data},
			OfferID:      offerID,
			DocType:      DATA_HASH,
			DataProvider: offer.Creator,
		}
	} else {
		err = json.Unmarshal(offerDataHashAsBytes, &offerDataHash)
		if err != nil {
			return err
		}
		offerDataHash.DataHashes = append(offerDataHash.DataHashes, data)
	}

	offerDataHashAsBytes, _ = json.Marshal(offerDataHash)
	err = ctx.GetStub().PutState(offerDataHash.ID, offerDataHashAsBytes)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`{"selector":{"%s":"%s","offer_id":"%s","state": %t}}`, DOC_TYPE, DATA_AGREEMENT, offerID, true)
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return err
		}

		var agreement DataAgreement
		err = json.Unmarshal(queryResponse.Value, &agreement)
		if err != nil {
			return err
		}

		agreement.OfferDataHashID = append(agreement.OfferDataHashID, hashID)

		agreementAsBytes, err := json.Marshal(agreement)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(agreement.ID, agreementAsBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

// func (pm *dataManagement) InsertTestHistoricalDataHash(ctx contractapi.TransactionContextInterface, offerID, hashID, dataHash, filename, entrydate string, offerDataHashID string) error {
// 	offerAsBytes, _ := ctx.GetStub().GetState(offerID)

// 	query := fmt.Sprintf(`{"selector":{"%s":"%s","offer_id":"%s"}}`, DOC_TYPE, DATA_HASH, offerID)
// 	resultsIterator, _ := ctx.GetStub().GetQueryResult(query)
// 	defer resultsIterator.Close()

// 	var offerDataHashes []*OfferDataHash
// 	for resultsIterator.HasNext() {
// 		queryResponse, _ := resultsIterator.Next()
// 		var hash OfferDataHash
// 		_ = json.Unmarshal(queryResponse.Value, &hash)
// 		offerDataHashes = append(offerDataHashes, &hash)
// 	}

// 	resultsIterator2, _ := ctx.GetStub().GetQueryResult(query)
// 	defer resultsIterator2.Close()

// 	var dataAgreement []*DataAgreement
// 	for resultsIterator2.HasNext() {
// 		queryResponse, _ := resultsIterator2.Next()
// 		var agreement DataAgreement
// 		_ = json.Unmarshal(queryResponse.Value, &agreement)
// 		dataAgreement = append(dataAgreement, &agreement)
// 	}

// 	data := DataHash{
// 		DocType:   DATA_HASH_VALUE,
// 		ID:        offerDataHashID,
// 		Hash:      dataHash,
// 		FileName:  filename,
// 		EntryDate: entrydate,
// 	}

// 	if len(offerDataHashes) == 0 {
// 		var offer HistoricalDataOffer
// 		_ = json.Unmarshal(offerAsBytes, &offer)
// 		offerDataHash := OfferDataHash{
// 			ID:           hashID,
// 			DataHashes:   []DataHash{data}, // Append the new data hash here
// 			OfferID:      offerID,
// 			DocType:      DATA_HASH,
// 			DataProvider: offer.Creator,
// 		}
// 		offerDataHash.DataHashes = append(offerDataHash.DataHashes, data)

// 		offerDataHashAsBytes, _ := json.Marshal(offerDataHash)
// 		_ = ctx.GetStub().PutState(offerDataHash.ID, offerDataHashAsBytes)
// 	} else {
// 		offerDataHashes[0].DataHashes = append(offerDataHashes[0].DataHashes, data)

// 		offerDataHashAsBytes, _ := json.Marshal(offerDataHashes[0])
// 		_ = ctx.GetStub().PutState(offerDataHashes[0].ID, offerDataHashAsBytes)
// 	}

// 	for _, _agreement := range dataAgreement {
// 		_agreement.OfferDataHashID = append(_agreement.OfferDataHashID, offerDataHashID)
// 		_aggrementAsBytes, _ := json.Marshal(_agreement)
// 		_ = ctx.GetStub().PutState(_agreement.ID, _aggrementAsBytes)
// 	}

// 	return nil
// }

func (pm *dataManagement) InsertHistoricalDataHash(ctx contractapi.TransactionContextInterface, offerID, hashID, dataHash, filename, entrydate string, offerDataHashID string) error {

	offerAsBytes, err := ctx.GetStub().GetState(offerID)

	logger.Info(fmt.Sprintf("OfferAsBytes: %s", offerAsBytes))
	logger.Info(fmt.Sprintf("err: %s", err))
	if err != nil {
		return err
	}
	if len(offerAsBytes) == 0 {
		return fmt.Errorf("No Such Offer Exists")
	}

	query := fmt.Sprintf(`{"selector":{"%s":"%s","offer_id":"%s"}}`, DOC_TYPE, DATA_HASH, offerID)
	logger.Info(query)

	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")
	logger.Info(fmt.Sprintf("Today: %s", today))
	query2 := fmt.Sprintf(`{"selector":{"%s":"%s","offer_id":"%s","state": %t}}`, DOC_TYPE, DATA_AGREEMENT, offerID, true)
	logger.Info(query2)
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	logger.Info(fmt.Sprintf("resultsIterator: %s", resultsIterator))
	logger.Info(fmt.Sprintf("err: %s", err))
	if err != nil {
		return err
	}
	defer resultsIterator.Close()

	var offerDataHashes []*OfferDataHash
	logger.Info(fmt.Sprintf("OfferDataHashes: %v", offerDataHashes))
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		logger.Info(fmt.Sprintf("queryResponse: %s", queryResponse))
		if err != nil {
			return err
		}

		var hash OfferDataHash
		logger.Info(fmt.Sprintf("hash: %v", hash))
		logger.Info(fmt.Sprintf("queryResponse.Value: %s", queryResponse.Value))
		logger.Info(fmt.Sprintf("offerDataHashes: %v", offerDataHashes))
		err = json.Unmarshal(queryResponse.Value, &hash)
		if err != nil {
			return err
		}
		offerDataHashes = append(offerDataHashes, &hash)
	}
	logger.Info(query)
	logger.Info(fmt.Sprintf("offerDataHashes1: %v", offerDataHashes))
	resultsIterator2, err := ctx.GetStub().GetQueryResult(query2)
	if err != nil {
		return err
	}
	defer resultsIterator2.Close()
	logger.Info(fmt.Sprintf("resultsIterator2: %s", resultsIterator2))
	var dataAgreement []*DataAgreement
	logger.Info(fmt.Sprintf("dataAgreement: %v", dataAgreement))
	for resultsIterator2.HasNext() {
		queryResponse, err := resultsIterator2.Next()
		if err != nil {
			return err
		}

		var agreement DataAgreement
		logger.Info(fmt.Sprintf("agreement: %v", agreement))
		logger.Info(fmt.Sprintf("queryResponse.Value: %s", queryResponse.Value))
		logger.Info(fmt.Sprintf("dataAgreement: %v", dataAgreement))
		err = json.Unmarshal(queryResponse.Value, &agreement)
		if err != nil {
			return err
		}
		dataAgreement = append(dataAgreement, &agreement)
	}
	logger.Info(query2)
	logger.Info(fmt.Sprintf("dataAgreement1: %v", dataAgreement))
	data := DataHash{
		DocType:   DATA_HASH_VALUE,
		ID:        offerDataHashID,
		Hash:      dataHash,
		FileName:  filename,
		EntryDate: entrydate,
	}
	logger.Info(fmt.Sprintf("data: %v", data))
	if len(offerDataHashes) == 0 {
		// TODO Create new Offer datahash
		var offer HistoricalDataOffer
		logger.Info(fmt.Sprintf("offer: %v", offer))
		logger.Info(fmt.Sprintf("offerAsBytes: %s", offerAsBytes))
		logger.Info(fmt.Sprintf("historicalDataOffer: %v", offer))
		err := json.Unmarshal(offerAsBytes, &offer)
		if err != nil {
			return err
		}
		offerDataHash := OfferDataHash{
			ID:           hashID,
			DataHashes:   []DataHash{data}, // Append the new data hash here
			OfferID:      offerID,
			DocType:      DATA_HASH,
			DataProvider: offer.Creator,
		}
		logger.Info(fmt.Sprintf("offerDataHash: %v", offerDataHash))

		offerDataHash.DataHashes = append(offerDataHash.DataHashes, data)
		logger.Info(fmt.Sprintf("offerDataHash: %v", offerDataHash))
		offerDataHashAsBytes, err := json.Marshal(offerDataHash)
		if err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("offerDataHashAsBytes: %v", offerDataHashAsBytes))
		err = ctx.GetStub().PutState(offerDataHash.ID, offerDataHashAsBytes)
		logger.Info(fmt.Sprintf("err: %v", err))
		if err != nil {
			return err
		}
	} else {

		offerDataHashes[0].DataHashes = append(offerDataHashes[0].DataHashes, data)
		logger.Info(fmt.Sprintf("offerDataHashes: %v", offerDataHashes))
		logger.Info(fmt.Sprintf("offerDataHashes[0]: %v", offerDataHashes[0]))
		logger.Info(fmt.Sprintf("data: %v", data))

		offerDataHashAsBytes, err := json.Marshal(offerDataHashes[0])
		if err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("offerDataHashAsBytes: %v", offerDataHashAsBytes))
		logger.Info(fmt.Sprintf("offerDataHashes[0].ID: %v", offerDataHashes[0].ID))
		err = ctx.GetStub().PutState(offerDataHashes[0].ID, offerDataHashAsBytes)
		if err != nil {
			return err
		}
	}

	for _, _agreement := range dataAgreement {
		logger.Info(fmt.Sprintf("_agreement: %v", _agreement))
		logger.Info(fmt.Sprintf("_agreement.OfferDataHashID: %v", _agreement.OfferDataHashID))
		logger.Info(fmt.Sprintf("hashID: %v", hashID))
		_agreement.OfferDataHashID = append(_agreement.OfferDataHashID, offerDataHashID)
		logger.Info(fmt.Sprintf("_agreement.OfferDataHashID: %v", _agreement.OfferDataHashID))
		_aggrementAsBytes, err := json.Marshal(_agreement)
		logger.Info(fmt.Sprintf("_aggrementAsBytes: %v", _aggrementAsBytes))
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(_agreement.ID, _aggrementAsBytes)
		if err != nil {
			return err
		}
	}
	logger.Info("Data Hash Inserted Successfully")

	return nil
}

func (pm *dataManagement) GetAllDataHashes(ctx contractapi.TransactionContextInterface) ([]*OfferDataHash, error) {

	queryString := fmt.Sprintf(`{"selector":{"%s": "%s"}}`, DOC_TYPE, DATA_HASH_VALUE)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var offerRequest []*OfferDataHash
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var request OfferDataHash
		err = json.Unmarshal(queryResponse.Value, &request)
		if err != nil {
			return nil, err
		}
		offerRequest = append(offerRequest, &request)
	}

	return offerRequest, nil

}

// GetDataHashByOfferID to retrive the hashs of spsific offer
func (pm *dataManagement) GetDataHashByOfferID(ctx contractapi.TransactionContextInterface, id string, provider string) ([]*OfferDataHash, error) {

	query := fmt.Sprintf(`{"selector":{"%s":"%s","offer_id":"%s","data_provider":"%s"}}`, DOC_TYPE, DATA_HASH, id, provider)
	logger.Info(query)
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var offerDataHashes []*OfferDataHash
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var hash OfferDataHash
		err = json.Unmarshal(queryResponse.Value, &hash)
		if err != nil {
			return nil, err
		}
		offerDataHashes = append(offerDataHashes, &hash)
	}
	logger.Info(len(offerDataHashes))
	return offerDataHashes, nil

}

// GetDataHashByAgreementID to retrive the hashs of spsific agreement
func (pm *dataManagement) GetDataHashByAgreementID(ctx contractapi.TransactionContextInterface, agreementID string) (*AgreementHash, error) {
	logger.Info("GetDataHashByAgreementID")
	var agreement DataAgreement
	agreementAsBytes, err := ctx.GetStub().GetState(agreementID)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(agreementAsBytes, &agreement)
	if err != nil {
		return nil, err
	}
	if len(agreement.OfferDataHashID) == 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`{"selector":{"%s":"%s","offer_id": "%s" }}`, DOC_TYPE, DATA_HASH, agreement.OfferID)
	logger.Info(query)
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var dataHashes []*OfferDataHash
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var hash OfferDataHash
		err = json.Unmarshal(queryResponse.Value, &hash)
		if err != nil {
			return nil, err
		}
		dataHashes = append(dataHashes, &hash)
	}
	result := &AgreementHash{
		Hashes:    dataHashes,
		Agreement: &agreement,
	}
	return result, nil

}

func (pm *dataManagement) GetDataHashesByAgreementID(ctx contractapi.TransactionContextInterface, agreementID string, multipleOfferID []string) (*AgreementHash, error) {
	logger.Info("GetDataHashesByAgreementID")
	var agreement DataAgreement
	agreementAsBytes, err := ctx.GetStub().GetState(agreementID)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(agreementAsBytes, &agreement)
	if err != nil {
		return nil, err
	}

	// Initialize a slice to store the data hashes for each offer ID
	var dataHashes []*OfferDataHash

	// Loop through each offer ID in the multipleOfferID array and query the data hashes
	for _, offerID := range multipleOfferID {
		query := fmt.Sprintf(`{"selector":{"%s":"%s","offer_id": "%s" }}`, DOC_TYPE, DATA_HASH, offerID)
		logger.Info(query)
		resultsIterator, err := ctx.GetStub().GetQueryResult(query)
		if err != nil {
			return nil, err
		}
		defer resultsIterator.Close()

		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return nil, err
			}

			var hash OfferDataHash
			err = json.Unmarshal(queryResponse.Value, &hash)
			if err != nil {
				return nil, err
			}
			dataHashes = append(dataHashes, &hash)
		}
	}

	result := &AgreementHash{
		Hashes:    dataHashes,
		Agreement: &agreement,
	}
	return result, nil
}

// GetAllAgreements to get all agreemets
func (pm *dataManagement) GetAllAgreements(ctx contractapi.TransactionContextInterface, dataProvider, dataConsumer, operator string) ([]*DataAgreement, error) {

	if len(operator) == 0 {
		operator = "$or"
		logger.Info("No Operator found.... Defaulting to $OR operator")
	}
	// timestamp, _ := ctx.GetStub().GetTxTimestamp()
	// currentTime := time.Unix(timestamp.GetSeconds(), 0)
	// today := currentTime.Format("2006-01-02 15:04")
	// queryString := fmt.Sprintf(`{"selector":{"%s":"%s","%s":[{"dataConsumer":"%s"},{"dataProvider":"%s"}],"state": %t,"end_date":{"$gte": "%s"}}}`, DOC_TYPE, DATA_AGREEMENT, operator, dataConsumer, dataProvider, true, today)

	queryString := fmt.Sprintf(`
	{"selector":{
		"%s":"%s",
		"%s":[
			{"dataConsumer":"%s"},{"dataProvider":"%s"}
		   ]
		}	
	}`, DOC_TYPE, DATA_AGREEMENT, operator, dataConsumer, dataProvider)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var offerRequest []*DataAgreement
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var request DataAgreement
		err = json.Unmarshal(queryResponse.Value, &request)
		if err != nil {
			return nil, err
		}
		offerRequest = append(offerRequest, &request)
	}

	return offerRequest, nil
}

// get Aggremenets by Agreement ID
func (pm *dataManagement) GetAgreementByID(ctx contractapi.TransactionContextInterface, agreementID string) (*DataAgreement, error) {
	logger.Info("GetAgreementByID")
	var agreement DataAgreement
	agreementAsBytes, err := ctx.GetStub().GetState(agreementID)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(agreementAsBytes, &agreement)
	if err != nil {
		return nil, err
	}
	return &agreement, nil
}

func (pm *dataManagement) GetAllCost(ctx contractapi.TransactionContextInterface, dataProvider, dataConsumer, operator string) ([]*OfferRequest, error) {

	if len(operator) == 0 {
		operator = "$or"
		logger.Info("No Operator found.... Defaulting to $OR operator")
	}
	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")
	queryString := fmt.Sprintf(`{"selector":{"%s":"%s","%s":[{"dataConsumer":"%s"},{"dataProvider":"%s"}], "endDate":{"$lte":"%s"}}}`, DOC_TYPE, OFFER_REQUEST, operator, dataConsumer, dataProvider, today)
	logger.Info(queryString)
	fmt.Println(queryString)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var offerRequest []*OfferRequest
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var request OfferRequest
		err = json.Unmarshal(queryResponse.Value, &request)
		if err != nil {
			return nil, err
		}
		offerRequest = append(offerRequest, &request)
	}

	return offerRequest, nil
}

// this function is deprecated
func (pm *dataManagement) GetCostFromEscrow(ctx contractapi.TransactionContextInterface, dataProvider, dataConsumer, operator string) ([]*Escrow, error) {
	if len(operator) == 0 {
		operator = "$or"
		logger.Info("No Operator found.... Defaulting to $OR operator")
	}
	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")
	queryString := fmt.Sprintf(`{"selector":{"%s":"%s","%s":[{"provider":"%s"},{"consumer":"%s"}], "endDate":{"$lte":"%s"}}}`, DOC_TYPE, ESCROW, operator, dataConsumer, dataProvider, today)
	logger.Info(queryString)
	fmt.Println(queryString)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var escrows []*Escrow
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var escrow Escrow
		err = json.Unmarshal(queryResponse.Value, &escrow)
		if err != nil {
			return nil, err
		}
		escrows = append(escrows, &escrow)
	}

	return escrows, nil
}

func (pm *dataManagement) GetTotalCost(ctx contractapi.TransactionContextInterface, dataProvider, dataConsumer, operator string) ([]*Costs, error) {

	queryString := fmt.Sprintf(`{"selector":{"%s":"%s","%s":[{"dataProvider":"%s"},{"dataConsumer":"%s"}]}}`, DOC_TYPE, COST, operator, dataProvider, dataConsumer)
	logger.Info(queryString)
	fmt.Println(queryString)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var costs []*Costs
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var cost Costs
		err = json.Unmarshal(queryResponse.Value, &cost)
		if err != nil {
			return nil, err
		}
		costs = append(costs, &cost)
	}

	return costs, nil

}

func (pm *dataManagement) GetAllCosts(ctx contractapi.TransactionContextInterface) ([]*Costs, error) {
	queryString := fmt.Sprintf(`{"selector":{"%s":"%s"}}`, DOC_TYPE, COST)
	logger.Info(queryString)
	fmt.Println(queryString)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var costs []*Costs
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var cost Costs
		err = json.Unmarshal(queryResponse.Value, &cost)
		if err != nil {
			return nil, err
		}
		costs = append(costs, &cost)
	}

	return costs, nil
}

// RevokeAgreement revokes an agreement ,update escrow record, and find costs
func (pm *dataManagement) RevokeAgreement(ctx contractapi.TransactionContextInterface, agreementId string, isProvider bool) error {

	var agreement DataAgreement

	exists, err := pm.ObjectExists(ctx, agreementId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("No Such agreement exist")
	}

	agreementAsBytes, err := ctx.GetStub().GetState(agreementId)
	if err != nil {
		return err
	}
	err = json.Unmarshal(agreementAsBytes, &agreement)
	if err != nil {
		return err
	}
	agreement.State = false
	var escrow Escrow

	escrowAsBytes, err := ctx.GetStub().GetState(agreement.EscrowID)
	if err != nil {
		return err
	}
	err = json.Unmarshal(escrowAsBytes, &escrow)
	if err != nil {
		return err
	}
	// update the escrow record
	escrow.Status = REVOKED
	escrow.Released = true

	if err != nil {
		return err
	}

	agreementAsBytes, err = json.Marshal(agreement)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(agreement.ID, agreementAsBytes)
	if err != nil {
		return err
	}

	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")
	cost := Costs{
		CreatedAt:             today,
		DocType:               COST,
		ID:                    rotateLeft(escrow.ID, 4),
		Agreement:             escrow.AgreementID,
		ProviderReimbursement: escrow.ProviderDeposit,
		ConsumerRefund:        escrow.ConsumerDeposit,
		EscrowID:              escrow.ID,
		DataConsumer:          escrow.Consumer,
		DataProvider:          escrow.Provider,
		OfferRequestID:        escrow.OfferRequestID,
	}
	// in case the end date is reached already
	if escrow.Status == ACTIVE && escrow.EndDate == today {

		escrow.Status = EXPIRED

		cost.ProviderReimbursement = roundOff((escrow.ProviderDeposit + escrow.ConsumerPayment))
		cost.ConsumerRefund = escrow.ConsumerDeposit
		ctx.GetStub().PutState(escrow.ID, escrowAsBytes)

	} else { // in case the end date is not reached yet
		var offerRequest OfferRequest
		offerRequestAsBytes, _ := ctx.GetStub().GetState(escrow.OfferRequestID)

		json.Unmarshal(offerRequestAsBytes, &offerRequest)
		hours := calculateHours(escrow.StartDate, escrow.EndDate)
		logger.Info("hour %f", hours)
		pricePerHour := math.Round((offerRequest.Price/hours)*100) / 100
		logger.Info("pricePerHour %f", pricePerHour)
		revokedHours := calculateHours(escrow.StartDate, today)
		logger.Info("revokedHours %f", revokedHours)

		cost.ProviderReimbursement = roundOff(pricePerHour*revokedHours) + escrow.ProviderDeposit
		logger.Info("cost.ProviderReimbursement %f", cost.ProviderReimbursement)
		remainingHours := hours - revokedHours
		cost.ConsumerRefund = roundOff(remainingHours*pricePerHour) + escrow.ConsumerDeposit
		logger.Info("cost.ConsumerRefund %f", cost.ConsumerRefund)
		if isProvider {

			cost.ConsumerRefund = escrow.ProviderDeposit + roundOff(remainingHours*pricePerHour) + escrow.ConsumerDeposit
			cost.ProviderReimbursement = roundOff(revokedHours * pricePerHour)

		} else {
			cost.ConsumerRefund = roundOff(remainingHours*pricePerHour) + escrow.ConsumerDeposit
			cost.ProviderReimbursement = escrow.ProviderDeposit + roundOff(revokedHours*pricePerHour)
		}
	}

	escrowAsBytes, _ = json.Marshal(escrow)
	ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
	costAsBytes, _ := json.Marshal(cost)

	ctx.GetStub().PutState(cost.ID, costAsBytes)

	return nil
}

func (pm *dataManagement) RevokeAgreementNew(ctx contractapi.TransactionContextInterface, agreementId string, isProvider bool) error {
	var agreement DataAgreement

	exists, err := pm.ObjectExists(ctx, agreementId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("No Such agreement exist")
	}

	agreementAsBytes, err := ctx.GetStub().GetState(agreementId)
	if err != nil {
		return err
	}
	err = json.Unmarshal(agreementAsBytes, &agreement)
	if err != nil {
		return err
	}
	agreement.State = false
	var escrow Escrow

	escrowAsBytes, err := ctx.GetStub().GetState(agreement.EscrowID)
	if err != nil {
		return err
	}
	err = json.Unmarshal(escrowAsBytes, &escrow)
	if err != nil {
		return err
	}
	// update the escrow record
	escrow.Status = REVOKED
	escrow.Released = true

	if err != nil {
		return err
	}

	agreementAsBytes, err = json.Marshal(agreement)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(agreement.ID, agreementAsBytes)
	if err != nil {
		return err
	}

	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")
	cost := Costs{
		CreatedAt:             today,
		DocType:               COST,
		ID:                    rotateLeft(escrow.ID, 4),
		Agreement:             escrow.AgreementID,
		ProviderReimbursement: escrow.ProviderDeposit,
		ConsumerRefund:        escrow.ConsumerDeposit,
		EscrowID:              escrow.ID,
		DataConsumer:          escrow.Consumer,
		DataProvider:          escrow.Provider,
		OfferRequestID:        escrow.OfferRequestID,
	}
	// Calculate remaining hours until the end date
	remainingHours := calculateHours(today, escrow.EndDate)

	if remainingHours < 24 {
		if isProvider {
			// If less than 24 hours are remaining and provider revokes the agreement
			// Transfer the provider deposit to the consumer
			// escrow.ConsumerDeposit += escrow.ProviderDeposit
			// escrow.ProviderDeposit = 0
			cost.ConsumerRefund = escrow.ProviderDeposit + escrow.ConsumerDeposit + escrow.ConsumerPayment
			cost.ProviderReimbursement = 0

		} else {
			cost.ConsumerRefund = escrow.ConsumerPayment
			cost.ProviderReimbursement = escrow.ProviderDeposit + escrow.ConsumerDeposit

			// If less than 24 hours are remaining and consumer revokes the agreement
			// Transfer the consumer deposit to the provider
			// escrow.ProviderDeposit += escrow.ConsumerDeposit
			// escrow.ConsumerDeposit = 0
		}
	}
	if remainingHours > 24 {
		if isProvider {
			// If less than 24 hours are remaining and provider revokes the agreement
			// Transfer the provider deposit to the consumer
			// escrow.ConsumerDeposit += escrow.ProviderDeposit
			// escrow.ProviderDeposit = 0
			cost.ConsumerRefund = escrow.ConsumerDeposit + escrow.ConsumerPayment
			cost.ProviderReimbursement = escrow.ProviderDeposit

		} else {
			cost.ConsumerRefund = escrow.ConsumerPayment + escrow.ConsumerDeposit
			cost.ProviderReimbursement = escrow.ProviderDeposit

			// If less than 24 hours are remaining and consumer revokes the agreement
			// Transfer the consumer deposit to the provider
			// escrow.ProviderDeposit += escrow.ConsumerDeposit
			// escrow.ConsumerDeposit = 0
		}
	}

	/////////////////////////////
	// cost.ProviderReimbursement = roundOff(pricePerHour*revokedHours) + escrow.ProviderDeposit
	// 	logger.Info("cost.ProviderReimbursement %f", cost.ProviderReimbursement)
	// 	remainingHours := hours - revokedHours
	// 	cost.ConsumerRefund = roundOff(remainingHours*pricePerHour) + escrow.ConsumerDeposit
	// 	logger.Info("cost.ConsumerRefund %f", cost.ConsumerRefund)
	// 	if isProvider {

	// 		cost.ConsumerRefund = escrow.ProviderDeposit + escrow.ConsumerDeposit + escrow.ConsumerPayment
	// 		cost.ProviderReimbursement = 0

	// 	} else {
	// 		cost.ConsumerRefund = escrow.ConsumerPayment
	// 		cost.ProviderReimbursement = escrow.ProviderDeposit + escrow.ConsumerDeposit
	// 	}

	/////////////////////////////

	escrowAsBytes, _ = json.Marshal(escrow)
	ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
	costAsBytes, _ := json.Marshal(cost)

	ctx.GetStub().PutState(cost.ID, costAsBytes)

	return nil
}

// ReleaseEscrow will be invoked when the agreement is expired (scheduled job) from the agenda
func (pm *dataManagement) ReleaseEscrow(ctx contractapi.TransactionContextInterface, escrowId string, releaseCostID string) error {
	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")
	costQuery := fmt.Sprintf(`{
		"selector": {
			"docType": "%s",
			"escrow_id":"%s"
		}
	}`, COST, escrowId)

	costExist := pm.checkCostExist(ctx, costQuery)
	if costExist {
		return fmt.Errorf("Cost already exist for selected escrow id %s", escrowId)
	}

	var escrow Escrow
	escrowAsBytes, err := ctx.GetStub().GetState(escrowId)

	if err != nil {
		return err
	}
	err = json.Unmarshal(escrowAsBytes, &escrow)
	if err != nil {
		return err
	}

	cost := Costs{
		CreatedAt:             today,
		DocType:               COST,
		ID:                    releaseCostID,
		Agreement:             escrow.AgreementID,
		ProviderReimbursement: escrow.ProviderDeposit,
		ConsumerRefund:        escrow.ConsumerDeposit,
		EscrowID:              escrowId,
		DataConsumer:          escrow.Consumer,
		DataProvider:          escrow.Provider,
		OfferRequestID:        escrow.OfferRequestID,
	}

	if escrow.Status == ACTIVE && escrow.EndDate == today {
		escrow.Released = true
		escrow.Status = EXPIRED
		escrowAsBytes, err = json.Marshal(escrow)
		if err != nil {
			return err
		}
		cost.ProviderReimbursement = roundOff((escrow.ProviderDeposit + escrow.ConsumerPayment))
		cost.ConsumerRefund = escrow.ConsumerDeposit
		ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
		var aggrement DataAgreement
		aggrementAsBytes, _ := ctx.GetStub().GetState(escrow.AgreementID)
		json.Unmarshal(aggrementAsBytes, &aggrement)
		aggrement.State = false
		aggrementAsBytes, err = json.Marshal(aggrement)
		ctx.GetStub().PutState(aggrement.ID, aggrementAsBytes)
		if err != nil {
			return err
		}
	} else {
		var offerRequest OfferRequest
		offerRequestAsBytes, _ := ctx.GetStub().GetState(escrow.OfferRequestID)

		json.Unmarshal(offerRequestAsBytes, &offerRequest)
		hours := calculateHours(escrow.StartDate, escrow.EndDate)
		logger.Info("hour %f", hours)
		pricePerHour := math.Round((offerRequest.Price/hours)*100) / 100
		logger.Info("pricePerHour %f", pricePerHour)
		revokedHours := calculateHours(escrow.StartDate, today)
		logger.Info("revokedHours %f", revokedHours)
		cost.ProviderReimbursement = roundOff(pricePerHour * revokedHours)
		logger.Info("cost.ProviderReimbursement %f", cost.ProviderReimbursement)
		remainingHours := hours - revokedHours
		cost.ConsumerRefund = roundOff(remainingHours * pricePerHour)
	}

	costAsBytes, err := json.Marshal(cost)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(cost.ID, costAsBytes)
	if err != nil {
		return err
	}
	return nil
}

// func (pm *dataManagement) FalsifyClaim2(ctx contractapi.TransactionContextInterface, offerId, _hashes, agreementId, costId string) (map[string]interface{}, error) {
// 	response := make(map[string]interface{})
// 	response["txId"] = ctx.GetStub().GetTxID()
// 	hashes := strings.Split(_hashes, ",")
// 	costQuery := fmt.Sprintf(`{
// 		"selector": {
// 			"docType": "%s",
// 			"agreement":"%s"
// 		}
// 	}`, COST, agreementId)
// 	// in case the request is sent to manipulate existing calculated costs
// 	costExist := pm.checkCostExist(ctx, costQuery)
// 	if costExist {
// 		response["message"] = fmt.Sprintf("Cost already exists for selected agreement id %s", agreementId)
// 		return response, nil
// 	}

// 	queryString := fmt.Sprintf(`{
// 		"selector": {
// 		   "docType": "%s",
// 		   "offer_id": "%s"
// 		}
// 	 }`, DATA_HASH, offerId)

// 	logger.Info(queryString)

// 	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
// 	if err != nil {
// 		response["message"] = "resultsIterator failed"
// 		return response, err
// 	}
// 	defer resultsIterator.Close()

// 	var offerDataHashes []*OfferDataHash
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			response["message"] = "resultsIterator failed"
// 			return response, err
// 		}

// 		var offerDataHash OfferDataHash
// 		err = json.Unmarshal(queryResponse.Value, &offerDataHash)
// 		if err != nil {
// 			response["message"] = "Failed to Unmarshal"
// 			return response, err
// 		}
// 		offerDataHashes = append(offerDataHashes, &offerDataHash)
// 	}
// 	agreementAsBytes, _ := ctx.GetStub().GetState(agreementId)
// 	var dataAgreement DataAgreement
// 	json.Unmarshal(agreementAsBytes, &dataAgreement)

// 	timestamp, _ := ctx.GetStub().GetTxTimestamp()
// 	currentTime := time.Unix(timestamp.GetSeconds(), 0)
// 	today := currentTime.Format("2006-01-02 15:04")

// 	cost := Costs{
// 		CreatedAt:      today,
// 		DocType:        COST,
// 		ID:             costId,
// 		Agreement:      agreementId,
// 		OfferRequestID: dataAgreement.OfferRequestID,
// 		EscrowID:       dataAgreement.EscrowID,
// 		DataConsumer:   dataAgreement.DataConsumer,
// 		DataProvider:   dataAgreement.DataProvider,
// 	}

// 	var falsify int
// 	for _, hash := range hashes {
// 		for _, aggHash := range offerDataHashes[0].DataHashes {
// 			if hash == aggHash.Hash {
// 				falsify++
// 			}
// 		}
// 	}

// 	falsify = len(hashes) - falsify

// 	cost.FalsifyCount = falsify
// 	if falsify <= 0 {
// 		logger.Info("Wrong falsify claim")
// 		response["message"] = "Wrong falsify claim"

// 		// No hourly penalty, just calculate reimbursement and refund based on deposits
// 		cost.ProviderReimbursement = dataAgreement.ProviderDeposit + dataAgreement.ConsumerDeposit
// 		cost.ConsumerRefund = dataAgreement.Price

// 	} else if falsify > 0 {
// 		logger.Info("Valid falsify claim")
// 		response["message"] = "Valid falsify claim found"
// 		cost.ProviderReimbursement = 0
// 		cost.ConsumerRefund = dataAgreement.ProviderDeposit + dataAgreement.ConsumerDeposit + dataAgreement.Price
// 	}

// 	costAsBytes, _ := json.Marshal(cost)
// 	ctx.GetStub().PutState(cost.ID, costAsBytes)

// 	var escrow Escrow
// 	escrowAsBytes, _ := ctx.GetStub().GetState(dataAgreement.EscrowID)
// 	json.Unmarshal(escrowAsBytes, &escrow)
// 	escrow.Released = true
// 	escrow.Status = REVOKED
// 	escrowAsBytes, _ = json.Marshal(escrow)
// 	ctx.GetStub().PutState(escrow.ID, escrowAsBytes)

// 	dataAgreement.State = false

// 	agreementAsBytes, _ = json.Marshal(dataAgreement)
// 	ctx.GetStub().PutState(dataAgreement.ID, agreementAsBytes)

// 	return response, nil
// }

func (pm *dataManagement) FalsifyClaim2(ctx contractapi.TransactionContextInterface, offerId, _hashes, agreementId, costId string) (map[string]interface{}, error) {
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	hashes := strings.Split(_hashes, ",")
	costQuery := fmt.Sprintf(`{
        "selector": {
            "docType": "%s",
            "agreement":"%s"
        }
    }`, COST, agreementId)
	// in case the request is sent to manipulate existing calculated costs
	costExist := pm.checkCostExist(ctx, costQuery)
	if costExist {
		response["message"] = fmt.Sprintf("Cost already exists for selected agreement id %s", agreementId)
		return response, nil
	}

	queryString := fmt.Sprintf(`{
        "selector": {
           "docType": "%s",
           "offer_id": "%s"
        }
     }`, DATA_HASH, offerId)

	logger.Info(queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		response["message"] = "resultsIterator failed"
		return response, err
	}
	defer resultsIterator.Close()

	var offerDataHashes []*OfferDataHash
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			response["message"] = "resultsIterator failed"
			return response, err
		}

		var offerDataHash OfferDataHash
		err = json.Unmarshal(queryResponse.Value, &offerDataHash)
		if err != nil {
			response["message"] = "Failed to Unmarshal"
			return response, err
		}
		offerDataHashes = append(offerDataHashes, &offerDataHash)
	}
	agreementAsBytes, _ := ctx.GetStub().GetState(agreementId)
	var dataAgreement DataAgreement
	err = json.Unmarshal(agreementAsBytes, &dataAgreement)
	if err != nil {
		return response, err
	}

	var escrow Escrow
	escrowAsBytes, err := ctx.GetStub().GetState(dataAgreement.EscrowID)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(escrowAsBytes, &escrow)
	if err != nil {
		return response, err
	}
	escrow.Released = true
	escrow.Status = REVOKED
	escrowAsBytes, err = json.Marshal(escrow)
	if err != nil {
		return response, err
	}
	err = ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
	if err != nil {
		return response, err
	}

	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")

	cost := Costs{
		CreatedAt:      today,
		DocType:        COST,
		ID:             costId,
		Agreement:      agreementId,
		OfferRequestID: dataAgreement.OfferRequestID,
		EscrowID:       dataAgreement.EscrowID,
		DataConsumer:   dataAgreement.DataConsumer,
		DataProvider:   dataAgreement.DataProvider,
	}

	var falsify int
	for _, hash := range hashes {
		for _, aggHash := range offerDataHashes[0].DataHashes {
			if hash == aggHash.Hash {
				falsify++
			}
		}
	}

	falsify = len(hashes) - falsify

	cost.FalsifyCount = falsify
	if falsify <= 0 {
		logger.Info("Wrong falsify claim")
		response["message"] = "Wrong falsify claim"

		// No hourly penalty, just calculate reimbursement and refund based on deposits
		cost.ProviderReimbursement = dataAgreement.ProviderDeposit + dataAgreement.ConsumerDeposit + dataAgreement.Price
		cost.ConsumerRefund = 0

	} else if falsify > 0 {
		logger.Info("Valid falsify claim")
		response["message"] = "Valid falsify claim found"
		cost.ProviderReimbursement = 0
		cost.ConsumerRefund = dataAgreement.ProviderDeposit + dataAgreement.ConsumerDeposit + dataAgreement.Price
	}

	costAsBytes, _ := json.Marshal(cost)
	ctx.GetStub().PutState(cost.ID, costAsBytes)

	dataAgreement.State = false
	agreementAsBytes, _ = json.Marshal(dataAgreement)
	ctx.GetStub().PutState(dataAgreement.ID, agreementAsBytes)

	return response, nil
}
func (pm *dataManagement) FalsifyClaimForHistorical(ctx contractapi.TransactionContextInterface, offerId, _hashes, agreementId, costId string) (map[string]interface{}, error) {
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	hashes := strings.Split(_hashes, ",")
	costQuery := fmt.Sprintf(`{
        "selector": {
            "docType": "%s",
            "agreement":"%s"
        }
    }`, COST, agreementId)
	// in case the request is sent to manipulate existing calculated costs
	costExist := pm.checkCostExist(ctx, costQuery)
	if costExist {
		response["message"] = fmt.Sprintf("Cost already exists for selected agreement id %s", agreementId)
		return response, nil
	}

	queryString := fmt.Sprintf(`{
        "selector": {
           "docType": "%s",
           "offer_id": "%s"
        }
     }`, DATA_HASH, offerId)

	logger.Info(queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		response["message"] = "resultsIterator failed"
		return response, err
	}
	defer resultsIterator.Close()

	var offerDataHashes []*OfferDataHash
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			response["message"] = "resultsIterator failed"
			return response, err
		}

		var offerDataHash OfferDataHash
		err = json.Unmarshal(queryResponse.Value, &offerDataHash)
		if err != nil {
			response["message"] = "Failed to Unmarshal"
			return response, err
		}
		offerDataHashes = append(offerDataHashes, &offerDataHash)
	}
	agreementAsBytes, _ := ctx.GetStub().GetState(agreementId)
	var dataAgreement DataAgreement
	err = json.Unmarshal(agreementAsBytes, &dataAgreement)
	if err != nil {
		return response, err
	}

	var escrow Escrow
	escrowAsBytes, err := ctx.GetStub().GetState(dataAgreement.EscrowID)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(escrowAsBytes, &escrow)
	if err != nil {
		return response, err
	}
	escrow.Released = true
	escrow.Status = REVOKED
	escrowAsBytes, err = json.Marshal(escrow)
	if err != nil {
		return response, err
	}
	err = ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
	if err != nil {
		return response, err
	}

	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")

	cost := Costs{
		CreatedAt:      today,
		DocType:        COST,
		ID:             costId,
		Agreement:      agreementId,
		OfferRequestID: dataAgreement.OfferRequestID,
		EscrowID:       dataAgreement.EscrowID,
		DataConsumer:   dataAgreement.DataConsumer,
		DataProvider:   dataAgreement.DataProvider,
	}

	var falsify int
	for _, hash := range hashes {
		for _, aggHash := range offerDataHashes[0].DataHashes {
			if hash == aggHash.Hash {
				falsify++
			}
		}
	}

	falsify = len(hashes) - falsify

	cost.FalsifyCount = falsify
	if falsify <= 0 {
		logger.Info("Wrong falsify claim")
		response["message"] = "Wrong falsify claim"

		// No hourly penalty, just calculate reimbursement and refund based on deposits
		cost.ProviderReimbursement = dataAgreement.ProviderDeposit + dataAgreement.ConsumerDeposit
		cost.ConsumerRefund = dataAgreement.Price

	} else if falsify > 0 {
		logger.Info("Valid falsify claim")
		response["message"] = "Valid falsify claim found"
		cost.ProviderReimbursement = 0
		cost.ConsumerRefund = dataAgreement.ProviderDeposit + dataAgreement.ConsumerDeposit + dataAgreement.Price
	}

	costAsBytes, _ := json.Marshal(cost)
	ctx.GetStub().PutState(cost.ID, costAsBytes)

	dataAgreement.State = false
	agreementAsBytes, _ = json.Marshal(dataAgreement)
	ctx.GetStub().PutState(dataAgreement.ID, agreementAsBytes)

	return response, nil
}

func (pm *dataManagement) ReleaseEscrow2(ctx contractapi.TransactionContextInterface, escrowId string, releaseCostID string, hours int, minutes int) error {
	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")
	costQuery := fmt.Sprintf(`{
		"selector": {
			"docType": "%s",
			"escrow_id":"%s"
		}
	}`, COST, escrowId)

	costExist := pm.checkCostExist(ctx, costQuery)
	if costExist {
		return fmt.Errorf("Cost already exists for selected escrow id %s", escrowId)
	}

	var escrow Escrow
	escrowAsBytes, err := ctx.GetStub().GetState(escrowId)

	if err != nil {
		return err
	}
	err = json.Unmarshal(escrowAsBytes, &escrow)
	if err != nil {
		return err
	}

	cost := Costs{
		CreatedAt:             today,
		DocType:               COST,
		ID:                    releaseCostID,
		Agreement:             escrow.AgreementID,
		ProviderReimbursement: escrow.ProviderDeposit,
		ConsumerRefund:        escrow.ConsumerDeposit,
		EscrowID:              escrowId,
		DataConsumer:          escrow.Consumer,
		DataProvider:          escrow.Provider,
		OfferRequestID:        escrow.OfferRequestID,
	}

	// // Parse the original escrow.EndDate
	// endDate, _ := time.Parse("2006-01-02 15:04", escrow.EndDate)

	// // Add hours and minutes to the escrow.EndDate
	// hoursToAdd := hours     // Adjust this value as needed
	// minutesToAdd := minutes // Adjust this value as needed
	// endDate = endDate.Add(time.Duration(hoursToAdd) * time.Hour)
	// endDate = endDate.Add(time.Duration(minutesToAdd) * time.Minute)

	// // Update the escrow.EndDate with the updated date and time
	// escrow.EndDate = endDate.Format("2006-01-02 15:04")

	newEndDate, _ := time.Parse("2006-01-02 15:04", escrow.EndDate)
	if hours > 0 {
		newEndDate = newEndDate.Add(time.Duration(hours) * time.Hour)
	}
	if minutes > 0 {
		newEndDate = newEndDate.Add(time.Duration(minutes) * time.Minute)
	}

	logger.Info("newEndDate %s", newEndDate.Format("2006-01-02 15:04"))
	logger.Info("today %s", today)

	if escrow.Status == ACTIVE && newEndDate.Format("2006-01-02 15:04") == today {
		escrow.Released = true
		escrow.Status = EXPIRED
		escrowAsBytes, err = json.Marshal(escrow)
		if err != nil {
			return err
		}
		cost.ProviderReimbursement = escrow.ProviderDeposit + escrow.ConsumerPayment
		cost.ConsumerRefund = escrow.ConsumerDeposit
		ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
		var agreement DataAgreement
		agreementAsBytes, _ := ctx.GetStub().GetState(escrow.AgreementID)
		json.Unmarshal(agreementAsBytes, &agreement)
		agreement.State = false
		agreementAsBytes, err = json.Marshal(agreement)
		ctx.GetStub().PutState(agreement.ID, agreementAsBytes)
		if err != nil {
			return err
		}
	}

	costAsBytes, err := json.Marshal(cost)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(cost.ID, costAsBytes)
	if err != nil {
		return err
	}
	return nil
}

// CalculateCost this func is not invoked in this version .. skip it
func (pm *dataManagement) CalculateCost(ctx contractapi.TransactionContextInterface, agreementID, actionCode string) bool {

	docID := generateUUID()

	agreementAsBytes, _ := ctx.GetStub().GetState(agreementID)

	var dataAgreement DataAgreement
	err := json.Unmarshal(agreementAsBytes, &dataAgreement)
	if err != nil {
		return false
	}

	if dataAgreement.State == false {
		return false
	}

	escrowAsBytes, _ := ctx.GetStub().GetState(dataAgreement.EscrowID)

	var escrow Escrow
	err = json.Unmarshal(escrowAsBytes, &escrow)
	if err != nil {
		return false
	}

	rate := escrow.ProviderDeposit + escrow.ConsumerDeposit
	dPrice := rate / 30

	var actPrice float64
	var reimb float64
	var refund float64
	ep := escrow.ProviderDeposit
	ec := escrow.ConsumerDeposit

	switch actionCode {
	case CODE_A:
		actPrice = dPrice * calculateDays(dataAgreement.StartDate, dataAgreement.EndDate)
		reimb = actPrice
		refund = (escrow.ConsumerPayment - actPrice) + ep + ec
	case CODE_B:
		actPrice = dPrice * calculateDays(dataAgreement.StartDate, dataAgreement.EndDate)
		reimb = actPrice + ep + ec
		refund = escrow.ConsumerPayment - actPrice
	case CODE_C:
		actPrice = dPrice * calculateDays(dataAgreement.StartDate, dataAgreement.EndDate)
		reimb = actPrice + ep
		refund = (escrow.ConsumerPayment - actPrice) + ec
	default:
		fmt.Println("Something went wrong")
		return false
	}

	reimb = float64(int(reimb*100)) / 100
	refund = float64(int(refund*100)) / 100

	var details Costs

	details.ID = docID
	details.DocType = COST
	details.Agreement = agreementID
	details.ConsumerRefund = refund
	details.ProviderReimbursement = reimb
	details.DataProvider = escrow.Provider
	details.DataConsumer = escrow.Consumer
	details.EscrowID = escrow.ID
	details.OfferRequestID = escrow.OfferRequestID

	detailsAsBytes, _ := json.Marshal(details)
	err = ctx.GetStub().PutState(docID, detailsAsBytes)
	if err != nil {
		return false
	}

	dataAgreement.State = false
	ddataAgreementAsBytes, err := json.Marshal(dataAgreement)
	err = ctx.GetStub().PutState(agreementID, ddataAgreementAsBytes)
	if err != nil {
		return false
	}

	escrow.Released = true
	escrowAsBytes, _ = json.Marshal(escrow)
	err = ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
	if err != nil {
		return false
	}

	return true
}

// not used yet
func (pm *dataManagement) DeleteData(ctx contractapi.TransactionContextInterface, keys string) {

	s := strings.Split(keys, ",")
	logger.Info(s)
	for _, key := range s {
		err := ctx.GetStub().DelState(key)
		if err != nil {
			logger.Error(err)
			return
		}
	}
	return
}

// LatencyClaim to check the frequency in appending the hash values and the time window will be around 9 mins
func (pm *dataManagement) LatencyClaim(ctx contractapi.TransactionContextInterface, offerId, agreementId, costId string) (map[string]interface{}, error) {

	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	costQuery := fmt.Sprintf(`{
		"selector": {
			"docType": "%s",
			"agreement":"%s"
		}
	}`, COST, agreementId)

	costExist := pm.checkCostExist(ctx, costQuery)
	if costExist {
		response["message"] = fmt.Sprintf("Cost already exist for selected agreement id %s", agreementId)
		return response, nil
	}

	queryString := fmt.Sprintf(`{
		"selector": {
		   "docType": "%s",
		   "offer_id": "%s"
		}
	 }`, DATA_HASH, offerId)

	logger.Info(queryString)
	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		response["message"] = "resultsIterator failed"
		return response, err
	}
	defer resultsIterator.Close()

	var offerDataHashes []*OfferDataHash
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			response["message"] = "resultsIterator failed"
			return response, err
		}

		var offerDataHash OfferDataHash
		err = json.Unmarshal(queryResponse.Value, &offerDataHash)
		if err != nil {
			response["message"] = "Failed to Unmarshal"
			return response, err
		}
		offerDataHashes = append(offerDataHashes, &offerDataHash)
	}
	agreementAsBytes, _ := ctx.GetStub().GetState(agreementId)

	var dataAgreement DataAgreement
	json.Unmarshal(agreementAsBytes, &dataAgreement)
	if len(offerDataHashes) == 1 {
		offerDatahash := offerDataHashes[0]

		offerDataHashesToCheck := filterByDataHashes(offerDatahash.DataHashes, dataAgreement.OfferDataHashID)

		if len(offerDataHashesToCheck) < 2 {
			fmt.Println("not found hash id not found")
			response["message"] = "not found hash id not found"
			return response, nil
		}
		var latencyCount int
		var ProviderReimbursement float64
		var ConsumerRefund float64
		var FalseProviderReimbursement float64
		var FalseConsumerRefund float64

		for idx := range offerDataHashesToCheck {

			if (idx - 1) >= 0 {
				// TODO check the time diff
				fmt.Println("TODO check the time diff")

				hours := calculateHours(offerDataHashesToCheck[idx-1].EntryDate, offerDataHashesToCheck[idx].EntryDate)

				offerRequestAsBytes, _ := ctx.GetStub().GetState(dataAgreement.OfferRequestID)

				var offerRequest OfferRequest
				json.Unmarshal(offerRequestAsBytes, &offerRequest)

				totalHours := calculateHours(dataAgreement.StartDate, dataAgreement.EndDate)
				pricePerHour := math.Round((offerRequest.Price/totalHours)*100) / 100

				logger.Info("Today", today)
				consumedHours := calculateHours(dataAgreement.StartDate, today)
				logger.Info("consumedHours", consumedHours)
				remainingHours := totalHours - consumedHours
				if hours > LATENCY_TIME {
					latencyCount++
					logger.Info("Latency exists")
					response["message"] = "Latency claim found to be valid"
					ProviderReimbursement += roundOff(pricePerHour * consumedHours)
					ConsumerRefund += roundOff((dataAgreement.ProviderDeposit + dataAgreement.ConsumerDeposit) + (remainingHours * pricePerHour))
				} else {
					response["message"] = "False claim raised"
					FalseProviderReimbursement += roundOff((dataAgreement.ProviderDeposit + dataAgreement.ConsumerDeposit) + (consumedHours * pricePerHour))
					FalseConsumerRefund += roundOff(remainingHours * pricePerHour)
				}

			} else {
				fmt.Println("Previous record not found")
				response["message"] = "NO OLDER RECORD FOUND TO CHECK LATENCY CLAIM"

			}
		}
		var escrow Escrow
		escrowAsBytes, _ := ctx.GetStub().GetState(dataAgreement.EscrowID)
		json.Unmarshal(escrowAsBytes, &escrow)
		if latencyCount == 0 {
			ConsumerRefund = FalseConsumerRefund
			ProviderReimbursement = FalseProviderReimbursement

		}

		escrow.Released = true
		escrow.Status = REVOKED
		dataAgreement.State = false
		ctx.GetStub().PutState(escrow.ID, escrowAsBytes)

		agreementAsBytes, _ = json.Marshal(dataAgreement)
		ctx.GetStub().PutState(dataAgreement.ID, agreementAsBytes)

		escrowAsBytes, _ = json.Marshal(escrow)
		ctx.GetStub().PutState(escrow.ID, escrowAsBytes)

		cost := Costs{
			LatencyCount:          latencyCount,
			CreatedAt:             today,
			DocType:               COST,
			ID:                    costId,
			Agreement:             agreementId,
			OfferRequestID:        dataAgreement.OfferRequestID,
			EscrowID:              dataAgreement.EscrowID,
			DataConsumer:          dataAgreement.DataConsumer,
			DataProvider:          dataAgreement.DataProvider,
			ConsumerRefund:        ConsumerRefund,
			ProviderReimbursement: ProviderReimbursement,
		}
		costAsBytes, _ := json.Marshal(cost)
		err = ctx.GetStub().PutState(cost.ID, costAsBytes)
		if err != nil {
			response["message"] = "Putstate failed"
			return response, err
		}

	} else {
		response["message"] = fmt.Sprintf("Multiple offer data hash found for offer id %s", offerId)
		return response, nil
	}
	return response, nil

}

// FalsifyClaim to check if the hash values are different from the one stored on BC
func (pm *dataManagement) FalsifyClaim(ctx contractapi.TransactionContextInterface, offerId, _hashes, agreementId, costId string) (map[string]interface{}, error) {
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	hashes := strings.Split(_hashes, ",")
	costQuery := fmt.Sprintf(`{
		"selector": {
			"docType": "%s",
			"agreement":"%s"
		}
	}`, COST, agreementId)
	// in case the request is sent to manipulate existing calculted costs
	costExist := pm.checkCostExist(ctx, costQuery)
	if costExist {
		response["message"] = fmt.Sprintf("Cost already exist for selected agreement id %s", agreementId)
		return response, nil
	}

	queryString := fmt.Sprintf(`{
		"selector": {
		   "docType": "%s",
		   "offer_id": "%s"
		}
	 }`, DATA_HASH, offerId)

	logger.Info(queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		response["message"] = "resultsIterator failed"
		return response, err
	}
	defer resultsIterator.Close()

	var offerDataHashes []*OfferDataHash
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			response["message"] = "resultsIterator failed"
			return response, err
		}

		var offerDataHash OfferDataHash
		err = json.Unmarshal(queryResponse.Value, &offerDataHash)
		if err != nil {
			response["message"] = "Failed to Unmarshal"
			return response, err
		}
		offerDataHashes = append(offerDataHashes, &offerDataHash)
	}
	agreementAsBytes, _ := ctx.GetStub().GetState(agreementId)
	var dataAgreement DataAgreement
	json.Unmarshal(agreementAsBytes, &dataAgreement)

	offerRequestAsBytes, _ := ctx.GetStub().GetState(dataAgreement.OfferRequestID)
	var offerRequest OfferRequest
	json.Unmarshal(offerRequestAsBytes, &offerRequest)
	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")

	cost := Costs{
		CreatedAt:      today,
		DocType:        COST,
		ID:             costId,
		Agreement:      agreementId,
		OfferRequestID: dataAgreement.OfferRequestID,
		EscrowID:       dataAgreement.EscrowID,
		DataConsumer:   dataAgreement.DataConsumer,
		DataProvider:   dataAgreement.DataProvider,
	}
	totalHours := calculateHours(dataAgreement.StartDate, dataAgreement.EndDate)
	pricePerHour := math.Round((offerRequest.Price/totalHours)*100) / 100

	logger.Info("Today", today)
	consumedHours := calculateHours(dataAgreement.StartDate, today)
	logger.Info("consumedHours", consumedHours)
	remainingHours := totalHours - consumedHours

	//
	var falsify int = 0

	for _, hash := range hashes {
		for _, aggHash := range offerDataHashes[0].DataHashes {
			if hash == aggHash.Hash {
				falsify++
			}
		}
	}

	falsify = len(hashes) - falsify

	//
	cost.FalsifyCount = falsify
	if falsify <= 0 {
		logger.Info("Wrong falsify claim")
		response["message"] = "Wrong falsify claim"

		cost.ProviderReimbursement = roundOff((dataAgreement.ProviderDeposit + dataAgreement.ConsumerDeposit) + (consumedHours * pricePerHour))
		cost.ConsumerRefund = roundOff(remainingHours * pricePerHour)

	} else if falsify > 0 {
		logger.Info("Valid falsify claim")
		response["message"] = "Valid falsify claim found"
		cost.ProviderReimbursement = roundOff(consumedHours * pricePerHour)
		cost.ConsumerRefund = roundOff((dataAgreement.ProviderDeposit + dataAgreement.ConsumerDeposit) + (remainingHours * pricePerHour))
	}

	costAsBytes, _ := json.Marshal(cost)
	ctx.GetStub().PutState(cost.ID, costAsBytes)

	var escrow Escrow
	escrowAsBytes, _ := ctx.GetStub().GetState(dataAgreement.EscrowID)
	json.Unmarshal(escrowAsBytes, &escrow)
	escrow.Released = true
	escrow.Status = REVOKED
	escrowAsBytes, _ = json.Marshal(escrow)
	ctx.GetStub().PutState(escrow.ID, escrowAsBytes)

	dataAgreement.State = false

	agreementAsBytes, _ = json.Marshal(dataAgreement)
	ctx.GetStub().PutState(dataAgreement.ID, agreementAsBytes)

	return response, nil
}

// insertDataOffer inserts HistoricalDataOffer into ledger
func (pm *dataManagement) InsertHistoricalDataOffer(ctx contractapi.TransactionContextInterface, dataOffer string) (map[string]interface{}, error) {

	var details HistoricalDataOffer
	logger.Info(dataOffer)
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	err := json.Unmarshal([]byte(dataOffer), &details)
	if err != nil {
		response["message"] = "Error occured while Unmarshal"
		return response, fmt.Errorf("Failed while unmarshling Order %s", err.Error())
	}

	exists, err := pm.ObjectExists(ctx, details.ID)
	if err != nil {
		response["message"] = "failed to get Offer"
		return response, fmt.Errorf("failed to get Offer: %v", err)
	}
	if exists {
		response["message"] = "DataOffer already exists"
		return response, fmt.Errorf("DataOffer already exists")
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return response, fmt.Errorf("failed to get verified OrgID: %v", err)
	}
	details.DocType = HISTORICALOFFER
	details.OwnerOrg = clientOrgID
	detailsAsBytes, _ := json.Marshal(details)
	err = ctx.GetStub().PutState(details.ID, detailsAsBytes)
	if err != nil {
		response["message"] = "Error occured while storing data"
		return response, fmt.Errorf("Failed to add DataOffer catch: %s", details.ID)
	}

	return response, nil
}

func (pm *dataManagement) UpdateHistoicalDataOffer(ctx contractapi.TransactionContextInterface, historicalOffer string) (map[string]interface{}, error) {

	var details HistoricalDataOffer
	logger.Info(historicalOffer)
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	err := json.Unmarshal([]byte(historicalOffer), &details)
	if err != nil {
		response["message"] = "Error occured while Unmarshal"
		return response, fmt.Errorf("Failed while unmarshling Order. %s", err.Error())
	}
	logger.Info(details)
	exists, err := pm.ObjectExists(ctx, details.ID)
	if err != nil {
		response["message"] = "failed to get journey"
		return response, fmt.Errorf("failed to get journey: %v", err)
	}
	if !exists {
		response["message"] = "No Record found"
		return response, fmt.Errorf("Journey does not exists")
	}

	var _offer HistoricalDataOffer
	_journeyAsBytes, _ := ctx.GetStub().GetState(details.ID)
	err = json.Unmarshal(_journeyAsBytes, &_offer)
	if err != nil {
		return response, err
	}
	details.DocType = HISTORICALOFFER
	detailsAsBytes, _ := json.Marshal(details)
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return response, fmt.Errorf("failed to get verified OrgID: %v", err)
	}
	logger.Info("clientOrgID : %s", clientOrgID)
	//logger.Info("details.OwnerOrg : %s", details.OwnerOrg)

	err = ctx.GetStub().PutState(details.ID, detailsAsBytes)
	if err != nil {
		response["message"] = "Error occured while updating data"
		return response, fmt.Errorf("Failed to update journey catch: %s", details.ID)
	}

	return response, nil
}
func (pm *dataManagement) GetAllHistoricalOffer(ctx contractapi.TransactionContextInterface) ([]*HistoricalDataOffer, error) {

	queryString := fmt.Sprintf(`{"selector":{"docType":"%s"}}`, HISTORICALOFFER)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var journeySchedule []*HistoricalDataOffer
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var journey HistoricalDataOffer
		err = json.Unmarshal(queryResponse.Value, &journey)
		if err != nil {
			return nil, err
		}
		logger.Info(journey.ID)

		journeySchedule = append(journeySchedule, &journey)
	}

	return journeySchedule, nil

}

func (pm *dataManagement) GetAllHistoricalDataOffer(ctx contractapi.TransactionContextInterface, creator string) ([]HistoricalQuery, error) {

	var queryString string
	logger.Info(creator)
	if len(creator) == 0 {
		queryString = fmt.Sprintf(`{"selector":{"docType":"%s"}}`, HISTORICALOFFER)
	} else {
		queryString = fmt.Sprintf(`{"selector":{"docType":"%s","creator":"%s"}}`, HISTORICALOFFER, creator)
	}

	logger.Info(queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var results []HistoricalQuery

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var historicalOffer *HistoricalDataOffer
		err = json.Unmarshal(queryResponse.Value, &historicalOffer)
		if err != nil {
			return nil, err
		}

		queryResultt := HistoricalQuery{Key: queryResponse.Key, Records: historicalOffer}
		results = append(results, queryResultt)
	}

	return results, nil
}

func (pm *dataManagement) FalsifyClaimUseCase2(ctx contractapi.TransactionContextInterface, offerId, _hashes, agreementId, costId string) (map[string]interface{}, error) {
	response := make(map[string]interface{})
	response["txId"] = ctx.GetStub().GetTxID()
	hashes := strings.Split(_hashes, ",")
	costQuery := fmt.Sprintf(`{
        "selector": {
            "docType": "%s",
            "agreement":"%s"
        }
    }`, COST, agreementId)
	// Check if the request is sent to manipulate existing calculated costs
	costExist := pm.checkCostExist(ctx, costQuery)
	if costExist {
		response["message"] = fmt.Sprintf("Cost already exists for selected agreement id %s", agreementId)
		return response, nil
	}

	queryString := fmt.Sprintf(`{
        "selector": {
           "docType": "%s",
           "offer_id": "%s"
        }
     }`, DATA_HASH, offerId)

	logger.Info(queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		response["message"] = "resultsIterator failed"
		return response, err
	}
	defer resultsIterator.Close()

	var offerDataHashes []*OfferDataHash
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			response["message"] = "resultsIterator failed"
			return response, err
		}

		var offerDataHash OfferDataHash
		err = json.Unmarshal(queryResponse.Value, &offerDataHash)
		if err != nil {
			response["message"] = "Failed to Unmarshal"
			return response, err
		}
		offerDataHashes = append(offerDataHashes, &offerDataHash)
	}
	agreementAsBytes, _ := ctx.GetStub().GetState(agreementId)
	var dataAgreement DataAgreement
	err = json.Unmarshal(agreementAsBytes, &dataAgreement)
	if err != nil {
		return response, err
	}

	// Update escrow status and release before applying penalty
	var escrow Escrow
	escrowAsBytes, err := ctx.GetStub().GetState(dataAgreement.EscrowID)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(escrowAsBytes, &escrow)
	if err != nil {
		return response, err
	}
	escrow.Released = true
	escrow.Status = REVOKED
	escrowAsBytes, err = json.Marshal(escrow)
	if err != nil {
		return response, err
	}
	err = ctx.GetStub().PutState(escrow.ID, escrowAsBytes)
	if err != nil {
		return response, err
	}

	logger.Info("offerDataHashes", offerDataHashes)
	logger.Info("dataAgreement.OfferDataHashID", dataAgreement.OfferDataHashID)

	// Calculate the number of missing data hashes
	missingDataHashes := make([]string, 0)
	for _, hash := range dataAgreement.OfferDataHashID {
		found := false
		for _, offerDataHash := range offerDataHashes {
			if hash == offerDataHash.DataHashes[0].ID {
				found = true
				break
			}
		}
		if !found {
			missingDataHashes = append(missingDataHashes, hash)
		}
	}
	timestamp, _ := ctx.GetStub().GetTxTimestamp()
	currentTime := time.Unix(timestamp.GetSeconds(), 0)
	today := currentTime.Format("2006-01-02 15:04")

	cost := Costs{
		CreatedAt:      today,
		DocType:        COST,
		ID:             costId,
		Agreement:      agreementId,
		OfferRequestID: dataAgreement.OfferRequestID,
		EscrowID:       dataAgreement.EscrowID,
		DataConsumer:   dataAgreement.DataConsumer,
		DataProvider:   dataAgreement.DataProvider,
		FalsifyCount:   len(hashes) - len(missingDataHashes),
	}

	// Check if there are any missing data hashes and apply the penalty
	if len(missingDataHashes) > 0 {
		logger.Info("Missing data hash(es) - applying penalty")
		response["message"] = "Missing data hash(es) - applying penalty"

		// Calculate the refund amount for the consumer
		var refundAmount float64

		missingFilesRefund := float64(len(missingDataHashes)) / float64(len(hashes)) // You can adjust this amount as needed

		penalty := math.Round(dataAgreement.Price * missingFilesRefund)
		refundAmount = math.Round(penalty + dataAgreement.ConsumerDeposit + dataAgreement.ProviderDeposit)
		// refundAmount = penalty + dataAgreement.ConsumerDeposit + dataAgreement.ProviderDeposit

		cost.ProviderReimbursement = dataAgreement.Price - penalty
		cost.ConsumerRefund = refundAmount

	} else {
		logger.Info("No missing data hashes")
		response["message"] = "No missing data hashes"

		// No penalty, just calculate reimbursement and refund based on deposits
		cost.ProviderReimbursement = dataAgreement.ProviderDeposit + dataAgreement.ConsumerDeposit + dataAgreement.Price
		cost.ConsumerRefund = 0
	}

	costAsBytes, err := json.Marshal(cost)
	if err != nil {
		return response, err
	}
	err = ctx.GetStub().PutState(cost.ID, costAsBytes)
	if err != nil {
		return response, err
	}

	dataAgreement.State = false
	agreementAsBytes, err = json.Marshal(dataAgreement)
	if err != nil {
		return response, err
	}
	err = ctx.GetStub().PutState(dataAgreement.ID, agreementAsBytes)
	if err != nil {
		return response, err
	}

	return response, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&dataManagement{})
	if err != nil {
		logger.Info("Error creating chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		logger.Info("Error starting chaincode: %v", err)
	}
}
