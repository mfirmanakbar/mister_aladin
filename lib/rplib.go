package lib

import (
	"reflect"
	"strings"
	"time"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"

	log "test/logging"

)

/*
	Library created by Raditya Pratama
	for everything
*/

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func EncryptText(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func DecryptText(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func encryptFile(filename string, data []byte, passphrase string) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(EncryptText(data, passphrase))
}

func decryptFile(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return DecryptText(data, passphrase)
}

func (ftp FtpData) openFtp(url string) (*goftp.FTP, error) {
	if ftp.Dt, ftp.err = goftp.Connect(url + ":21"); ftp.err != nil {
		log.Errorf("FTP not connected, check your connection")
		return nil, ftp.err
		// panic(err)
	}

	return ftp.Dt, nil
}
func (ftp FtpData) auth(uname, passwd string) error {
	if ftp.err = ftp.Dt.Login(uname, passwd); ftp.err != nil {
		// panic(err)
		log.Errorf("Account not Authorized")
		return ftp.err
	}
	return nil
}

func (ftp FtpData) ConnectFtp(url, uname, passwd string) (*goftp.FTP, error) {
	myFtp, err := ftp.openFtp(url)
	if err != nil {
		return nil, err
	}
	defer ftp.closeFtp()
	log.Logf("Connected")
	err = ftp.auth(uname, passwd)
	if err != nil {
		return nil, err
	}
	log.Logf("Auth")

	return myFtp, nil
}

func (ftp FtpData) closeFtp() error {
	return ftp.Dt.Close()
}

/*
	Example of using StringToDate Function
	Created by Raditya Pratama
	ex :
		GetDate()
		return : current timestamp with format YYYY-mm-dd HH:ii:ss

		GetDate("yyyy-mm-dd") || StringToDate("y-m-d")
		return : current timestamp with format yyyy-mm-dd

		GetDate("ymd", "20181231")
		return : 20181231
*/

func GetDate(optionParam ...string) (string, time.Time) {
	// var err error
	location, _ := time.LoadLocation("Asia/Jakarta")

	dateTimeToConvert, formatDateTime := time.Now().In(location), "2006-01-02 15:04:05"
	if isset(optionParam, 0) {
		sTime := optionParam[0]
		formatDateTime = ReplaceFormat(sTime)
	}
	if isset(optionParam, 1) {
		sTime := optionParam[1]
		layoutParse := getParseLayout(sTime)
		if layoutParse == "" {
			log.Logf("Invalid String Date Format %s", sTime)
			return "", time.Time{}
		}
		dateTimeToConvert, _ = time.ParseInLocation(layoutParse, sTime, location)
	}
	return dateTimeToConvert.Format(formatDateTime), dateTimeToConvert
}

// DateDiff is used to get diff between two string date
// created by Raditya Pratama
func DateDiff(dateParam ...string) (map[string]float64) {
	totalParam := len(dateParam)
	if(totalParam < 1 || totalParam > 2){
		log.Logf("Unsufficient Param DateDiff")
		return nil
	}

	defaultFormat := "y-m-d h:i:s"

	
	strTime1, time1 := GetDate(defaultFormat, dateParam[0])
	strTime2, time2 := "", time.Time{}
	if !isset(dateParam, 1) {
		strTime2, time2 = GetDate(defaultFormat)
	} else {
		strTime2, time2 = GetDate(defaultFormat, dateParam[1])
	}

	if strTime1 == "" || strTime2 == "" {
		return nil
	}
	duration := time2.Sub(time1)

	
	var diffResult = map[string]float64{
		"seconds": duration.Seconds(),
		"minutes" : duration.Minutes(),
		"hours" : duration.Hours(),
	}

	diffResult["days"] = diffResult["hours"]/24
	diffResult["mounths"] = diffResult["days"]/30
	diffResult["years"] = diffResult["mounths"]/12
	/*var diff float64
	if reqDuration == "all" {
		
	}else if reqDuration == "s" {
		diff = duration.Seconds()
	} else if reqDuration == "i" {
		diff = duration.Minutes()
	} else if reqDuration == "h" {
		diff = duration.Hours()
	} else if reqDuration == "d" {
		diff = (duration.Hours()/24)
	} else if reqDuration == "m" {
		diff = (duration.Hours()/24/30)
	}
	else if reqDuration == "y" {
		diff = (duration.Hours()/24/30/12)
	}*/

	return diffResult
}

func getParseLayout(dateString string) (layoutParse string) {
	// if time.Parse("2016-01-02", dateString)
	var err error
	var formatLists = []string{
		"2006-01-02 15:04:05",
		"2006-02-01 15:04:05",
		"02-01-2006 15:04:05",
		"01-02-2006 15:04:05",
		"2006-01-02 15:04",
		"2006-02-01 15:04",
		"02-01-2006 15:04",
		"01-02-2006 15:04",

		"2006/01/02 15:04:05",
		"2006/02/01 15:04:05",
		"02/01/2006 15:04:05",
		"01/02/2006 15:04:05",
		"2006/01/02 15:04",
		"2006/02/01 15:04",
		"02/01/2006 15:04",
		"01/02/2006 15:04",

		"02012006",
		"01022006",
		"20060102",
		"20060201",

		"02-01-2006",
		"01-02-2006",
		"2006-01-02",
		"2006-02-01",

		"02/01/2006",
		"01/02/2006",
		"2006/01/02",
		"2006/02/01",

		"02012006 150405",
		"01022006 150405",
		"20060102 150405",
		"20060201 150405",

		"02/01/2006 150405",
		"01/02/2006 150405",
		"2006/01/02 150405",
		"2006/02/01 150405",

		"02012006 1504",
		"01022006 1504",
		"20060102 1504",
		"20060201 1504",

		"02/01/2006 1504",
		"01/02/2006 1504",
		"2006/01/02 1504",
		"2006/02/01 1504",
	}

	for _, format := range formatLists {
		_, err = time.Parse(format, dateString)
		if err == nil {
			layoutParse = format
			break
		}
	}
	return
}

func ReplaceByArr(old []string, new string, source string) string {
	for _, e := range old {
		source = strings.Replace(source, e, new, -1)
	}
	return source
}

func ReplaceFormat(dateFormat string) (formatDateTime string) {
	sFormat := dateFormat
	sFormat = ReplaceByArr([]string{"YYYY", "yyyy", "Y", "y"}, "2006", sFormat)
	sFormat = ReplaceByArr([]string{"MM", "mm", "M", "m"}, "01", sFormat)
	sFormat = ReplaceByArr([]string{"DD", "dd", "D", "d"}, "02", sFormat)
	sFormat = ReplaceByArr([]string{"HH", "hh", "H", "h"}, "15", sFormat)
	sFormat = ReplaceByArr([]string{"II", "ii", "I", "i"}, "04", sFormat)
	sFormat = ReplaceByArr([]string{"SS", "ss", "S", "s"}, "05", sFormat)
	sFormat = ReplaceByArr([]string{"W"}, "Monday", sFormat)
	sFormat = ReplaceByArr([]string{"w"}, "Mon", sFormat)
	sFormat = ReplaceByArr([]string{"F"}, "January", sFormat)
	sFormat = ReplaceByArr([]string{"f"}, "Jan", sFormat)

	formatDateTime = sFormat
	return
}

func DeFormatDate(dateFormat string) (formatDateTime string, reformatDate string) {
	// log.Logf("%s", dateFormat)
	reformatDate = getParseLayout(dateFormat)
	// dateTimeToConvert, _ := time.Parse(reformatDate, dateFormat)
	sFormat := reformatDate

	sFormat = strings.Replace(sFormat, "2006", "yyyy" , -1)
	sFormat = strings.Replace(sFormat, "01", "mm", -1)
	sFormat = strings.Replace(sFormat, "02", "dd", -1)
	sFormat = strings.Replace(sFormat, "15", "hh", -1)
	sFormat = strings.Replace(sFormat, "04", "ii", -1)
	sFormat = strings.Replace(sFormat, "05", "ss", -1)
	
	formatDateTime = sFormat
	// log.Logf("%s", dateFormat, reformatDate, sFormat, dateTimeToConvert)
	return
}

/*
	example of using in_array
	names := []string{"Mary", "Anna", "Beth", "Johnny", "Beth"}
    fmt.Println(in_array("Anna", names)) // results true, 1
    fmt.Println(in_array("Jon", names)) // results false, -1
*/
func in_array(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
