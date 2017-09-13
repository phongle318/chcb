package text

import (
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dustin/go-humanize"
	"github.com/michlabs/bottext"
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/config"
)

var t bottext.BotTextFunc

func Init(filePath string, locale string) {
	bottext.MustLoad(filePath)
	t = bottext.New(locale)
}

func T(msgKey string, u *fbbot.User) string {
	if u == nil {
		return t(msgKey)
	} else {
		return Personalize(t(msgKey), u)
	}
}

func Personalize(text string, u *fbbot.User) string {
	if strings.Contains(text, "@Gender") || strings.Contains(text, "@gender") {
		switch u.Gender() {
		case "male":
			text = strings.Replace(text, "@Gender", "Anh", -1)
			text = strings.Replace(text, "@gender", "anh", -1)
		case "female":
			text = strings.Replace(text, "@Gender", "Chị", -1)
			text = strings.Replace(text, "@gender", "chị", -1)
		default:
			text = strings.Replace(text, "@Gender", "Bạn", -1)
			text = strings.Replace(text, "@gender", "bạn", -1)
		}
	}

	if strings.Contains(text, "@full_name") {
		text = strings.Replace(text, "@full_name", u.FullName(), -1)
	}

	return text
}

func FromTime(time time.Time) string {
	return strconv.FormatInt(time.Unix(), 10)
}

func ToTime(text string) time.Time {
	intValue, _ := strconv.ParseInt(text, 10, 0)
	return time.Unix(intValue, 0)
}

func ToInt(text string) int {
	if text == "" {
		return 0
	}
	intVal, _ := strconv.Atoi(text)
	return intVal
}

func ToIntSlice(text string) ([]int, error) {
	strVals := strings.Split(text, ",")
	var intVals []int
	for _, id := range strVals {
		intVal, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		intVals = append(intVals, intVal)
	}
	return intVals, nil
}

func FromIntSlice(intSlice []int) string {
	if len(intSlice) == 0 {
		return ""
	}
	var result []string
	for _, value := range intSlice {
		result = append(result, strconv.Itoa(value))
	}
	return strings.Join(result, ",")
}

func FormatPriceVND(price int) string {
	return humanize.FormatInteger("#.###,", price) + "₫"
}

func GetResizedImgURL(imgUrl string) string {
	if strings.HasPrefix(imgUrl, "http") {
		// Normalize url
		imgUrl = strings.TrimPrefix(imgUrl, "http://cdn.fptshop.com.vn/Uploads/Originals/")
		imgUrl = strings.TrimPrefix(imgUrl, "https://cdn.fptshop.com.vn/Uploads/Originals/")
	}
	// Image resizing service should have a "_originals" alias pointing to https://cdn.fptshop.com.vn/Uploads/Originals/
	// Temporary disable for now
	//imgUrl = "_originals/" + imgUrl
	u := &url.URL{Path: imgUrl}
	return fmt.Sprintf(config.Env.ResizedImgUrl, u.String())
}

func GetColorImgURL(colorHex string) string {
	if config.Env.ColorImgUrl == "" {
		return ""
	}
	return fmt.Sprintf(config.Env.ColorImgUrl, strings.TrimPrefix(colorHex, "#"))
}

func Sanitize(htmlStr string) string {
	r, _ := regexp.Compile(`[ \t]*<li>`)
	result := r.ReplaceAllString(htmlStr, "&bull; ") // Replace <li> with bullet character
	r, _ = regexp.Compile(`<[^>]+>`)
	result = r.ReplaceAllString(result, "") // Remove all html tags
	result = html.UnescapeString(result)    // Convert escaped html entities to Unicode character
	r, _ = regexp.Compile(`[\r\n]+`)
	result = r.ReplaceAllString(result, "\n") // Remove multiple line endings
	result = strings.Replace(result, "Đặt Online nhận ngay khuyến mãi:", "", -1)
	result = strings.Replace(result, "Xem chi tiết", "", -1)
	result = strings.Trim(result, "\n")
	return result
}

func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func MarshalToString(data map[string]interface{}) string {
	result, err := json.Marshal(data)
	if err != nil {
		log.Error("Error in FromJSON: ", err)
		return ""
	}
	return string(result)
}
