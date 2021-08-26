package properties

const (
	// DbName :資料庫名
	DbName = "leapsy_env" //DB

	// CollectionNameOfCheckInRecord :Collection名稱:打卡紀錄
	CollectionNameOfCheckInRecord = "check_in_record" //Collection

	// CollectionNameOfCheckInStatistics :Collection名:打卡統計
	CollectionNameOfCheckInStatistics = "check_in_statistics" //Collection

	// const CollectionName = "persion"                                //Collection //範例程式

	// API Port
	PortOfAPI = 8082 //API port

	// PortOfMongoDB :MongoDB的Port
	PortOfMongoDB string = "27017"

	//docker HostOfMongoDB string = "172.17.0.2"
	HostOfMongoDB string = "localhost"
)
