package json_comment

import (
	"encoding/json"
	"io"
	"strings"
	"testing"
)

func TestStrippedReaderRead(t *testing.T) {
	makeReader := func(s string) io.Reader {
		return NewStrippedReader(strings.NewReader(s))
	}

	shouldNotFail := func(input string) {
		var m struct{}
		dec := json.NewDecoder(makeReader(input))
		if err := dec.Decode(&m); err != nil {
			t.Errorf("should not fail: %s\n %v\n", input, err)
		}
	}

	normalTestcases := []struct{ input string }{
		{`{}`},
		{`{
                        // foobar
                        # lol
                }`},
		{`{
                        "a_comment_with#_hash": "// this should pass",
                        "/* this is a key */": [
                                "this //"
                                /* is
                                 * a value
                                 */
                        ]
                }`},
		{`{
                        /* multiple
                           line comment works too */ "key": "value"
                 }`},
	}

	for _, testcase := range normalTestcases {
		shouldNotFail(testcase.input)
	}
}
