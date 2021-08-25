package controller

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"

	"leapsy.com/db"
	// "my-rest-api/db"
	settings "leapsy.com/properties"
	// "my-rest-api/settings"
	//"leapsy.com/model"
	//"my-rest-api/model" //只有在create(insert)或update 才需要import model
)

// const dbName = "leapsy_env"                                     //DB
// const collectionNameOfCheckInRecord = "check_in_record"         //Collection
// const collectionNameOfCheckInStatistics = "check_in_statistics" //Collection
// const collectionName = "persion"                                //Collection
// const port = 8081                                               //API port
// const port = 8000 //API port

// 建立GET POST 路徑
func NewPersonController() {

	fmt.Println("測試")
	app := fiber.New()

	/*建立 checkInRecord 路徑*/
	app.Get("/checkInRecord/query/:date", getCheckInRecord) //應到人員資料
	app.Get("/checkInRecord/query/:year/:month/:day", getCheckInRecordYMD)
	// app.Get("/checkInRecord/attendance/:date?", getAttendanceOfCheckInStatistics) //實到人員資料
	// app.Get("/checkInRecord/notArrived/:date?", getNotArrivedOfCheckInStatistics) //未到人員資料
	//app.Post("/person", createPerson)
	//app.Put("/person/:id", updatePerson)
	//app.Delete("/person/:id", deletePerson)

	/*建立 checkInStatistics 路徑*/
	// app.Get("/checkInStatistics/query/:date?", getCheckInStatistics) //統計資料
	//app.Post("/person", createPerson)
	//app.Put("/person/:id", updatePerson)
	//app.Delete("/person/:id", deletePerson)

	/*建立範例 person 路徑*/
	// app.Get("/person/:id?", getPerson)
	// app.Post("/person", createPerson)
	// app.Put("/person/:id", updatePerson)
	// app.Delete("/person/:id", deletePerson)

	app.Listen(settings.PortOfAPI)
}

/* 以下為 CheckInRecord 相關 functions */
// 取得指定日期<應到>人員資料
func getCheckInRecord(c *fiber.Ctx) {

	// 取得 collection
	collection, err := db.GetMongoDbCollection(settings.DbName, settings.CollectionNameOfCheckInRecord)

	// 若連線有誤
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	// 預設 filter 為空
	var filter bson.M = bson.M{}

	// 若給參數-指定日期，設定filter，過濾出指定日期的資料
	if c.Params("date") != "" {

		// 取出參數
		parameterDate := c.Params("date")
		fmt.Println("查詢日期為 ", parameterDate)

		// 若日期為今天，則進行時間判斷
		if isToday(parameterDate) {

			//現在
			now := time.Now().Local() //In(time.FixedZone("", 8*60*60))

			//今天
			today := time.Date(
				now.Year(),
				time.Month(now.Month()),
				now.Day(),
				0,
				0,
				0,
				0,
				time.Local,
			)

			fmt.Println("now=", now)     //filter 型態 Map[date:2020-01-01]
			fmt.Println("today=", today) //filter 型態 Map[date:2020-01-01]
			// 過濾器
			filter = bson.M{
				"datetimetoday":   today,
				"datetimecheckin": bson.M{`$lt`: now},
			}
			fmt.Println("filter=", filter) //filter 型態 Map[date:2020-01-01]

		} else {

			// 若日期非今天，則直接設定filter

			// 過濾器
			filter = bson.M{"date": parameterDate}
			fmt.Println("filter=", filter) //filter 型態 Map[date:2020-01-01]
		}

	}

	// 結果
	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	// 若error
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	// 所有結果裝進result
	cur.All(context.Background(), &results)

	for i, e := range results {
		fmt.Printf("印出第一階段結果：Result[%d]=%s\n", i, e["checkintime"]) //results[0]["checkintime"]
	}

	//最後正確結果
	// correctResult := []primitive.M{}

	//進行晚時間的過濾
	// for i, e := range results {

	// 	fmt.Printf("巡迴結果：Result[%d]=%+v \n", i, e["checkintime"]) //results[0]["checkintime"]

	// 	strTime := fmt.Sprintf("%v", e["checkintime"]) // 轉成string
	// 	strDate := fmt.Sprintf("%v", e["date"])        // 轉成string

	// 	// 加入正確結果:沒請假+是現在時間
	// 	if strTime != "" && !isFutureTime(strDate, strTime) {
	// 		// 正確結果:就加入另外一個
	// 		correctResult = append(correctResult, results[i])
	// 	} else if strTime == "" {
	// 		// 若有病假+事假也要加入前端自己判斷數量
	// 		correctResult = append(correctResult, results[i])
	// 	}
	// }

	// for i, e := range correctResult {
	// 	fmt.Printf("最後結果")
	// 	if e != nil {
	// 		fmt.Printf("最後結果：Result[%d]=%+v \n", i, correctResult[i]) //results[0]["checkintime"]
	// 	}
	// }

	// 若查無資料
	// if correctResult == nil {
	// 	c.SendStatus(404)
	// 	return
	// }

	// json, _ := json.Marshal(correctResult)
	// c.Send(json)
	c.JSON(results)
}

// (優化寫法)取得指定日期<應到>人員資料
func getCheckInRecordYMD(c *fiber.Ctx) {

	// 取得 collection
	collection, err := db.GetMongoDbCollection(settings.DbName, settings.CollectionNameOfCheckInRecord)

	// 若連線有誤
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	// 預設 filter 為空
	var filter bson.M = bson.M{}

	// 取出參數：年月日
	parameterYear := c.Params("year")
	parameterMonth := c.Params("month")
	parameterDay := c.Params("day")

	fmt.Printf("查詢%s年%s月%s日：", parameterYear, parameterMonth, parameterDay)

	// 年月日非空
	if parameterYear != "" && parameterMonth != "" && parameterDay != "" {

		fmt.Printf("轉換沒錯\n")

		// 轉換成數字
		year, errYear := strconv.Atoi(parameterYear)
		month, errMonth := strconv.Atoi(parameterMonth)
		day, errDay := strconv.Atoi(parameterDay)

		fmt.Printf("轉換後 %d年 %d月%d日\n", year, month, day)

		fmt.Printf("轉換後 年err=%s 月err=%s 日err=%s\n", errYear, errMonth, errDay)

		// 轉換數字沒錯誤
		if errYear == nil && errMonth == nil && errDay == nil {

			fmt.Printf("轉換沒錯")

			// 若日期為今天，則要過濾出
			// if isTodayYMD(year, month, day) {

			//今天
			today := time.Date(
				year,
				time.Month(month),
				day,
				0,
				0,
				0,
				0,
				time.Local,
			)

			fmt.Printf("Today %v \n", today)

			//現在時間
			now := time.Now().Local() //.In(time.FixedZone("", 8*60*60))

			// 過濾器
			filter = bson.M{
				"datetimetoday":   today,               //日期為當天
				"datetimecheckin": bson.M{`$lte`: now}, //未來的打卡時間不可出現
			}

			fmt.Println("測試 ", primitive.NewDateTimeFromTime(today).Time()) //filter 型態 Map[date:2020-01-01]
			fmt.Println("filter=", filter)                                  //filter 型態 Map[date:2020-01-01]

			// } else {

			// 	// 若日期非今天，則直接設定filter

			// 	//今天
			// 	today := time.Date(
			// 		year,
			// 		time.Month(month),
			// 		day,
			// 		8, //+8時區 早上8點
			// 		0,
			// 		0,
			// 		0,
			// 		time.Local,
			// 	)

			// 	// 過濾器
			// 	filter = bson.M{"datetimetoday": today}
			// 	fmt.Println("filter=", filter) //filter 型態 Map[date:2020-01-01]
			// }
		}

	} else {
		//年月日有空報錯

		c.Status(500).Send(err)
		return
	}

	// 查詢結果
	var results []bson.M

	// curser
	cur, err := collection.Find(context.Background(), filter)

	// 巡迴結束釋放記憶體
	defer cur.Close(context.Background())

	// 若巡迴error
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	// 查詢結果裝進result
	cur.All(context.Background(), &results)

	//印出所有結果
	for i, e := range results {
		fmt.Printf("印出第一階段結果：Result[%d]=%s\n", i, e["checkintime"]) //results[0]["checkintime"]
	}

	//最後正確結果
	// correctResult := []primitive.M{}

	//進行晚時間的過濾
	// for i, e := range results {

	// 	fmt.Printf("巡迴結果：Result[%d]=%+v \n", i, e["checkintime"]) //results[0]["checkintime"]

	// 	strTime := fmt.Sprintf("%v", e["checkintime"]) // 轉成string
	// 	strDate := fmt.Sprintf("%v", e["date"])        // 轉成string

	// 	// 加入正確結果:沒請假+是現在時間
	// 	if strTime != "" && !isFutureTime(strDate, strTime) {
	// 		// 正確結果:就加入另外一個
	// 		correctResult = append(correctResult, results[i])
	// 	} else if strTime == "" {
	// 		// 若有病假+事假也要加入前端自己判斷數量
	// 		correctResult = append(correctResult, results[i])
	// 	}
	// }

	// for i, e := range correctResult {
	// 	fmt.Printf("最後結果")
	// 	if e != nil {
	// 		fmt.Printf("最後結果：Result[%d]=%+v \n", i, correctResult[i]) //results[0]["checkintime"]
	// 	}
	// }

	// 若查無資料
	// if correctResult == nil {
	// 	c.SendStatus(404)
	// 	return
	// }

	// json, _ := json.Marshal(correctResult)
	// c.Send(json)
	c.JSON(results)
}

func isToday(date string) (result bool) {

	// 判斷是否為今天
	year, _ := strconv.Atoi(date[0:4])
	month, _ := strconv.Atoi(date[5:7])
	day, _ := strconv.Atoi(date[8:10])
	now := time.Now().Local() //.In(time.FixedZone("", 8*60*60))
	return now.Year() == year && (int)(now.Month()) == month && day == now.Day()
}

func isTodayYMD(year int, month int, day int) (result bool) {
	// 判斷是否為今天

	now := time.Now().Local() // 取本機時區 In(time.FixedZone("", 8*60*60))
	return now.Year() == year && (int)(now.Month()) == month && day == now.Day()
}

// // 判斷打卡時間 是否比currentTime更晚
// func isFutureTime(date string, checkintime string) bool {

// 	//現在時間
// 	t := time.Now().In(time.FixedZone("", 8*60*60))
// 	fmt.Println("現在時間=", t)

// 	//拆解打卡年月日
// 	fmt.Println("date=", date)
// 	year, _ := strconv.Atoi(date[0:4])
// 	fmt.Println("年", year)
// 	month, _ := strconv.Atoi(date[5:7])
// 	fmt.Println("月", month)
// 	day, _ := strconv.Atoi(date[8:10])
// 	fmt.Println("日", day)

// 	//拆解打卡時分秒
// 	fmt.Println("checkintime=", checkintime)
// 	hour, _ := strconv.Atoi(checkintime[0:2])
// 	min, _ := strconv.Atoi(checkintime[3:5])
// 	sec, _ := strconv.Atoi(checkintime[6:8])
// 	fmt.Println("時", hour)
// 	fmt.Println("分", min)
// 	fmt.Println("秒", sec)

// 	if year > t.Year() {
// 		return true
// 	} else if year < t.Year() {
// 		return false
// 	} else {
// 		//同年
// 		if month > int(t.Month()) {
// 			return true
// 		} else if month < int(t.Month()) {
// 			return false
// 		} else {
// 			//同月
// 			if day > t.Day() {
// 				return true
// 			} else if day < t.Day() {
// 				return false
// 			} else {
// 				//同日
// 				if hour > t.Hour() {
// 					return true
// 				} else if hour < t.Hour() {
// 					return false
// 				} else {
// 					//同時
// 					if min > t.Minute() {
// 						return true
// 					} else if hour < t.Hour() {
// 						return false
// 					} else {
// 						//同分
// 						if sec > t.Second() {
// 							return true
// 						} else {
// 							return false
// 						}
// 					}
// 				}
// 			}

// 		}
// 	}

// }

// //刪除從資料從一個array
// func remove(array []primitive.M, s int) []primitive.M {

// 	// 當還有元素時
// 	if len(array) > s {
// 		return append(array[:s], array[s+1:]...)
// 	}
// 	emptyArray := []primitive.M{}
// 	return emptyArray
// }

// // 取得指定日期<實到>人員資料
// func getAttendanceOfCheckInStatistics(c *fiber.Ctx) {

// 	// 取得 collection
// 	collection, err := db.GetMongoDbCollection(settings.DbName, settings.CollectionNameOfCheckInRecord)

// 	// 若連線有誤
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	var filter bson.M = bson.M{}

// 	// 若有給date
// 	if c.Params("date") != "" {

// 		/*按照date取出當筆資料*/

// 		//取出date參數
// 		myDate := c.Params("date")

// 		fmt.Println("查詢日期(實到)=", myDate)

// 		//bson.M{} 裡面所用的欄位名稱 必須使用mongoDb欄位名稱 而非struct的欄位名稱 (與JAVA相異)
// 		filter = bson.M{"date": myDate, "leavetype": ""} //應到:leave_type is NULL
// 		fmt.Println("filter=", filter)                   //filter 型態 Map[date:2020-01-01]

// 	}

// 	var results []bson.M
// 	cur, err := collection.Find(context.Background(), filter)
// 	defer cur.Close(context.Background())

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	cur.All(context.Background(), &results)

// 	// 若查無資料
// 	if results == nil {
// 		c.SendStatus(404)
// 		return
// 	}

// 	json, _ := json.Marshal(results)
// 	c.Send(json)
// }

// // 取得指定日期<未到>人員資料
// func getNotArrivedOfCheckInStatistics(c *fiber.Ctx) {

// 	// 取得 collection
// 	collection, err := db.GetMongoDbCollection(settings.DbName, settings.CollectionNameOfCheckInRecord)

// 	// 若連線有誤
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	var filter bson.M = bson.M{}

// 	// 若有給date
// 	if c.Params("date") != "" {

// 		/*按照date取出當筆資料*/

// 		//取出date參數
// 		myDate := c.Params("date")
// 		fmt.Println("查詢日期(未到)=", myDate)

// 		//bson.M{} 裡面所用的欄位名稱 必須使用mongoDb欄位名稱 而非struct的欄位名稱 (與JAVA相異)
// 		filter = bson.M{"date": myDate, "leavetype": bson.M{"$ne": ""}} //應到:leave_type is NOT Equal NULL
// 		fmt.Println("filter=", filter)

// 	}

// 	var results []bson.M
// 	cur, err := collection.Find(context.Background(), filter)
// 	defer cur.Close(context.Background())

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	cur.All(context.Background(), &results)

// 	// 若查無資料
// 	if results == nil {
// 		c.SendStatus(404)
// 		return
// 	}

// 	json, _ := json.Marshal(results)
// 	c.Send(json)
// }

// /* 以下為 CheckInStatistics 相關 functions */
// // 取得指定日期統計資料
// func getCheckInStatistics(c *fiber.Ctx) {

// 	// 取得 collection
// 	collection, err := db.GetMongoDbCollection(settings.DbName, settings.CollectionNameOfCheckInStatistics)

// 	// 若連線有誤
// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	var filter bson.M = bson.M{}

// 	// 若有給date
// 	if c.Params("date") != "" {

// 		/*按照date取出當筆資料*/

// 		//取出date參數
// 		myDate := c.Params("date")
// 		fmt.Println("查詢日期(統計)=", myDate)

// 		filter = bson.M{"date": myDate}
// 		fmt.Println("filter=", filter) //filter 型態 Map[date:2020-01-01]

// 	}

// 	var results []bson.M
// 	cur, err := collection.Find(context.Background(), filter)
// 	defer cur.Close(context.Background())

// 	if err != nil {
// 		c.Status(500).Send(err)
// 		return
// 	}

// 	cur.All(context.Background(), &results)

// 	// 若查無資料
// 	if results == nil {
// 		c.SendStatus(404)
// 		return
// 	}

// 	json, _ := json.Marshal(results)
// 	c.Send(json)
// }
