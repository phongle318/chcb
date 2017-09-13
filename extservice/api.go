package extservice

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

func doRequest(url string, isPost bool, postData url.Values) ([]byte, error) {
	var resp *http.Response
	var err error
	start := time.Now()
	log.Debug("Start API call. Url=", url)
	if isPost {
		resp, err = http.PostForm(url, postData)
	} else {
		resp, err = http.Get(url)
	}
	log.Debugf("End API call. Time=%.1f", time.Now().Sub(start).Seconds())
	if err != nil {
		log.Error("Error happened: ", err)
		return nil, err
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Error("Status code: ", resp.StatusCode)
			htmlData, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				log.Error("Body: ", string(htmlData))
			}
			return nil, errors.New("Error happened when calling order API.")
		}
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}
