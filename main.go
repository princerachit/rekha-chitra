package main

func main() {
	m := map[string]int32{"a": 1, "b": 2, "c": 3}
	m2 := map[string]int64{"a": 1, "b": 2, "c": 3}
	println(sumOfLongOrInts(m))
	println(sumOfLongOrInts(m2))
	print("hello world")
}

func sumOfLongOrInts[K comparable, V int32 | int64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}
