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
	value          []byte
	currentByte    byte
	destIndex      int
	isQuoted       bool
	isEscaped      bool
	extraSpaceLen  int
	lastCharIsPipe bool
}

func (p *preprocessor) preprocess(value []byte) []byte {
	p.value = value
	p.destIndex = 0
	p.isQuoted = false
	p.isEscaped = false
	p.extraSpaceLen = 0
	p.lastCharIsPipe = false

	for _, currentByte := range p.value {
		p.currentByte = currentByte
		p.processByte()

	}
	return value[:p.destIndex]
}

func (p *preprocessor) processByte() {
	if p.currentByte == escape {
		if p.isEscaped {
			p.value[p.destIndex] = escape
			p.destIndex++
			p.isEscaped = false
		} else {
			p.extraSpaceLen = 0
			p.lastCharIsPipe = false
			p.isEscaped = true
		}
		return
	}

	if p.isEscaped {
		switch p.currentByte {
		case 'r':
			p.value[p.destIndex] = '\r'
		case 'n':
			p.value[p.destIndex] = '\n'
		case dblQuote:
			p.value[p.destIndex] = dblQuote
		default:
			panic(fmt.Sprintf(`invalid escape sequence \%v`, string(p.currentByte)))
		}
		p.destIndex++
		p.isEscaped = false
		return
	}

	if p.currentByte == dblQuote {
		p.isQuoted = !p.isQuoted
		p.extraSpaceLen = 0
		p.lastCharIsPipe = false
		return
	}

	if p.currentByte == space {
		if p.isQuoted {
			p.value[p.destIndex] = p.currentByte
			p.destIndex++
		} else {
			p.extraSpaceLen++
		}
		return
	}

	if p.extraSpaceLen > 0 && p.currentByte != pipe && !p.lastCharIsPipe {
		for index := 0; index < p.extraSpaceLen; index++ {
			p.value[p.destIndex] = space
			p.destIndex++
		}
	}
	p.extraSpaceLen = 0
	if p.currentByte == pipe && !p.isQuoted {
		p.lastCharIsPipe = true
		p.value[p.destIndex] = null
	} else {
		p.lastCharIsPipe = false
		p.value[p.destIndex] = p.currentByte
	}
	p.destIndex++

}
