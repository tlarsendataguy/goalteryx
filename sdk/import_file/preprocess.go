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
	sourceIndex    int
	value          []byte
	destIndex      int
	isQuoted       bool
	isEscaped      bool
	extraSpaceLen  int
	lastCharIsPipe bool
}

func (p *preprocessor) preprocess(value []byte) []byte {
	p.resetFor(value)
	for p.processByte() {
	}
	return value[:p.destIndex]
}

func (p *preprocessor) resetFor(value []byte) {
	p.sourceIndex = 0
	p.value = value
	p.destIndex = 0
	p.isQuoted = false
	p.isEscaped = false
	p.extraSpaceLen = 0
	p.lastCharIsPipe = false
}

func (p *preprocessor) processByte() bool {
	if p.sourceIndex >= len(p.value) {
		return false
	}

	currentByte := p.value[p.sourceIndex]
	p.sourceIndex++

	if currentByte == pipe {
		return p.processPipe()
	}

	if currentByte == escape {
		return p.processEscape()
	}

	if p.isEscaped {
		return p.processCharAfterEscape(currentByte)
	}

	if currentByte == dblQuote {
		return p.processDoubleQuote()
	}

	if currentByte == space {
		return p.processSpace()
	}

	if p.whiteSpaceIsSignificant(currentByte) {
		p.catchUpSignificantWhitespace()
	}

	return p.processNormalChar(currentByte)
}

func (p *preprocessor) processPipe() bool {
	if p.isQuoted {
		return p.processNormalChar('|')
	}
	p.lastCharIsPipe = true
	p.value[p.destIndex] = null
	p.destIndex++
	return true
}

func (p *preprocessor) processEscape() bool {
	if p.isEscaped {
		p.value[p.destIndex] = escape
		p.destIndex++
		p.isEscaped = false
	} else {
		p.extraSpaceLen = 0
		p.lastCharIsPipe = false
		p.isEscaped = true
	}
	return true
}

func (p *preprocessor) processCharAfterEscape(currentByte byte) bool {
	switch currentByte {
	case 'r':
		p.value[p.destIndex] = '\r'
	case 'n':
		p.value[p.destIndex] = '\n'
	case dblQuote:
		p.value[p.destIndex] = dblQuote
	default:
		panic(fmt.Sprintf(`invalid escape sequence \%v at position %v`, string(currentByte), p.sourceIndex))
	}
	p.destIndex++
	p.isEscaped = false
	return true
}

func (p *preprocessor) processDoubleQuote() bool {
	p.isQuoted = !p.isQuoted
	p.extraSpaceLen = 0
	p.lastCharIsPipe = false
	p.value[p.destIndex] = '"'
	p.destIndex++
	return true
}

func (p *preprocessor) processSpace() bool {
	if p.isQuoted {
		p.value[p.destIndex] = space
		p.destIndex++
	} else {
		p.extraSpaceLen++
	}
	return true
}

func (p *preprocessor) processNormalChar(currentByte byte) bool {
	p.extraSpaceLen = 0
	p.lastCharIsPipe = false
	p.value[p.destIndex] = currentByte
	p.destIndex++
	return true
}

func (p *preprocessor) whiteSpaceIsSignificant(currentByte byte) bool {
	return p.extraSpaceLen > 0 && currentByte != pipe && !p.lastCharIsPipe
}

func (p *preprocessor) catchUpSignificantWhitespace() {
	for index := 0; index < p.extraSpaceLen; index++ {
		p.value[p.destIndex] = space
		p.destIndex++
	}
}
