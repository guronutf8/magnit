package Url

import (
	"magnit/internal/DB"
	"math"
	"net/http"
	neturl "net/url"
	"regexp"
)

// Данные кодирования короткой ссылки
var codes = map[string]int{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16, "h": 17, "i": 18, "j": 19, "k": 20, "l": 21, "m": 22, "n": 23, "o": 24, "p": 25, "q": 26, "r": 27, "s": 28, "t": 29, "u": 30, "v": 31, "w": 32, "x": 33, "y": 34, "z": 35, "A": 36, "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42, "H": 43, "I": 44, "J": 45, "K": 46, "L": 47, "M": 48, "N": 49, "O": 50, "P": 51, "Q": 52, "R": 53, "S": 54, "T": 55, "U": 56, "V": 57, "W": 58, "X": 59, "Y": 60, "Z": 61}

const digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Url struct {
	Short string
	Long  string
}

// Save Колдуем, сохраняем оригинальную ссылку, ID записи превращаем в короткий код. Отдаем его
func (uM Url) Save(db *DB.DB) (string, int) {
	uNet, err := neturl.Parse(uM.Long)
	if err != nil {
		panic(err)
	}

	if !uM.checkOriginal(uNet) {
		return "", http.StatusBadRequest
	}

	lastId := db.AddUrl(uNet.Scheme+"://"+uNet.Host, uNet.Path)

	short := uM.formatToShort(uint64(lastId))

	return uNet.Scheme + "://localhost/" + short, 200
}

// Get колдуем, берем которую ссылку и превращаем ее в ID, что бы быстро найти в БД
func (uM Url) Get(db *DB.DB) (string, int) {
	uNet, err := neturl.Parse(uM.Short)
	if err != nil {
		panic(err)
	}

	if !uM.checkShort(uNet) {
		return "", http.StatusBadRequest
	}

	short := []rune(uNet.Path)
	index := uM.formatToIndex(string(short[1:]))
	original, err := db.GetOriginal(index)
	if err != nil {
		return "", http.StatusNotFound
	}

	return original, http.StatusMovedPermanently

}

func (uM Url) checkShort(uNet *neturl.URL) bool {
	if uNet.Scheme != "http" && uNet.Scheme != "https" {
		return false
	}
	if uNet.Host != "localhost" || uNet.Path == "" || uNet.Path == "/" {
		return false
	}

	//проверка, если найдем символы не a-Z 0-9, то ошибка
	matched, _ := regexp.MatchString(`\W`, uNet.Path[1:])
	if matched {
		return false
	}
	return true
}
func (uM Url) checkOriginal(uNet *neturl.URL) bool {
	// todo не уверен как отработают домены .рф
	if uNet.Scheme != "http" && uNet.Scheme != "https" {
		return false
	}
	if uNet.Host == "" || uNet.Path == "" || uNet.Path == "/" {
		return false
	}
	return true
}

// todo max 200000
/**
Вся соль задачи, превращение строки в число
*/
func (uM Url) formatToIndex(s string) (u uint64) {

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		sign := string(runes[len(runes)-i-1])
		if i > 1 {
			a := (i * 61 * codes[sign]) + codes[sign]

			for i2 := 1; i2 <= codes[sign]; i2++ {
				i3 := int(math.Pow(61, 2))
				a += i3
			}
			u += uint64(a)
		} else if i == 1 {

			a := codes[sign]*61 + codes[sign]
			u += uint64(a)
		} else {
			u += uint64(codes[sign])
		}
	}
	return
}

/**
Вся соль, превращение число в строку, гуглить "strconv.formatint", просто улучшенная, т.к. оригинальная могла форматировать лишь в 36тиричную систему
*/
func (uM Url) formatToShort(u uint64) (s string) {
	var a [64]byte
	i := len(a)
	var b uint64 = 62
	for u >= b {
		i--
		q := u / b
		a[i] = digits[uint(u-q*b)]
		u = q
	}
	i--
	a[i] = digits[uint(u)]
	s = string(a[i:])
	return
}
