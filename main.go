package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Data []struct {
		ConfirmDate    string      `json:"ConfirmDate"`
		No             string      `json:"No"`
		Age            interface{} `json:"Age"`
		Gender         string      `json:"Gender"`
		GenderEn       string      `json:"GenderEn"`
		Nation         string      `json:"Nation"`
		NationEn       string      `json:"NationEn"`
		Province       string      `json:"Province"`
		ProvinceID     int         `json:"ProvinceId"`
		District       string      `json:"District"`
		ProvinceEn     string      `json:"ProvinceEn"`
		StatQuarantine int         `json:"StatQuarantine"`
	} `json:"Data"`
}

func main() {
	r := gin.Default()
	r.GET("/covid/summary", covidSummary)
	r.Run(":8080")
}

func covidSummary(c *gin.Context) {
	// Call Data
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	response, err := http.Get("https://static.wongnai.com/devinterview/covid-cases.json")
	if err != nil {
		log.Fatal(err)
	}
	//Read Data
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
		c.JSON(200, gin.H{
			"message": "Data source invalid",
		})
	}
	var data ResponseData
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
		c.JSON(200, gin.H{
			"message": "Data invalid",
		})
	}
	//Define count variable
	provinceCount := map[string]int{}
	ageGroupCount := map[string]int{
		"0-30":  0,
		"31-60": 0,
		"61+":   0,
		"N/A":   0,
	}
	//Summary Data//
	for i := 0; i < len(data.Data); i++ {
		if data.Data[i].Age != nil {
			if data.Data[i].Age.(float64) <= float64(30) {
				ageGroupCount["0-30"] = ageGroupCount["0-30"] + 1
			} else if data.Data[i].Age.(float64) <= float64(60) {
				ageGroupCount["31-60"] = ageGroupCount["31-60"] + 1
			} else if data.Data[i].Age.(float64) > float64(60) {
				ageGroupCount["61+"] = ageGroupCount["61+"] + 1
			}
		} else {
			ageGroupCount["N/A"] = ageGroupCount["N/A"] + 1
		}
		if data.Data[i].Province == "" {
			provinceCount["N/A"] = provinceCount["N/A"] + 1
		} else {
			province_name := data.Data[i].Province
			provinceCount[province_name] = provinceCount[province_name] + 1
		}

	}
	//Send response
	c.JSON(200, gin.H{
		"Province": provinceCount,
		"AgeGroup": ageGroupCount,
	})

}
