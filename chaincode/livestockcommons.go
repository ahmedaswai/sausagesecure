package chaincode

import (
	"encoding/json"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"reflect"
	"sausagesecure/chaincode/model"
	"strconv"
	"strings"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) addAsset(ctx contractapi.TransactionContextInterface, asset model.IAsset, transNum int) (model.IAsset, error) {

	id := strings.Join([]string{ctx.GetStub().GetTxID(), strconv.Itoa(transNum + 1)}, "")
	asset.SetID(id)

	assetObjectType := reflect.TypeOf(asset)
	asset.SetObjectName(assetObjectType.Elem().Name())

	jsonAsset, err := json.Marshal(asset)
	if err != nil {
		return nil, err
	}

	ledgerErr := ctx.GetStub().PutState(asset.GetID(), jsonAsset)

	if err != nil {
		return nil, ledgerErr
	}

	return asset, nil
}

func (s *SmartContract) createObjectTypeQuery(objectType string) string {

	query := model.CouchDBQuery{Selector: model.ObjectQuerySelector{ObjectName: objectType}, IndexDescriptor: []string{"_design/objectTypeDoc", "objectTypeIdx"}}
	queryJson, err := json.Marshal(query)
	if err != nil {
		return ""
	}
	return string(queryJson)
}

func (s *SmartContract) fetchMeatPackagePriceQuery(meatPackageID string) string {

	query := model.CouchMeatPackagePriceQuery{Selector: model.MeatPackagePriceQuery{MeatPackageID: meatPackageID}}
	queryJson, err := json.Marshal(query)
	if err != nil {
		return ""
	}
	return string(queryJson)
}

func (s *SmartContract) fetchAssetsByType(ctx contractapi.TransactionContextInterface, objectType string) ([]model.IAsset, error) {
	query := s.createObjectTypeQuery(objectType)
	return s.fetchAssetsByQuery(ctx, query, objectType)

}

func (s *SmartContract) fetchAssetsByQuery(ctx contractapi.TransactionContextInterface, query string, objectType string) ([]model.IAsset, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}

	defer func(resultsIterator shim.StateQueryIteratorInterface) {
		err := resultsIterator.Close()
		if err != nil {
			log.Panic("Error in Close StateQueryIteratorInterface", err)
		}
	}(resultsIterator)

	var assets []model.IAsset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset = s.fetchTypeName(objectType)
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

func (s *SmartContract) fetchTypeName(objectType string) model.IAsset {
	switch objectType {
	case "Farm":
		return new(model.Farm)
	case "LiveStock":
		return new(model.LiveStock)
	case "LiveStockHistory":
		return new(model.LiveStockHistory)

	case "LiveStockFeedRecord":
		return new(model.LiveStockFeedRecord)

	case "VaccinationRecord":
		return new(model.VaccinationRecord)

	case "ButcheryTransaction":
		return new(model.ButcheryTransaction)

	case "MeatPackage":
		return new(model.MeatPackage)

	default:
		return nil

	}
}

func (s *SmartContract) fetchAssetsByLiveStockId(ctx contractapi.TransactionContextInterface, liveStockId string, objectType string) ([]model.IAsset, error) {
	query := model.CouchLiveStockDBQuery{Selector: model.LiveStockIdSelector{LiveStockID: liveStockId, ObjectName: objectType},
		IndexDescriptor: []string{"_design/liveStockIDDoc", "liveStockIDIdx"}}
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	assets, err := s.fetchAssetsByQuery(ctx, string(queryJson), objectType)

	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (s *SmartContract) fetchAssetById(ctx contractapi.TransactionContextInterface, assetId string, objectType string) (model.IAsset, error) {

	assetJson, err := ctx.GetStub().GetState(assetId)
	if err != nil {
		return nil, err
	}
	var assetObj = s.fetchTypeName(objectType)
	err = json.Unmarshal(assetJson, assetObj)
	if err != nil {
		return nil, err
	}
	return assetObj, nil
}

func (s *SmartContract) fetchMother(ctx contractapi.TransactionContextInterface, motherID string) (*model.LiveStock, error) {
	if len(motherID) <= 0 {
		return nil, nil
	}

	motherLiveStock, err := s.fetchAssetById(ctx, motherID, "LiveStock")

	if err != nil {
		return nil, err
	}
	motherObject := motherLiveStock.(*model.LiveStock)

	return motherObject, nil
}

func (s *SmartContract) fetchMeatPackage(ctx contractapi.TransactionContextInterface, meatPackageID string) (*model.MeatPackage, *model.MeatPackagePrice, error) {

	meatPackageAsset, err := s.fetchAssetById(ctx, meatPackageID, "MeatPackage")

	if err != nil {
		return nil, nil, err
	}
	meatPackage, _ := meatPackageAsset.(*model.MeatPackage)

	mspId, err := ctx.GetClientIdentity().GetMSPID()

	if err != nil || mspId != "Org1MSP" {
		return meatPackage, nil, err
	}

	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(PRIVATE_DATA_COLLECTION, s.fetchMeatPackagePriceQuery(meatPackageID))

	if err != nil {
		return meatPackage, nil, err
	}

	defer func(resultsIterator shim.StateQueryIteratorInterface) {
		err := resultsIterator.Close()
		if err != nil {
			log.Panic("Error in Close StateQueryIteratorInterface", err)
		}
	}(resultsIterator)

	var meatPackagePrice *model.MeatPackagePrice

	if resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return meatPackage, nil, err
		}
		err = json.Unmarshal(queryResponse.Value, &meatPackagePrice)
		if err != nil {
			return meatPackage, nil, err
		}

	}
	return meatPackage, meatPackagePrice, nil
}
