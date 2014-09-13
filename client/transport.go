package client

import "io"

type Transport interface {
	io.Writer
}
