package main

import "bytes"

type commandBuffer struct{ bytes.Buffer }

func (cb *commandBuffer) Close() error {
	cb.Reset()
	return nil
}
