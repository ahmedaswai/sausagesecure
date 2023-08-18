package model

type ObjectQuerySelector struct {
	ObjectName string
}

type MeatPackagePriceQuery struct {
	MeatPackageID string
}

type LiveStockIdSelector struct {
	ObjectName  string
	LiveStockID string
}
type CouchDBQuery struct {
	Selector        ObjectQuerySelector `json:"selector"`
	IndexDescriptor []string            `json:"use_index"`
}

type LiveMeatPackageTracking struct {
	MeatPackage         *MeatPackage
	LiveStock           *LiveStock
	LiveStockMother     *LiveStock
	MeatPackagePrice    *MeatPackagePrice
	Farm                *Farm
	LiveStockHistory    []*LiveStockHistory
	VaccinationRecord   []*VaccinationRecord
	ButcheryTransaction *ButcheryTransaction
	LiveStockFeedRecord []*LiveStockFeedRecord
}

type CouchLiveStockDBQuery struct {
	Selector        LiveStockIdSelector `json:"selector"`
	IndexDescriptor []string            `json:"use_index"`
}
type CouchMeatPackagePriceQuery struct {
	Selector MeatPackagePriceQuery `json:"selector"`
}
type IAsset interface {
	GetID() string
	SetID(id string)

	SetObjectName(objectName string)
}

type Farm struct {
	City        string
	FullAddress string
	ID          string
	Latitude    float64
	Longitude   float64
	ObjectName  string
}

type LiveStock struct {
	DateOfBirth uint64
	FarmID      string
	ID          string
	MotherID    string
	ObjectName  string
	Species     string
	StartWeight int
}

type LiveStockHistory struct {
	CurrentWeight   int
	HealthCondition string // Intact , Sick
	ID              string
	LiveStockID     string
	ObjectName      string
	TransactionDate uint64
}

type LiveStockFeedRecord struct {
	Amount             int
	FeedType           string
	FrequencyInHours   int
	ID                 string
	LiveStockID        string
	NutritionalContent string // This attribute provides information about the nutritional
	// content of the feed, such as the levels
	// of protein, carbohydrates, and fiber
	ObjectName      string
	TransactionDate uint64
}

type VaccinationRecord struct {
	DoseAmount      int
	VaccinationDate uint64
	VaccinationName string
	ID              string
	LiveStockID     string
	ObjectName      string
}

type ButcheryTransaction struct {
	ButcheryLocation    string
	ID                  string
	LiveStockID         string
	LiveStockPureWeight int
	ObjectName          string
	TransactionDate     uint64
}

type MeatPackage struct {
	CutImageLink  string
	CutName       string
	ID            string
	LiveStockID   string
	ObjectName    string
	PackageWeight float32
	PackagingDate uint64
}

// private data collection to manage the prices between multiple Organizations

type MeatPackagePrice struct {
	ID            string
	LiveStockID   string
	MeatPackageID string
	ObjectName    string
	Price         int
}

type LiveStockUpdateLog struct {
	LiveStock       *LiveStock
	TransactionId   string
	IsDeleted       bool
	TransactionTime int32
}
