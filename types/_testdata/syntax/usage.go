package fixture

func Use(ch chan<- Number) (Number, error) {
	r := Record{Name: "x", Values: []Number{1, 2, 3}}
	total := Number(0)
outer:
	for i, value := range r.Values {
		switch {
		case i == 1:
			continue
		case value > 2:
			break outer
		default:
			total += value
		}
	}
	defer func() { Global[r.Name] = total }()
	select {
	case ch <- total:
	default:
	}
	return total, nil
}

func Builtins() int {
	values := make([]int, 0)
	values = append(values, 1)
	clear(values)
	return min(max(len(values), 1), 2)
}
