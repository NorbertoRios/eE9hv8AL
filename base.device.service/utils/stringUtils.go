package utils

//SplitString splits string using delimiters
func SplitString(input string, delims []string) []string {
	splits := []string{}

	prev := 0
	for pos, r1 := range input {
		for _, d := range delims {
			if d != "" && d == string(r1) {
				if input[prev:pos] != "" {
					splits = append(splits, input[prev:pos])
				}
				prev = pos + 1
			}
		}
	}
	if input[prev:] != "" {
		splits = append(splits, input[prev:])
	}
	return splits
}
