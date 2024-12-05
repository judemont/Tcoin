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
	var currency = Currency{name: "USD", rate: 1.0, symbol: "$"}
	var cryptoSymbol string
	if len(os.Args) <= 1 {
		printHelp()
	}

	if slices.Contains(os.Args, "-h") || slices.Contains(os.Args, "--help") {
		printHelp()
		os.Exit(0)
	}

	cryptoSymbol = os.Args[1]

	var cryptoData, err = getCoinData(cryptoSymbol)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var cryptoChart, chartErr = getChart(cryptoSymbol, "6m")
	if chartErr != nil {
		fmt.Println(chartErr)
		os.Exit(1)
	}

	printCoinData(cryptoData, cryptoChart, currency)
}

func printCoinData(coin Crypto, chartData [][]float64, currency Currency) {
	var style = lipgloss.NewStyle().
		Bold(true).
		MarginTop(1)

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

	fmt.Println(slc.View())
}

func printHelp() {
	fmt.Println(
		`
Usage: tcoin <crypto>

Options:
 -h, --help 	 Display this help message

Example: tcoin bitcoin
Example: tcoin eth

 `)
	os.Exit(0)

}
