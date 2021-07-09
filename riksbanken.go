package riksbanken

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CrossName struct {
	Description string `xml:"seriesdescription"`
	ID          string `xml:"seriesid"`
	Name        string `xml:"seriesname"`
}

type SEKPair struct {
	Date  string `xml:"date"`
	Value string `xml:"value"`
}

func GetAllCrossNames() ([]CrossName, error) {
	type Response struct {
		Return []CrossName `xml:"return"`
	}

	type Body struct {
		Response Response `xml:"http://swea.riksbank.se/xsd getAllCrossNamesResponse"`
	}

	type Envelope struct {
		Body Body `xml:"http://www.w3.org/2003/05/soap-envelope Body"`
	}

	client := &http.Client{}
	payload := []byte(`<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope" xmlns:xsd="http://swea.riksbank.se/xsd">
	<soap:Header/>
	<soap:Body>
		<xsd:getAllCrossNames>
			<languageid>en</languageid>
		</xsd:getAllCrossNames>
	</soap:Body>
</soap:Envelope>
`)
	req, err := http.NewRequest("POST", "http://swea.riksbank.se/sweaWS/services/SweaWebServiceHttpSoap12Endpoint", bytes.NewBuffer(payload))

	if err != nil {
		return []CrossName{}, err
	}
	req.Header.Set("Content-Type", "application/soap+xml;charset=utf-8;action=urn:getAllCrossNames")

	resp, err := client.Do(req)
	if err != nil {
		return []CrossName{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	cn := Envelope{}
	xml.Unmarshal(body, &cn)

	return cn.Body.Response.Return, nil
}

func GetInterestAndExchangeRates() ([]SEKPair, error) {
	fromDate := "2021-06-01"
	toDate := "2021-06-30"

	type Series struct {
		ID      string    `xml:"seriesid"`
		Name    string    `xml:"seriesname"`
		Results []SEKPair `xml:"resultrows"`
	}
	type Groups struct {
		ID     int    `xml:"groupid"`
		Name   string `xml:"groupname"`
		Series Series `xml:"series"`
	}
	type Return struct {
		DateFrom string `xml:"datefrom"`
		DateTo   string `xml:"dateto"`
		Groups   Groups `xml:"groups"`
	}
	type Response struct {
		Return Return `xml:"return"`
	}
	type Body struct {
		Response Response `xml:"http://swea.riksbank.se/xsd getInterestAndExchangeRatesResponse"`
	}

	type Envelope struct {
		Body Body `xml:"http://www.w3.org/2003/05/soap-envelope Body"`
	}

	client := &http.Client{}
	payload := []byte(fmt.Sprintf(`<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope" xmlns:xsd="http://swea.riksbank.se/xsd">
	<soap:Header/>
	<soap:Body>
		<xsd:getInterestAndExchangeRates>
			<searchRequestParameters>
				<aggregateMethod>D</aggregateMethod>
				<datefrom>%s</datefrom>
				<dateto>%s</dateto>
				<languageid>en</languageid>
				<min>false</min>
				<avg>true</avg>
				<max>true</max>
				<ultimo>false</ultimo>
				<searchGroupSeries>
					<groupid>11</groupid>
					<seriesid>SEKEURPMI</seriesid>
				</searchGroupSeries>
			</searchRequestParameters>
		</xsd:getInterestAndExchangeRates>
	</soap:Body>
</soap:Envelope>
`, fromDate, toDate))

	req, err := http.NewRequest("POST", "http://swea.riksbank.se/sweaWS/services/SweaWebServiceHttpSoap12Endpoint", bytes.NewBuffer(payload))

	if err != nil {
		return []SEKPair{}, err
	}
	req.Header.Set("Content-Type", "application/soap+xml;charset=utf-8;action=urn:getInterestAndExchangeRates")

	resp, err := client.Do(req)
	if err != nil {
		return []SEKPair{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []SEKPair{}, err
	}

	cn := Envelope{}
	xml.Unmarshal(body, &cn)

	return cn.Body.Response.Return.Groups.Series.Results, nil
}
