// Strip comments from JSON input.
//
// Supported comment types:
//
// - Single line:
//
//   {
//     // this is a single line comment,
//     # this one is a single line comment, too
//   }
//
// - Multiple lines:
//
//   {
//     /* multiple
//      * lines comments
//      * works
//      */
//   }
package json_comment

import (
	"io"
)

type STATE uint64

const (
	S_OTHER = iota
	S_STRING
	S_SL_COMMENT // single line comment
	S_ML_COMMENT // multiple lines comment

	SYM_LINE_BREAK  = "\n"
	SYM_QUOTE       = "\""
	SYM_HASH        = "#"
	SYM_SLASH_SLASH = "//"
	SYM_SLASH_STAR  = "/*"
	SYM_STAR_SLASH  = "*/"
)

func IsQuoteMark(char byte) bool {
	return string(char) == SYM_QUOTE
}

func IsLineBreak(char byte) bool {
	return string(char) == SYM_LINE_BREAK
}

func IsSingleLineComment(char, peakChar byte) bool {
	sc := string(char)
	return ((sc == SYM_HASH) || (sc+string(peakChar) == SYM_SLASH_SLASH))
}

func IsMultiLinesCommentStarts(char, peakChar byte) bool {
	return string(char)+string(peakChar) == SYM_SLASH_STAR
}

func IsMultiLinesCommentEnds(char, peakChar byte) bool {
	return string(char)+string(peakChar) == SYM_STAR_SLASH
}

func readSome(from io.Reader, n int) ([]byte, int, error) {
	buf := make([]byte, n)
	actualRead, err := from.Read(buf)
	return buf, actualRead, err
}

type StrippedReader struct {
	input io.Reader
	state STATE
}

func (s *StrippedReader) Read(p []byte) (int, error) {
	var (
		needLength = len(p)
		readLength = 0

		i              int
		char, peekChar byte

		buf []byte
		err error
	)

	for readLength < needLength {
		buf, _, err = readSome(s.input, needLength-readLength)
		if err != nil {
			return readLength, err
		}

		for i = 0; i < len(buf)-1; i = i + 1 {
			char = buf[i]
			peekChar = buf[i+1]

			switch s.state {
			case S_OTHER:
				if IsSingleLineComment(char, peekChar) {
					s.state = S_SL_COMMENT
					continue // starts to skip comment
				}
				if IsMultiLinesCommentStarts(char, peekChar) {
					s.state = S_ML_COMMENT
					continue // starts to skip comment
				}

				if IsQuoteMark(char) {
					s.state = S_STRING
				}

				p[readLength] = char
				readLength = readLength + 1
			case S_STRING:
				if IsQuoteMark(char) {
					s.state = S_OTHER
				}

				p[readLength] = char
				readLength = readLength + 1
			case S_SL_COMMENT:
				if IsLineBreak(char) {
					s.state = S_OTHER

					// keep line break
					p[readLength] = char
					readLength = readLength + 1
				}
			case S_ML_COMMENT:
				if IsMultiLinesCommentEnds(char, peekChar) {
					s.state = S_OTHER

					i = i + 1 // skip peeked character
				}
			}
		}

		// read last character
		char = peekChar
		if s.state == S_OTHER || s.state == S_STRING {
			p[readLength] = char
			readLength = readLength + 1
		}
	}
	return readLength, nil
}

func NewStrippedReader(input io.Reader) io.Reader {
	return &StrippedReader{
		input: input,
		state: S_OTHER,
	}
}
