package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCurrencyConversionPair_ConvertToCrypto(t *testing.T) {
	amountToConvert = 100
	exchangeRates = map[string]string{
		"BTC": "0.0000239383749626",
		"ETH": "0.0004388168619767",
	}

	Convey("TestCurrencyConversionPair_ConvertToCrypto", t, func() {
		amountToConvert = 100
		exchangeRates = map[string]string{
			"BTC": "0.0000239383749626",
			"ETH": "0.0004388168619767",
		}

		Convey("BTC", func() {
			ccp := CurrencyConversionPair{
				CryptoName:  "BTC",
				PctOfAmount: 70,
			}

			got, err := ccp.ConvertToCrypto()

			So(err, ShouldBeNil)
			So(got, ShouldEqual, 0.001675686247382)
		})

		Convey("ETH", func() {
			ccp := CurrencyConversionPair{
				CryptoName:  "ETH",
				PctOfAmount: 30,
			}

			got, err := ccp.ConvertToCrypto()

			So(err, ShouldBeNil)
			So(got, ShouldEqual, 0.013164505859301)
		})
	})
}

func TestCurrencyConversionPair_amountToConvertPCT(t *testing.T) {
	type fields struct {
		CryptoName  string
		PctOfAmount float64
	}
	Convey("TestCurrencyConversionPair_amountToConvertPCT", t, func() {
		Convey("70% @ 100", func() {
			tt := struct {
				name   string
				amt    float64
				fields fields
				want   float64
			}{
				name: "70% @ 100",
				amt:  100.0,
				fields: fields{
					CryptoName:  "BTC",
					PctOfAmount: 70.0,
				},
				want: 70,
			}

			amountToConvert = tt.amt
			ccp := CurrencyConversionPair{
				CryptoName:  tt.fields.CryptoName,
				PctOfAmount: tt.fields.PctOfAmount,
			}
			got := ccp.amountToConvertPCT()
			So(got, ShouldEqual, tt.want)

		})
		Convey("70% @ 200", func() {
			tt := struct {
				name   string
				amt    float64
				fields fields
				want   float64
			}{
				name: "70% @ 200",
				amt:  200.0,
				fields: fields{
					CryptoName:  "BTC",
					PctOfAmount: 70,
				},
				want: 140,
			}

			amountToConvert = tt.amt
			ccp := CurrencyConversionPair{
				CryptoName:  tt.fields.CryptoName,
				PctOfAmount: tt.fields.PctOfAmount,
			}
			got := ccp.amountToConvertPCT()

			So(got, ShouldEqual, tt.want)

		})
		Convey("30% @ 300", func() {
			tt := struct {
				name   string
				amt    float64
				fields fields

				want float64
			}{
				name: "30% @ 300",
				amt:  300.0,
				fields: fields{
					CryptoName:  "BTC",
					PctOfAmount: 30,
				},
				want: 90,
			}

			amountToConvert = tt.amt
			ccp := CurrencyConversionPair{
				CryptoName:  tt.fields.CryptoName,
				PctOfAmount: tt.fields.PctOfAmount,
			}
			got := ccp.amountToConvertPCT()
			So(got, ShouldEqual, tt.want)

		})
	})
}

func TestMatchAndCreateStructs(t *testing.T) {
	type args struct {
		stringsList []string
		numbersList []string
	}

	Convey("TestMatchAndCreateStructs", t, func() {

		tt := struct {
			name    string
			args    args
			want    []CurrencyConversionPair
			wantErr bool
		}{
			name: "BTC,ETh",
			args: args{
				stringsList: []string{"BTC", "ETH"},
				numbersList: []string{"70", "30"},
			},
			want: []CurrencyConversionPair{
				{
					CryptoName:  "BTC",
					PctOfAmount: 70,
				},
				{
					CryptoName:  "ETH",
					PctOfAmount: 30,
				},
			},
			wantErr: false,
		}

		got, err := MatchAndCreateStructs(tt.args.stringsList, tt.args.numbersList)
		So(err, ShouldBeNil)
		So(got, ShouldResemble, tt.want)

	})

}

func Test_fetchExchangeRates(t *testing.T) {
	Convey("Test_fetchExchangeRates", t, func() {

		err := fetchExchangeRates()
		So(err, ShouldBeNil)
	})
}

func TestValidateAnProcess(t *testing.T) {
	Convey("Validate and Process", t, func() {
		Convey("Happy Path", func() {
			want := []CurrencyConversionPair{
				{
					CryptoName:  "BTC",
					PctOfAmount: 70,
				},
				{
					CryptoName:  "ETH",
					PctOfAmount: 30,
				},
			}
			got, err := ValidateAnProcess("70,30", "BTC,ETH")
			So(got, ShouldResemble, want)
			So(err, ShouldBeNil)

		})
	})
}

func TestPrintCurrencyPairs(t *testing.T) {

	Convey("PrintCurrencyPairs", t, func() {
		amountToConvert = 100
		exchangeRates = map[string]string{
			"BTC": "0.0000239383749626",
			"ETH": "0.0004388168619767",
		}

		x := []CurrencyConversionPair{
			{
				CryptoName:  "BTC",
				PctOfAmount: 70,
			},
			{
				CryptoName:  "ETH",
				PctOfAmount: 30,
			},
		}
		err := PrintCurrencyPairs(x)
		So(err, ShouldBeNil)

	})
}
