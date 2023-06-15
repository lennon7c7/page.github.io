package util

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
)

type MapResponse struct {
	Info struct {
		Type    int    `json:"type"`
		Error   int    `json:"error"`
		Time    int    `json:"time"`
		Message string `json:"message"`
	} `json:"info"`
	Detail struct {
		PoiCount int `json:"poi_count"`
		Poilist  []struct {
			Addr     string `json:"addr"`
			AddrInfo struct {
				Adcode    string `json:"adcode"`
				C         string `json:"c"`
				D         string `json:"d"`
				P         string `json:"p"`
				ShortAddr string `json:"short_addr"`
			} `json:"addr_info"`
			BaseMapInfo struct {
				BaseMapFlag string `json:"base_map_flag"`
			} `json:"base_map_info"`
			Catacode   string        `json:"catacode"`
			Catalog    string        `json:"catalog"`
			Direction  string        `json:"direction"`
			Dist       string        `json:"dist"`
			Dtype      string        `json:"dtype"`
			FinalScore string        `json:"final_score"`
			Name       string        `json:"name"`
			Pointx     string        `json:"pointx"`
			Pointy     string        `json:"pointy"`
			Tags       []interface{} `json:"tags"`
			UID        string        `json:"uid"`
			Weight     string        `json:"weight"`
		} `json:"poilist"`
		RequestID string `json:"request_id"`
		Results   []struct {
			Adcode               string `json:"adcode,omitempty"`
			C                    string `json:"c,omitempty"`
			CCht                 string `json:"c_cht,omitempty"`
			CEn                  string `json:"c_en,omitempty"`
			CityCode             string `json:"city_code,omitempty"`
			D                    string `json:"d,omitempty"`
			DCht                 string `json:"d_cht,omitempty"`
			DEn                  string `json:"d_en,omitempty"`
			Dtype                string `json:"dtype"`
			N                    string `json:"n,omitempty"`
			NCht                 string `json:"n_cht,omitempty"`
			NEn                  string `json:"n_en,omitempty"`
			Name                 string `json:"name,omitempty"`
			NationCode           string `json:"nation_code,omitempty"`
			P                    string `json:"p,omitempty"`
			PCht                 string `json:"p_cht,omitempty"`
			PEn                  string `json:"p_en,omitempty"`
			PhoneAreaCode        string `json:"phone_area_code,omitempty"`
			Pointx               string `json:"pointx,omitempty"`
			Pointy               string `json:"pointy,omitempty"`
			Stat                 int    `json:"stat,omitempty"`
			AddressChildrenScene string `json:"address_children_scene,omitempty"`
			AddressName          string `json:"address_name,omitempty"`
			AddressScene         string `json:"address_scene,omitempty"`
			RoughAddressName     string `json:"rough_address_name,omitempty"`
			SecondLandmark       struct {
				Addr          string `json:"addr"`
				AddressCht    string `json:"address_cht"`
				AddressEn     string `json:"address_en"`
				Catacode      string `json:"catacode"`
				Direction     string `json:"direction"`
				Dist          string `json:"dist"`
				Dtype         string `json:"dtype"`
				LandmarkLevel string `json:"landmark_level"`
				Name          string `json:"name"`
				Pointx        string `json:"pointx"`
				Pointy        string `json:"pointy"`
				UID           string `json:"uid"`
			} `json:"second_landmark,omitempty"`
			DescWeight         string `json:"desc_weight,omitempty"`
			Direction          string `json:"direction,omitempty"`
			Dist               string `json:"dist,omitempty"`
			UID                string `json:"uid,omitempty"`
			Street             string `json:"street,omitempty"`
			StreetNumber       string `json:"street_number,omitempty"`
			Kind               string `json:"kind,omitempty"`
			Tag                string `json:"tag,omitempty"`
			StandardAddress    string `json:"standard_address,omitempty"`
			StandardAddressCht string `json:"standard_address_cht,omitempty"`
			StandardAddressEn  string `json:"standard_address_en,omitempty"`
		} `json:"results"`
	} `json:"detail"`
}

// 经纬度 转 省市区地址
func MapXyToAddress(longitude float32, latitude float32) (Province string, City string, Area string, Community string) {
	if longitude <= 0 || latitude <= 0 {
		return
	}

	url := fmt.Sprintf("https://apis.map.qq.com/jsapi?qt=rgeoc&lnglat=%v,%v", longitude, latitude)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		return
	}
	utf8Reader := transform.NewReader(resp.Body,
		simplifiedchinese.GBK.NewDecoder())
	all, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		panic(err)
	}

	var result MapResponse
	err = json.Unmarshal(all, &result)
	if err != nil {
		log.Println("Unmarshal err: ", err)
		return
	}

	for _, value := range result.Detail.Results {
		if value.P != "" && value.C != "" && value.D != "" {
			Province = value.P
			City = value.C
			Area = value.D
		}

		if value.StandardAddress != "" {
			Community = value.StandardAddress
		}
	}

	return
}
