package collab

import "strings"

func wordwrap(s string, limit int) string {
	if len(s) < limit {
		return s
	}

	if strings.TrimSpace(s) == "" {
		return s
	}

	// convert string to slice
	strSlice := strings.Fields(s)

	var result string = ""

	for len(strSlice) >= 1 {
		// convert slice/array back to string
		// but insert \r\n at specified limit

		result = result + strings.Join(strSlice[:limit], " ") + "\r\n"

		// discard the elements that were copied over to result
		strSlice = strSlice[limit:]

		// change the limit
		// to cater for the last few words in
		//
		if len(strSlice) < limit {
			limit = len(strSlice)
		}
	}
	return result
}
