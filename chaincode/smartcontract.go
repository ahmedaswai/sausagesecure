package chaincode

import (
	"encoding/json"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	model "sausagesecure/chaincode/model"
)

const PRIVATE_DATA_COLLECTION = "Org1_PDC2"

// InitLedger adds a base set of assets to the ledger

func (s *SmartContract) AddLiveStockFeedRecord(ctx contractapi.TransactionContextInterface, history string) (*model.LiveStockFeedRecord, error) {

	var liveStock *model.LiveStockFeedRecord
	err := json.Unmarshal([]byte(history), &liveStock)
	if err != nil {
		return nil, err
	}
	ik, err := s.addAsset(ctx, liveStock, 0)
	stock, _ := ik.(*model.LiveStockFeedRecord)
	return stock, err

}

func (s *SmartContract) AddVaccinationRecord(ctx contractapi.TransactionContextInterface, asset string) (*model.VaccinationRecord, error) {

	var liveStock *model.VaccinationRecord
	err := json.Unmarshal([]byte(asset), &liveStock)
	if err != nil {
		return nil, err
	}
	ik, err := s.addAsset(ctx, liveStock, 0)
	stock, _ := ik.(*model.VaccinationRecord)
	return stock, err

}

func (s *SmartContract) AddButcheryTransaction(ctx contractapi.TransactionContextInterface, asset string) (*model.ButcheryTransaction, error) {

	var liveStock *model.ButcheryTransaction
	err := json.Unmarshal([]byte(asset), &liveStock)
	if err != nil {
		return nil, err
	}
	ik, err := s.addAsset(ctx, liveStock, 0)
	stock, _ := ik.(*model.ButcheryTransaction)
	return stock, err

}

func (s *SmartContract) FetchLiveStockUpdateLog(ctx contractapi.TransactionContextInterface, liveStockId string) ([]*model.LiveStockUpdateLog, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(liveStockId)
	if err != nil {
		return nil, err
	}
	defer func(resultsIterator shim.HistoryQueryIteratorInterface) {
		err := resultsIterator.Close()
		if err != nil {
			log.Panic("Error in Close HistoryQueryIteratorInterface", err)
		}
	}(resultsIterator)

	var liveStockHistoryList []*model.LiveStockUpdateLog

	for resultsIterator.HasNext() {
		if err != nil {
			return nil, err
		}
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var asset *model.LiveStock

		err = json.Unmarshal(queryResponse.Value, &asset)

		if err != nil {
			return nil, err
		}

		liveStockRecord := model.LiveStockUpdateLog{LiveStock: asset,
			TransactionId: queryResponse.TxId, IsDeleted: queryResponse.IsDelete,
			TransactionTime: queryResponse.Timestamp.Nanos}
		liveStockHistoryList = append(liveStockHistoryList, &liveStockRecord)
	}
	return liveStockHistoryList, nil
}
