package main

import (
	"flag"
	"fmt"

	distr "github.com/akundu/utilities/statistics/distribution"
)

func main() {
	switch *mode_ptr {
	case "r":
		urng := distr.NewuniformGenerator(*lower_bound_ptr, *upper_bound_ptr)
		for i := 0; i < *num_to_generate_ptr; i++ {
			fmt.Println(urng.GenerateNumber())
		}
	case "n":
		grng := distr.NewgaussianGenerator(*lower_bound_ptr, *upper_bound_ptr-1)
		num_list := grng.GenerateNumbers(*num_to_generate_ptr)
		for _, num := range num_list {
			fmt.Println(num)
		}
	default:
		panic("didnt get a valid mode type")
	}
}

var (
	mode_ptr            *string
	num_to_generate_ptr *int
	lower_bound_ptr     *int
	upper_bound_ptr     *int
)

func init() {
	mode_ptr = flag.String("m", "r", "random/uniform(r) or normal(n)")
	num_to_generate_ptr = flag.Int("n", 10000, "number of nums to generate")
	lower_bound_ptr = flag.Int("l", 0, "lower bound on number to generate")
	upper_bound_ptr = flag.Int("u", 101, "upper bound on number to generate")

	flag.Parse()
}
