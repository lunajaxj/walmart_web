package main

import (
	"fmt"
)

type Inventory struct {
	Date     string
	Quantity int
}

func calculateSales(inventory []Inventory) []int {
	sales := make([]int, len(inventory))

	// Start the loop from the end (most recent date)
	for i := len(inventory) - 1; i >= 0; i-- {
		// Calculate sales for the last 10 days only
		if len(inventory)-i > 10 {
			break
		}

		if i == len(inventory)-1 {
			// For the last day, no next data, so set sales as missing data
			sales[i] = -1
		} else {
			nextQuantity := inventory[i+1].Quantity
			currQuantity := inventory[i].Quantity

			if currQuantity == 999999999 {
				// Invalid inventory quantity, mark sales as error
				sales[i] = -1
			} else if nextQuantity == 999999999 {
				// If next day's quantity is invalid, search for the last valid quantity
				j := i + 1
				for j < len(inventory) && inventory[j].Quantity == 999999999 {
					j++
				}

				if j == len(inventory) {
					// No valid next data found, mark sales as missing data
					sales[i] = -1
				} else {
					nextQuantity = inventory[j].Quantity
					sales[i] = nextQuantity - currQuantity
					if sales[i] < 0 {
						sales[i] = 0
					}
				}
			} else {
				sales[i] = nextQuantity - currQuantity
				if sales[i] < 0 {
					sales[i] = 0
				}
			}
		}
	}

	return sales
}

func main() {
	// Example usage:
	inventory := []Inventory{
		{"2023-07-16", 999999999},
		{"2023-07-17", 100},
		{"2023-07-18", 200},
		{"2023-07-19", 80},
		{"2023-07-21", 999999999},
		{"2023-07-20", 999999999},
		{"2023-07-22", 50},
	}

	sales := calculateSales(inventory)

	for i, s := range sales {
		if s == -1 {
			fmt.Printf("Day %s: 缺少计算数据\n", inventory[i].Date)
		} else {
			fmt.Printf("Day %s: %d\n", inventory[i].Date, s)
		}
	}
}
