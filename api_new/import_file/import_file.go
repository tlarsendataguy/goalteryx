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
	currentByte    byte
	destIndex      int
	isQuoted       bool
	isEscaped      bool
	extraSpaceLen  int
	lastCharIsPipe bool
}

func (p *preprocessor) preprocess(value []byte) []byte {
	p.destIndex = 0
	p.isQuoted = false
	p.isEscaped = false
	p.extraSpaceLen = 0
	p.lastCharIsPipe = false

	for sourceIndex, currentByte := range value {
		p.currentByte = currentByte

		if p.currentByte == escape {
			if p.isEscaped {
				value[p.destIndex] = escape
				p.destIndex++
				p.isEscaped = false
			} else {
				p.extraSpaceLen = 0
				p.lastCharIsPipe = false
				p.isEscaped = true
			}
			continue
		}

		if p.isEscaped {
			switch p.currentByte {
			case 'r':
				value[p.destIndex] = '\r'
			case 'n':
				value[p.destIndex] = '\n'
			case dblQuote:
				value[p.destIndex] = dblQuote
			default:
				panic(fmt.Sprintf(`invalid escape sequence \%v at position %v`, string(p.currentByte), sourceIndex))
			}
			p.destIndex++
			p.isEscaped = false
			continue
		}

		if p.currentByte == dblQuote {
			p.isQuoted = !p.isQuoted
			p.extraSpaceLen = 0
			p.lastCharIsPipe = false
			continue
		}

		if p.currentByte == space {
			if p.isQuoted {
				value[p.destIndex] = p.currentByte
				p.destIndex++
			} else {
				p.extraSpaceLen++
			}
			continue
		}

		if p.extraSpaceLen > 0 && p.currentByte != pipe && !p.lastCharIsPipe {
			for index := 0; index < p.extraSpaceLen; index++ {
				value[p.destIndex] = space
				p.destIndex++
			}
		}
		p.extraSpaceLen = 0
		if p.currentByte == pipe && !p.isQuoted {
			p.lastCharIsPipe = true
			value[p.destIndex] = null
		} else {
			p.lastCharIsPipe = false
			value[p.destIndex] = p.currentByte
		}
		p.destIndex++
	}
	return value[:p.destIndex]
}
