// internal/images/rewriter.go

package images

import "bytes"

func Rewrite(md []byte, mapping map[string]string) []byte {
	out := md
	for from, to := range mapping {
		out = bytes.ReplaceAll(out, []byte(from), []byte(to))
	}
	return out
}
