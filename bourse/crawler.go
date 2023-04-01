package bourse

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"time"

	"github.com/sirupsen/logrus"
)

func Fetch(url string, savePath string, cachedTime time.Duration) (result string, err error) {
	defer func() {
		if err != nil {
			logrus.Debugf("getting url FAILED : %v , saved path: %v", url, savePath)
		} else {
			logrus.Debugf("getting url SUCCEED : %v , saved path: %v", url, savePath)
		}
	}()

	fi, err := os.Stat(savePath)

	mtime := time.Now().Truncate(time.Second * 10000)

	if err != nil {
		//return "", err
	} else {
		mtime = fi.ModTime()
	}

	var byteData []byte

	if time.Now().Before(mtime.Add(cachedTime)) {
		byteData, _ = os.ReadFile(savePath)
	}

	if len(byteData) < 5000 {

		client := &http.Client{Timeout: 5 * time.Second}

		//logrus.Debug("start getting url %v ", url)

		resp, err := client.Get(url)

		if err != nil {
			return "", err
		}

		if resp.StatusCode != http.StatusOK {
			//log.Fatal("crawling failed!")
			return "", errors.New("fetching failed")
		}

		byteData, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			//return "", errors.New("fetching failed")
		}

	}

	pageContent := string(byteData)

	//if savePath == "" {
	//	savePath = filepath.Dir(os.Args[0]) + "/ddd.txt"
	//}

	file, errCreate := os.Create(savePath)
	if errCreate != nil {
		return "", errors.New("error in storage file")
	}
	//helpers.CheckError(errCreate, "check if the path is right in windows and linux")
	defer file.Close()

	n, errWrite := file.Write(byteData)
	if errWrite != nil {
		return "", errors.New("error in storage file")
	}
	//helpers.CheckError(errWrite, "check if the file is writable and isnt locked")

	if n > 0 {
		return pageContent, nil
	}

	return pageContent, errors.New("Action Failed!")
}
