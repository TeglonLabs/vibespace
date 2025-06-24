package main
import ("fmt"; "math")
func main() {
  // Test case: center=0.5, neighbors=[0.1, 0.9] avg=0.5, diff=0
  fmt.Println("math.Exp(-0.0):", math.Exp(-0.0))
  // Actual failing case
  fmt.Println("math.Exp(-0.4):", math.Exp(-0.4))
}
