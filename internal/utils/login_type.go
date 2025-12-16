package utils

import (
	"fmt"
	"regexp"
)

var EmailRegex = regexp.MustCompile(`@.*\.`)
var PhoneRegex = regexp.MustCompile(`[\d\-\+\(\)\s]{8,}`)

func DetermineLoginType(login string) (string, error) {
	if IsEmail(login) {
		return "email", nil
	}

	if IsPhoneNumber(login) {
		return "phone", nil
	}

	return "", fmt.Errorf("unknown login type")
}

func IsEmail(login string) bool {
	return EmailRegex.MatchString(login)
}

func IsPhoneNumber(login string) bool {
	return PhoneRegex.MatchString(login)
}

// import (
// 	"github.com/nyaruka/phonenumbers"
// )

// type PhoneData struct {
// 	phoneNumber           string
// 	normalizedPhoneNumber string
// 	valid                 bool
// }

// func NewPhoneData(phoneNumber string) {

// }

// func (parser *PhoneParser) ParsePhoneNumber(phoneNumber string) *PhoneData {
// 	num, err := phonenumbers.Parse(parser.phoneNumber, "RU")

// 	if err != nil {
// 		return &PhoneData{valid: true, phoneNumber: phoneNumber,normalizedPhoneNumber: num.}
// 	} else {
// 		parser.valid = false
// 	}

// 	parser.normalizedPhoneNumber = num.String()

// 	return true
// }

// num, err := phonenumbers.Parse("+79161234567", "RU")
// if err != nil {
// 		panic(err)
// }
