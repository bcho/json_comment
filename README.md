# JSON Comment

[![Build Status](https://travis-ci.org/bcho/json_comment.svg?branch=master)](https://travis-ci.org/bcho/json_comment)
[![GoDoc](https://godoc.org/github.com/bcho/json_comment?status.svg)](https://godoc.org/github.com/bcho/json_comment)

Strip comments from JSON input.


## Comment Types

### Single Line

```json
{
  # this is a comment
  // this one is a comment, too
  "not_a_comment": "// this is not a comment",
  "neither_this_one": "# you know, for comment"
}
```


### Multiple Lines

```json
[
  /* comment can go
   across lines
   "inside_comment",
   */
   "but not contain this one" /* this is another comment */
]
```


## License

MIT
