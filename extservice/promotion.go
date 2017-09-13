package extservice

import (
	"encoding/json"
	"fmt"
	"github.com/fpt-corp/fptshop/config"
)

type Promotion struct {
	Title       string `json:"Name"`
	Url         string `json:"Link"`
	Image       string `json:"PictureUrl"`
	Description string
}

func (p Promotion) GetTitle() string       { return p.Title }
func (p Promotion) GetUrl() string         { return p.Url }
func (p Promotion) GetImage() string       { return p.Image }
func (p Promotion) GetDescription() string { return p.Description }

func GetHomepagePromotions() ([]Promotion, error) {
	url := config.Env.OrderApiUrl + fmt.Sprintf("/Order/GetListBannerHome?key=%s", config.Env.OrderApiKey)
	body, err := doRequest(url, false, nil)
	if err != nil {
		return nil, err
	}

	var promotions []Promotion
	err = json.Unmarshal(body, &promotions)
	if err != nil {
		return nil, err
	}

	return promotions, nil
}
