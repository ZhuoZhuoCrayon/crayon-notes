package main

func removeComments(source []string) []string {
	res := []string{}
	newline := []byte{}
	inBlock := false
	for _, line := range source {
		for i := 0; i < len(line); i++ {
			if inBlock {
				if i+1 < len(line) && line[i] == '*' && line[i+1] == '/' {
					inBlock = false
					i++
				}
			} else {
				if i+1 < len(line) && line[i] == '/' && line[i+1] == '*' {
					inBlock = true
					i++
				} else if i+1 < len(line) && line[i] == '/' && line[i+1] == '/' {
					break
				} else {
					newline = append(newline, line[i])
				}
			}
		}
		if !inBlock && len(newline) > 0 {
			res = append(res, string(newline))
			newline = []byte{}
		}
	}
	return res
}
