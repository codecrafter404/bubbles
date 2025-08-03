package graph

func CalculateTotal(customItems []*OrderCustomItem, items []*OrderItem) float64 {
	total := 0.0
	for _, x := range items {
		total += x.Item.Price * float64(x.Quantity)
	}

	for _, c := range customItems {
		for _, v := range c.CustomItem.Variants {
			total += v.Price * float64(c.Quantity)
		}
	}
	return total
}
