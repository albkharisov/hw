//go:build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("find empty domain name", func(t *testing.T) {
		_, err := GetDomainStat(bytes.NewBufferString(data), "")
		require.Equal(t, ErrDomain, err)
	})

	t.Run("find uncorrect domain name", func(t *testing.T) {
		_, err := GetDomainStat(bytes.NewBufferString(data), "broWsecat.com")
		require.Equal(t, ErrDomain, err)
	})
}

func TestGetDomainStatErrors(t *testing.T) {
	dataNoEmailField := `{"Id":1,"Name":"Janice Rose","Username":"KeithHart","Email":"abc@def.gov","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":2,"Name":"Howard Mendoza","Username":"0Oliver","mail":"abc@def.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":3,"Name":"Janice Rose","Username":"KeithHart","Email":"abc@hjk.gov","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":4,"Name":"Howard Mendoza","Username":"0Oliver","mail":"abc@def.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}`
	dataEmptyEmailField := `{"Id":1,"Name":"Janice Rose","Username":"KeithHart","Email":"abc@def.gov","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":2,"Name":"Howard Mendoza","Username":"0Oliver","Email":"","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":3,"Name":"Janice Rose","Username":"KeithHart","Email":"abc@hjk.gov","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":4,"Name":"Howard Mendoza","Username":"0Oliver","Email":"","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}`

	t.Run("parse JSON with some 'Email' fields missing", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(dataNoEmailField), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"def.gov": 1,
			"hjk.gov": 1,
		}, result)
	})

	t.Run("parse JSON with empty 'Email fields'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(dataEmptyEmailField), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"def.gov": 1,
			"hjk.gov": 1,
		}, result)
	})
}

func TestGetDomainStatJSONmalformed(t *testing.T) {
	dataJSONemptyLines := `{"Id":1,"Name":"Janice Rose","Username":"KeithHart","Email":"abc@def.gov","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}


{"Id":3,"Name":"Janice Rose","Username":"KeithHart","Email":"abc@hjk.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":4,"Name":"Howard Mendoza","Username":"0Oliver","Email":"abc@def.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}`
	dataMalformedJSON := `{"Id":1,"Name":"Janice Rose","Username":"KeithHart","Email":"abc@def.gov","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
	{"Id":3,"Name":"Janice Rose","Username":"KeithHart","Email":"abc@hjk.gov","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
	"Id":4,"Name":"Howard Mendoza","Username":"0Oliver","mail":"abc@def.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}`

	t.Run("parse JSON with empty lines", func(t *testing.T) {
		_, err := GetDomainStat(bytes.NewBufferString(dataJSONemptyLines), "gov")
		require.Equal(t, ErrMalformedJSON, err)
	})

	t.Run("parse malformed JSON", func(t *testing.T) {
		_, err := GetDomainStat(bytes.NewBufferString(dataMalformedJSON), "gov")
		require.Equal(t, ErrMalformedJSON, err)
	})
}
