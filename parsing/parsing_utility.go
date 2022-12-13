package parsing

import (
    "bytes"
)


func ScanSubstr(substr []byte) func (data []byte, atEOF bool) (advance int, token []byte, err error) {
    f := func(data []byte, atEOF bool) (int, []byte, error) {
        if atEOF && len(data) == 0 {
            return 0, nil, nil
        }
        // Look until we find substr
        index := -1
        for i := 0; i < len(data); i++ {
            ss := data[i:i+len(substr)]
            if bytes.Equal(ss, substr) {
                index = i
                return index + len(substr), []byte(data[:i+1]), nil
            }
        }

        // We spit out last output
        if atEOF {
            return len(data), data, nil
        }
        return 0, nil, nil
    }
    return f
}
