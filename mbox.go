// Package mbox parses and formats the mbox file format.
//
// As the mbox file format is not standardized this package expects the least
// common denominator, the so called mboxo format.
package mbox

var (
	header        = []byte("From ")
	escapedHeader = append([]byte{'>'}, header...)
)
