package main

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func main() {
	var chartPeriod string = "24h"
	var currency = Currency{name: "USD", rate: 1.0, symbol: "$"}
	var cryptoSymbol string
	if len(os.Args) <= 1 {
		printHelp()
	}

	if slices.Contains(os.Args, "-h") || slices.Contains(os.Args, "--help") {
		printHelp()
		os.Exit(0)
	}

	if slices.Contains(os.Args, "-p") || slices.Contains(os.Args, "--period") {
		chartPeriod = getArgValue(os.Args, "-p", "--period")
	}

	cryptoSymbol = os.Args[1]

	var cryptoData, err = getCoinData(cryptoSymbol)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var cryptoChart, chartErr = getChart(cryptoSymbol, chartPeriod)
	if chartErr != nil {
		fmt.Println(chartErr)
		os.Exit(1)
	}

	printCoinData(cryptoData, cryptoChart, currency)
}

func printCoinData(coin Crypto, chartData [][]float64, currency Currency) {
	var style = lipgloss.NewStyle().
		Bold(true)

	fmt.Print(style.Render(coin.Name + " (" + coin.Symbol + "): "))

	var priceColor string = "#067213"

	if coin.PriceChangePercentageDay < 0 {
		priceColor = "#ff0000"
	}

	style = lipgloss.NewStyle().
		Background(lipgloss.Color(priceColor))
	fmt.Println(style.Render(FormatPrice(currency.symbol, coin.Price)))

	printChart(chartData)

	dataTable := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		Width(70).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row%2 == 0 {
				return lipgloss.NewStyle().
					Align(lipgloss.Left).
					MarginLeft(3).Bold(true)

			} else {
				return lipgloss.NewStyle().
					Align(lipgloss.Left).
					MarginLeft(3).
					Bold(false)
			}

		})

	dataTable.Row("Market Cap", FormatPrice(currency.symbol, coin.MarketCap))
	dataTable.Row("Volume", FormatPrice(currency.symbol, coin.Volume))
	dataTable.Row("Max supply", FormatPrice(coin.Symbol, coin.TotalSupply))
	dataTable.Row("Circulating supply", FormatPrice(coin.Symbol, coin.CirculatingSupply))
	dataTable.Row("% of supply in circulation", fmt.Sprintf("%.2f %%", (coin.CirculatingSupply/coin.TotalSupply)*100))
	dataTable.Row("Homepage", coin.Website)

	// style = lipgloss.NewStyle().
	// 	Align(lipgloss.Center)

	fmt.Println(dataTable.String())
}

func printChart(chart [][]float64) {
	var minPrice float64 = chart[0][1]
	var maxPrice float64 = 0

	for i := 0; i < len(chart); i++ {
		if chart[i][1] < minPrice {
			minPrice = chart[i][1]
		} else if chart[i][1] > maxPrice {
			maxPrice = chart[i][1]
		}
	}

	slc := timeserieslinechart.New(100, 15, timeserieslinechart.WithYRange(minPrice, maxPrice))

	for i := 0; i < len(chart); i++ {
		slc.Push(timeserieslinechart.TimePoint{Time: time.Unix(int64(chart[i][0]), 0), Value: chart[i][1]})
	}

	slc.DrawBraille()

	var style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		Margin(1)

	fmt.Println(style.Render(slc.View()))
}

func getArgValue(args []string, shortFlag, longFlag string) string {
	for i, arg := range args {
		if arg == shortFlag || arg == longFlag {
			if i+1 < len(args) {
				return args[i+1]
			}
		}
	}
	return ""
}

func printHelp() {
	fmt.Println(
		`
Usage: tcoin <coin>

Options:
 -h, --help 	 Display this help message
 -p, --period    Period of the chart (24h, 1w, 1m, 3m, 6m, 1y, all)

Example: tcoin bitcoin
Example: tcoin ethereum -p 1y
Example: tcoin monero --period 1m

 `)
	os.Exit(0)

}
