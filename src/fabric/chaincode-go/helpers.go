package main

import (
	"encoding/json"
	"fmt"
	"time"

	"math"

	guuid "github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)
// getClientOrgID return clientOrgId for getting transaction initiator's Identity
func getClientOrgID(ctx contractapi.TransactionContextInterface, verifyOrg bool) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed getting client's orgID: %v", err)
	}

	if verifyOrg {
		err = verifyClientOrgMatchesPeerOrg(clientOrgID)
		if err != nil {
			return "", err
		}
	}

	return clientOrgID, nil
}

// verifyClientOrgMatchesPeerOrg checks the client org id matches the peer org id.
func verifyClientOrgMatchesPeerOrg(clientOrgID string) error {
	peerOrgID, err := shim.GetMSPID()
	if err != nil {
		return fmt.Errorf("failed getting peer's orgID: %v", err)
	}

	if clientOrgID != peerOrgID {
		return fmt.Errorf("client from org %s is not authorized to read or write private data from an org %s peer",
			clientOrgID,
			peerOrgID,
		)
	}

	return nil
}
// checkCostExist to avoid recalculating
func (pm *dataManagement) checkCostExist(ctx contractapi.TransactionContextInterface, queryString string) bool {
	fmt.Println(queryString)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return false
	}
	defer resultsIterator.Close()

	var costs []*Costs
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return false
		}

		var cost Costs
		err = json.Unmarshal(queryResponse.Value, &cost)
		if err != nil {
			return false
		}
		costs = append(costs, &cost)
	}

	return len(costs) > 0
}

func (pm *dataManagement) ObjectExists(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
	assetBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, fmt.Errorf("failed to read Object %s from world state. %v", key, err)
	}

	return assetBytes != nil, nil
}

func generateUUID() string {
	id := guuid.New()
	return id.String()
}

func calculateDays(from, to string) float64 {
	format := "2006-01-02 15:04 MST"
	from += " 00:00 UTC"
	to += " 00:00 UTC"
	start, err := time.Parse(format, from)
	if err != nil {
		fmt.Println(err.Error())
	}
	end, err := time.Parse(format, to)
	if err != nil {
		fmt.Println(err.Error())
	}
	diff := end.Sub(start)
	return (diff.Hours() / 24)
}

func rotateLeft(s string, n int) string {
	n = n % len(s)
	return s[n:] + s[0:n]
}

func calculateHours(from, to string) float64 {
	logger.Info("calculateHours:", from, to)
	format := "2006-01-02 15:04 MST"
	from += " UTC"
	to += " UTC"
	start, err := time.Parse(format, from)
	if err != nil {
		fmt.Println(err.Error())
	}
	end, err := time.Parse(format, to)
	if err != nil {
		fmt.Println(err.Error())
	}
	diff := end.Sub(start)
	return math.Round(diff.Hours()*100) / 100
}

func Search(length int, f func(index int) bool) int {
	for index := 0; index < length; index++ {
		if f(index) {
			return index
		}
	}
	return -1
}

func roundOff(num float64) float64 {
	return math.Round((num)*100) / 100

}

func filterByDataHashes(fu []DataHash, su []string) (out []DataHash) {
	f := make(map[string]struct{}, len(su))
	for _, u := range su {
		f[u] = struct{}{}
	}
	for _, u := range fu {
		if _, ok := f[u.ID]; ok {
			out = append(out, u)
		}
	}
	return
}
