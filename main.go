package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

var FetchExchangeRates = fetchExchangeRates

var amountToConvert float64

// Global variable to store exchange rates
var exchangeRates map[string]string

// Structs to match the JSON structure
type ExchangeRates struct {
	Data RatesData `json:"data"`
}

type RatesData struct {
	Currency string            `json:"currency"`
	Rates    map[string]string `json:"rates"`
}

type CurrencyConversionPair struct {
	CryptoName  string
	PctOfAmount float64
}

func (ccp CurrencyConversionPair) ConvertToCrypto() (float64, error) {
	rate, err := strconv.ParseFloat(exchangeRates[ccp.CryptoName], 64)
	if err != nil {
		return 0, err
	}

	res := amountToConvert * (ccp.PctOfAmount / float64(100)) * rate
	return res, nil

}

func (ccp CurrencyConversionPair) amountToConvertPCT() float64 {
	return amountToConvert * (ccp.PctOfAmount / 100)
}
func main() {

	err := FetchExchangeRates()
	if err != nil {
		log.Fatal("there was a problem with the API")
	}

	app := GetApp()

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var GetApp = getApp

func getApp() *cli.App {
	app := &cli.App{
		Name:  "CLI Application",
		Usage: "A simple CLI application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "strings",
				Value:   "BTC,ETH",
				Usage:   "a comma-separated list of strings",
				Aliases: []string{"s"},
			},
			&cli.StringFlag{
				Name:    "numbers",
				Value:   "70,30",
				Usage:   "a comma-separated list of numbers that must equal 100",
				Aliases: []string{"n"},
			},
		},
		Action: func(c *cli.Context) error {
			command := c.Args().Get(0)
			if command == "" {
				return fmt.Errorf("a command is required")
			}
			var erramt error
			amountToConvert, erramt = strconv.ParseFloat(c.Args().Get(0), 64)

			if erramt != nil {
				return erramt
			}

			stringList := c.String("strings")
			numberList := c.String("numbers")

			currencyPairs, err := ValidateAnProcess(numberList, stringList)
			if err != nil {
				return err
			}

			err = PrintCurrencyPairs(currencyPairs)

			if err != nil {
				return err
			}

			return nil
		},
	}
	return app
}

func PrintCurrencyPairs(currencyPairs []CurrencyConversionPair) error {
	for _, pair := range currencyPairs {
		convert, err := pair.ConvertToCrypto()

		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\n$%v => %v %v\n", pair.amountToConvertPCT(), convert, pair.CryptoName)
	}
	return nil
}

func ValidateAnProcess(numberList string, stringList string) ([]CurrencyConversionPair, error) {
	// Validate and process numbers
	numbers := strings.Split(numberList, ",")
	sum := 0
	for _, numStr := range numbers {
		num, err := strconv.Atoi(strings.TrimSpace(numStr))
		if err != nil {
			return nil, fmt.Errorf("invalid number: %v", numStr)
		}
		sum += num
	}

	if sum != 100 {
		return nil, fmt.Errorf("numbers must sum up to 100")
	}

	stringSlice := strings.Split(stringList, ",")
	numberSlice := strings.Split(numberList, ",")

	currencyPairs, err := MatchAndCreateStructs(stringSlice, numberSlice)

	if err != nil {
		log.Fatal(err)
	}
	return currencyPairs, nil
}

// Function to match strings and numbers by index and return a slice of structs
func MatchAndCreateStructs(stringsList []string, numbersList []string) ([]CurrencyConversionPair, error) {
	if len(stringsList) != len(numbersList) {
		return nil, fmt.Errorf("strings and numbers lists must be of the same length")
	}

	var pairs []CurrencyConversionPair

	for i, s := range stringsList {
		num, err := strconv.ParseFloat(numbersList[i], 64)
		if err != nil {
			return nil, fmt.Errorf("error converting number at index %d: %v", i, err)
		}

		pair := CurrencyConversionPair{CryptoName: s, PctOfAmount: num}
		pairs = append(pairs, pair)
	}

	return pairs, nil
}

// Function to fetch and unmarshal exchange rates
func fetchExchangeRates() error {
	url := "https://api.coinbase.com/v2/exchange-rates?currency=USD"

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	er := &ExchangeRates{}
	err = json.Unmarshal(body, &er)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	exchangeRates = er.Data.Rates

	return nil
}
