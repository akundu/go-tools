package main

import (
	"flag"
	"fmt"
	"math"
	"sort"
	"time"

	rng "github.com/leesper/go_rng"
)

func mappingValues(num, old_min, old_max, new_min, new_max int) int {
	old_range := old_max - old_min
	if old_range <= 0 {
		return 0
	}

	new_range := new_max - new_min
	if new_range <= 0 {
		return 0
	}

	return (((new_max - new_min) * (num - old_min)) / (old_max - old_min)) + new_min
}

func closestPowerOf10(num int) int {
	c_power_of_10 := float64(int(math.Log10(float64(num))))
	if math.Pow(10, c_power_of_10) < float64(num) {
		return int(math.Pow(10, c_power_of_10+1))
	}
	return int(math.Pow(10, c_power_of_10))
}

func main() {
	switch *mode_ptr {
	case "r":
		urng := rng.NewUniformGenerator(time.Now().UnixNano())
		for i := 0; i < *num_to_generate_ptr; i++ {
			fmt.Println(urng.Int32Range(int32(*lower_bound_ptr), int32(*upper_bound_ptr)))
		}
	case "n":
		rounded_power_of_10 := closestPowerOf10(*upper_bound_ptr - 1)

		grng := rng.NewGaussianGenerator(time.Now().UnixNano())
		num_list := make([]int, *num_to_generate_ptr)

		for i := 0; i < *num_to_generate_ptr; i++ {
			num_list[i] = int(grng.Gaussian(0, 3) * float64(rounded_power_of_10))
		}

		sort.Ints(num_list)
		for _, num := range num_list {
			fmt.Println(mappingValues(num, num_list[0], num_list[len(num_list)-1], *lower_bound_ptr, *upper_bound_ptr-1))
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
