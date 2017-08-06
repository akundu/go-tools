package main

import (
	"flag"
	"fmt"
	"time"

	rng "github.com/leesper/go_rng"
)

func main() {
	switch *mode_ptr {
	case "r":
		urng := rng.NewUniformGenerator(time.Now().UnixNano())
		for i := 0; i < *num_to_generate_ptr; i++ {
			fmt.Println(urng.Int32Range(int32(*lower_bound_ptr), int32(*upper_bound_ptr)))
		}
	case "n":
		grng := rng.NewGaussianGenerator(time.Now().UnixNano())
		for i := 0; i < *num_to_generate_ptr; i++ {
			fmt.Println(grng.StdGaussian())
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
