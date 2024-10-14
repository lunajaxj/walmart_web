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

	for i := 0; i < len(inventory); i++ {
		if i == len(inventory)-1 {
			//第一天没有以往数据，因此将销售额设为缺失数据
			sales[i] = -1
		} else {
			prevQuantity := inventory[i-1].Quantity
			currQuantity := inventory[i].Quantity

			if currQuantity == 999999999 {
				// 库存数量无效，将销售标记为错误
				sales[i] = -1
			} else if prevQuantity == 999999999 {
				// 如果前一天的数量无效，则搜索最后有效的数量
				j := i - 1
				for j >= 0 && inventory[j].Quantity == 999999999 {
					j--
				}

				if j < 0 {
					// 未找到有效的先前数据，将销售额标记为缺失数据
					sales[i] = -1
				} else {
					prevQuantity = inventory[j].Quantity
					sales[i] = prevQuantity - currQuantity
					if sales[i] < 0 {
						sales[i] = 0
					}
				}
			} else {
				sales[i] = prevQuantity - currQuantity
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
