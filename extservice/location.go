package extservice

import (
	"encoding/json"
	"fmt"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/fpt-corp/fptshop/config"
)

type ShopLocation struct {
	ShopName string  `json:"TenShop"`
	Location string  `json:"DiaChi"`
	Province string  `json:"TenTinh"`
	District string  `json:"TenQuan"`
	Qty      float64 `json:"SoLuongConLai"`
}

// Never return error since FPTShop API is unreliable
func GetNearestShopHasSufficientStock(address string, sku string) []ShopLocation {
	address = url.QueryEscape(address)
	url := config.Env.OrderApiUrl + fmt.Sprintf("/Order/GetListShopByAddress?desc=%s&sku=%s&key=%s", address, sku, config.Env.OrderApiKey)
	body, err := doRequest(url, false, nil)
	locations := []ShopLocation{}
	if err != nil {
		return locations
	}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		log.Error("Error when unmarshaling shop locations ", err)
	}
	var shopHasSufficientStock []ShopLocation
	for _, shop := range locations {
		if int(shop.Qty) >= config.Env.InStockThreshold {
			shopHasSufficientStock = append(shopHasSufficientStock, shop)
		}
	}
	return shopHasSufficientStock
}
