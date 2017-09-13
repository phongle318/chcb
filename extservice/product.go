package extservice

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/fpt-corp/fptshop/config"
	"github.com/pkg/errors"
)

type ProductSearchResult struct {
	Products []Product `json:"_searchResults"`
}

type Product struct { // For search API
	ID       int    `json:"ID"`
	Name     string `json:"NameProduct"`
	ImageUrl string `json:"UrlPicture"`
	Type     string `json:"TypeNameAscii"`
	UrlKey   string `json:"NameAscii"`
	Price    string `json:"Price"`
	SKU      string `json:"Sku"`
}

type Product2 struct { // For detail API
	ID        int    `json:"ID"`
	Promotion string `json:"Promotion"`
	Name      string `json:"Name"`
	Price     int    `json:"Price"`
}

type ProductDetail struct {
	Product Product2         `json:"product"`
	Variant []ProductVariant `json:"variant"`
}

type ProductVariant struct {
	ID        int    `json:"ID"`
	SKU       string `json:"sku"`
	ColorName string `json:"ColorName"`
	ColorHex  string `json:"ColorValue"`
}

func SearchForProducts(name string) ([]Product, error) {
	var result = new(ProductSearchResult)
	name = url.QueryEscape(name)
	searchUrl := config.Env.SearchApiUrl + fmt.Sprintf("?searchTerm=%s&searchField=NameProduct", name)

	body, err := doRequest(searchUrl, false, nil)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		return result.Products, nil
	}
}

func GetProductDetails(sku string) (product ProductDetail, err error) {
	result := map[string]interface{}{}
	productUrl := config.Env.ProductApiUrl + fmt.Sprintf("?key=%s&sku=%s", config.Env.OrderApiKey, sku)
	body, err := doRequest(productUrl, false, nil)
	if err != nil {
		return product, err
	} else {
		err = json.Unmarshal(body, &result)
		if err != nil {
			return product, err
		}
		data, ok := result["data"]
		if !ok {
			return product, errors.New("Error when getting product details. Sku=" + sku)
		}
		dataJson, err := json.Marshal(data)
		if err != nil {
			return product, err
		}
		err = json.Unmarshal(dataJson, &product)
		if err != nil {
			return product, err
		}
		return product, nil
	}
}
