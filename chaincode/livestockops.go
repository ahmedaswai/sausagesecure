package chaincode

import (
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/samber/lo"
	model "sausagesecure/chaincode/model"
)

// InitLedger adds a base set of assets to the ledger

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	farms, err := s.fetchAssetsByType(ctx, "Farm")

	if err != nil {
		return err
	}
	if len(farms) > 0 {
		return nil
	}

	newFarms := []model.Farm{

		{City: "Beni Suef", FullAddress: "مركز اهناسيا امام سيد مرعي بتاع الكاوتش", Latitude: 29.0661, Longitude: 31.0994},
		{City: "Giza", FullAddress: "مركز العياط شارع  عماد محصلة", Latitude: 29.3084, Longitude: 30.8428},
		{City: "Fayoum", FullAddress: "مركز طامية عزبة صوفي ابو طالب ", Latitude: 30.013056, Longitude: 31.208853},
	}
	for idx, asset := range newFarms {
		_, err := s.addAsset(ctx, &asset, idx)

		if err != nil {
			return err
		}

	}

	return nil
}

func (s *SmartContract) AddLiveStock(ctx contractapi.TransactionContextInterface, liveStockJson string) (*model.LiveStock, error) {

	var liveStock *model.LiveStock
	err := json.Unmarshal([]byte(liveStockJson), &liveStock)
	if err != nil {
		return nil, err
	}
	ik, err := s.addAsset(ctx, liveStock, 0)
	stock, _ := ik.(*model.LiveStock)
	return stock, err

}

func (s *SmartContract) UpdateLiveStockData(ctx contractapi.TransactionContextInterface, liveStockId string, farmId string) (*model.LiveStock, error) {
	liveStockAsset, err := s.fetchAssetById(ctx, liveStockId, "LiveStock")
	if err != nil {
		return nil, err
	}

	liveStock := liveStockAsset.(*model.LiveStock)

	liveStock.FarmID = farmId
	liveStockAssetJson, err := json.Marshal(liveStock)

	err = ctx.GetStub().PutState(liveStockId, liveStockAssetJson)

	if err != nil {
		return nil, err
	}
	return liveStock, nil

}

func (s *SmartContract) AddLiveStockHistory(ctx contractapi.TransactionContextInterface, history string) (*model.LiveStockHistory, error) {

	var liveStock *model.LiveStockHistory
	err := json.Unmarshal([]byte(history), &liveStock)
	if err != nil {
		return nil, err
	}
	ik, err := s.addAsset(ctx, liveStock, 0)
	stock, _ := ik.(*model.LiveStockHistory)
	return stock, err

}

func (s *SmartContract) FetchAllFarms(ctx contractapi.TransactionContextInterface) ([]*model.Farm, error) {
	assets, err := s.fetchAssetsByType(ctx, "Farm")
	if err != nil {
		return nil, err
	}
	farms := lo.Map(assets, func(item model.IAsset, index int) *model.Farm {
		return item.(*model.Farm)
	})
	return farms, nil
}
