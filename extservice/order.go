package extservice

import (
	"fmt"
	"net/url"

	"github.com/fpt-corp/fptshop/config"
)

func CreateOrder(productID string, phone string, customerName string, variantID string, shopName string) error {
	phone = url.QueryEscape(phone)
	customerName = url.QueryEscape(customerName)
	shopName = url.QueryEscape(shopName)
	postUrl := config.Env.OrderApiUrl + fmt.Sprintf("/addorder?proid=%s&name=%s&phone=%s&shop=%s&variantid=%s", productID, customerName, phone, shopName, variantID)
	postData := url.Values{
		"prodid":    {productID},
		"name":      {customerName},
		"phone":     {phone},
		"variantid": {variantID},
		"shop":      {shopName},
		"key":       {config.Env.OrderApiKey}}

	_, err := doRequest(postUrl, true, postData)
	return err
}
