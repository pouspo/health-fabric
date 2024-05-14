package chaincode

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// UniqueStrings Function to make a list of strings unique
func UniqueStrings(input []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueList []string

	// Iterate through the input list
	for _, str := range input {
		// If the string is not already in the map, add it
		if _, ok := uniqueMap[str]; !ok {
			uniqueMap[str] = true
			uniqueList = append(uniqueList, str)
		}
	}

	return uniqueList
}

func inArray(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

// Get the DN (distinguished name) associated with a pkix.Name.
// NOTE: This code is almost a direct copy of the String() function in
// https://go-review.googlesource.com/c/go/+/67270/1/src/crypto/x509/pkix/pkix.go#26
// which returns a DN as defined by RFC 2253.
func getDN(name *pkix.Name) string {
	r := name.ToRDNSequence()
	s := ""
	for i := 0; i < len(r); i++ {
		rdn := r[len(r)-1-i]
		if i > 0 {
			s += ","
		}
		for j, tv := range rdn {
			if j > 0 {
				s += "+"
			}
			typeString := tv.Type.String()
			typeName, ok := attributeTypeNames[typeString]
			if !ok {
				derBytes, err := asn1.Marshal(tv.Value)
				if err == nil {
					s += typeString + "=#" + hex.EncodeToString(derBytes)
					continue // No value escaping necessary.
				}
				typeName = typeString
			}
			valueString := fmt.Sprint(tv.Value)
			escaped := ""
			begin := 0
			for idx, c := range valueString {
				if (idx == 0 && (c == ' ' || c == '#')) ||
					(idx == len(valueString)-1 && c == ' ') {
					escaped += valueString[begin:idx]
					escaped += "\\" + string(c)
					begin = idx + 1
					continue
				}
				switch c {
				case ',', '+', '"', '\\', '<', '>', ';':
					escaped += valueString[begin:idx]
					escaped += "\\" + string(c)
					begin = idx + 1
				}
			}
			escaped += valueString[begin:]
			s += typeName + "=" + escaped
		}
	}
	return s
}

var attributeTypeNames = map[string]string{
	"2.5.4.6":  "C",
	"2.5.4.10": "O",
	"2.5.4.11": "OU",
	"2.5.4.3":  "CN",
	"2.5.4.5":  "SERIALNUMBER",
	"2.5.4.7":  "L",
	"2.5.4.8":  "ST",
	"2.5.4.9":  "STREET",
	"2.5.4.17": "POSTALCODE",
}

// Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}

func getGroupAttr(parsedCert *x509.Certificate) (string, error) {
	type T struct {
		Attrs struct {
			Groups string `json:"groups"`
		} `json:"attrs"`
	}

	var t T
	var extValue []byte

	// Finding the specified extension
	for _, ext := range parsedCert.Extensions {
		if ext.Id.String() == "1.2.3.4.5.6.7.8.1" {
			extValue = ext.Value
		}
	}

	if len(extValue) > 0 {
		if err := json.NewDecoder(bytes.NewReader(extValue)).Decode(&t); err != nil {
			return "", err
		}
	}

	return t.Attrs.Groups, nil
}
