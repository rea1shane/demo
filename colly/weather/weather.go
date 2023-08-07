package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"regexp"
	"strconv"
	"time"
)

type TodayDetails struct {
	TempHigh      int     // 最高温度
	TempLow       int     // 最低温度
	WindDirection int     // 风向
	WindSpeed     int     // 风速
	Humidity      int     // 湿度
	DewPoint      int     // 露点
	Pressure      float64 // 气压
	UVIndex       int     // 紫外线指数
	Visibility    float64 // 能见度
	MoonPhase     string  // 月相

	CollectTime time.Time // 采集时间
}

var intRe = regexp.MustCompile(`[0-9]+`)
var floatRe = regexp.MustCompile(`\d+\.\d+`)

func main() {
	c := colly.NewCollector()
	collectTodayDetails(c)
	c.Visit("https://weather.com/zh-CN/weather/today/l/7a4684e0789c881e79935986f2e9e5ab05b0104ac4310fd8818006dfb66092c3")
}

func collectTodayDetails(c *colly.Collector) {
	c.OnHTML(".TodayDetailsCard--detailsContainer--16Hg0", func(e *colly.HTMLElement) {
		todayDetails := &TodayDetails{}
		e.ForEach("div.ListItem--listItem--2wQRK", func(i int, elem *colly.HTMLElement) {
			switch i {
			// Number
			case 0, 1, 2, 3, 4, 5, 6:
				spanData := elem.ChildText("div > span")
				switch i {
				// Int
				case 0, 1, 2, 3, 5:
					spanInts := findInts(spanData)
					switch i {
					// Temp
					case 0:
						switch len(spanInts) {
						default:
							todayDetails.TempHigh = spanInts[0]
							todayDetails.TempLow = spanInts[1]
						case 1:
							todayDetails.TempHigh = -999
							todayDetails.TempLow = spanInts[0]
						}
					// Wind
					case 1:
						windDirectionStr := elem.ChildAttr("div > span > svg", "style")
						todayDetails.WindDirection = findInts(windDirectionStr)[0]
						todayDetails.WindSpeed = spanInts[0]
					// Humidity
					case 2:
						todayDetails.Humidity = spanInts[0]
					// DewPoint
					case 3:
						todayDetails.DewPoint = spanInts[0]
					// UVIndex
					case 5:
						todayDetails.UVIndex = spanInts[0]
					}
				// Float
				case 4, 6:
					spanFloats := findFloats(spanData)
					switch i {
					// Pressure
					case 4:
						todayDetails.Pressure = spanFloats[0]
					// Visibility
					case 6:
						todayDetails.Visibility = spanFloats[0]
					}
				}
			// String
			case 7:
				data := elem.ChildText("div.WeatherDetailsListItem--wxData--2s6HT")
				todayDetails.MoonPhase = data
			}
		})
		todayDetails.CollectTime = time.Now()
		fmt.Println(todayDetails)
	})
}

func findInts(str string) []int {
	intStrings := intRe.FindAllString(str, -1)
	res := make([]int, len(intStrings))
	for i, str := range intStrings {
		res[i], _ = strconv.Atoi(str)
	}
	return res
}

func findFloats(str string) []float64 {
	intStrings := floatRe.FindAllString(str, -1)
	res := make([]float64, len(intStrings))
	for i, str := range intStrings {
		res[i], _ = strconv.ParseFloat(str, 64)
	}
	return res
}
