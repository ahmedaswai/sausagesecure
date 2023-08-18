package chaincode

import (
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/samber/lo"
	"log"
	"sausagesecure/chaincode/model"
	"strconv"
	"strings"
)

func (s *SmartContract) AddMeatPackage(ctx contractapi.TransactionContextInterface, asset string) (*model.MeatPackage, error) {

	var liveStock *model.MeatPackage
	err := json.Unmarshal([]byte(asset), &liveStock)
	if err != nil {
		return nil, err
	}

	ik, err := s.addAsset(ctx, liveStock, 0)
	stock, _ := ik.(*model.MeatPackage)
	return stock, err

}

func (s *SmartContract) FetchMeatPackageHistory(ctx contractapi.TransactionContextInterface, meatPackageID string) (*model.LiveMeatPackageTracking, error) {

	meatPackage, meatPackagePrice, err := s.fetchMeatPackage(ctx, meatPackageID)

	if err != nil {
		return nil, err
	}
	liveStock, err := s.fetchAssetById(ctx, meatPackage.LiveStockID, "LiveStock")

	if err != nil {
		return nil, err
	}
	liveStockObj, _ := liveStock.(*model.LiveStock)

	farmLiveStock, err := s.fetchAssetById(ctx, liveStockObj.FarmID, "Farm")

	if err != nil {
		return nil, err
	}

	motherLiveStock, err := s.fetchMother(ctx, liveStockObj.MotherID)

	liveStockHistory, err := s.fetchAssetsByLiveStockId(ctx, meatPackage.LiveStockID, "LiveStockHistory")

	if err != nil {
		return nil, err
	}

	liveStockHistoryRs := lo.Map(liveStockHistory, func(item model.IAsset, index int) *model.LiveStockHistory {
		return item.(*model.LiveStockHistory)
	})

	vaccinationRecord, err := s.fetchAssetsByLiveStockId(ctx, meatPackage.LiveStockID, "VaccinationRecord")
	if err != nil {
		return nil, err
	}
	vaccinationRecordRs := lo.Map(vaccinationRecord, func(item model.IAsset, index int) *model.VaccinationRecord {
		return item.(*model.VaccinationRecord)
	})

	butcheryTransaction, err := s.fetchAssetsByLiveStockId(ctx, meatPackage.LiveStockID, "ButcheryTransaction")

	if err != nil {
		return nil, err
	}
	butcheryTransactionRs := lo.Map(butcheryTransaction, func(item model.IAsset, index int) *model.ButcheryTransaction {
		return item.(*model.ButcheryTransaction)
	})

	liveStockFeedRecord, err := s.fetchAssetsByLiveStockId(ctx, meatPackage.LiveStockID, "LiveStockFeedRecord")
	if err != nil {
		return nil, err
	}
	liveStockFeedRecordRs := lo.Map(liveStockFeedRecord, func(item model.IAsset, index int) *model.LiveStockFeedRecord {
		return item.(*model.LiveStockFeedRecord)
	})

	log.Printf("LiveMeatPackageTracking has the next props  meatpackage Id %s and liveStockId %s "+
		" size of liveStockFeedRecord %d size of butcheryTransaction %d size of butcheryTransaction %d",
		meatPackageID, meatPackage.LiveStockID, len(liveStockFeedRecordRs), len(butcheryTransactionRs), len(butcheryTransactionRs))

	return &model.LiveMeatPackageTracking{MeatPackage: meatPackage, LiveStock: liveStock.(*model.LiveStock),
		LiveStockMother: motherLiveStock, LiveStockHistory: liveStockHistoryRs,
		Farm: farmLiveStock.(*model.Farm), VaccinationRecord: vaccinationRecordRs, ButcheryTransaction: butcheryTransactionRs[0],
		LiveStockFeedRecord: liveStockFeedRecordRs, MeatPackagePrice: meatPackagePrice}, nil
}

func (s *SmartContract) AddMeatPackagePrice(ctx contractapi.TransactionContextInterface, meatPackagePriceJson string) (*model.MeatPackagePrice, error) {
	var meatPackageObj *model.MeatPackagePrice
	err := json.Unmarshal([]byte(meatPackagePriceJson), &meatPackageObj)
	if err != nil {
		return nil, err
	}
	id := strings.Join([]string{ctx.GetStub().GetTxID(), strconv.Itoa(1)}, "")
	meatPackageObj.SetID(id)
	meatPackageObj.SetObjectName("MeatPackagePrice")

	meatPackagePriceEncoded, err := json.Marshal(meatPackageObj)

	if err != nil {
		return nil, err
	}

	var errPutPrivateData = ctx.GetStub().PutPrivateData(PRIVATE_DATA_COLLECTION, id, meatPackagePriceEncoded)

	if errPutPrivateData != nil {
		return nil, errPutPrivateData
	}
	return meatPackageObj, nil

}

func (s *SmartContract) FetchMeatPackagePrice(ctx contractapi.TransactionContextInterface, meatPackageId string) (*model.MeatPackagePrice, error) {

	_, meatPackagePriceJson, err := s.fetchMeatPackage(ctx, meatPackageId)

	if err != nil {
		return nil, err
	}

	return meatPackagePriceJson, nil

}
