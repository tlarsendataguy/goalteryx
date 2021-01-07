package import_file

import "fmt"

const space byte = ' '
const dblQuote byte = '"'
const pipe byte = '|'
const escape byte = '\\'
const null byte = 0

func Preprocess(value []byte) []byte {
	p := &preprocessor{}
	return p.preprocess(value)
}

type preprocessor struct {
	length         int
	destIndex      int
	isQuoted       bool
	isEscaped      bool
	extraSpaceLen  int
	lastCharIsPipe bool
}

func (p *preprocessor) preprocess(value []byte) []byte {
	length := len(value)
	destIndex := 0
	isQuoted := false
	isEscaped := false
	extraSpaceLen := 0
	lastCharIsPipe := false

	for sourceIndex := 0; sourceIndex < length; sourceIndex++ {
		currentByte := value[sourceIndex]

		if currentByte == escape {
			if isEscaped {
				value[destIndex] = escape
				destIndex++
				isEscaped = false
			} else {
				extraSpaceLen = 0
				lastCharIsPipe = false
				isEscaped = true
			}
			continue
		}

		if isEscaped {
			switch currentByte {
			case 'r':
				value[destIndex] = '\r'
			case 'n':
				value[destIndex] = '\n'
			case dblQuote:
				value[destIndex] = dblQuote
			default:
				panic(fmt.Sprintf(`invalid escape sequence \%v at position %v`, string(currentByte), sourceIndex))
			}
			destIndex++
			isEscaped = false
			continue
		}

		if currentByte == dblQuote {
			isQuoted = !isQuoted
			extraSpaceLen = 0
			lastCharIsPipe = false
			continue
		}

		if currentByte == space {
			if isQuoted {
				value[destIndex] = currentByte
				destIndex++
			} else {
				extraSpaceLen++
			}
			continue
		}

		if extraSpaceLen > 0 && currentByte != pipe && !lastCharIsPipe {
			for index := 0; index < extraSpaceLen; index++ {
				value[destIndex] = space
				destIndex++
			}
		}
		extraSpaceLen = 0
		if currentByte == pipe && !isQuoted {
			lastCharIsPipe = true
			value[destIndex] = null
		} else {
			lastCharIsPipe = false
			value[destIndex] = currentByte
		}
		destIndex++
	}
	return value[:destIndex]
}
