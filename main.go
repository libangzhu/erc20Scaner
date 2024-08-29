package main

import (
	"flag"
)

var (
	rawUrl     = flag.String("u", "https://mainnet.bityuan.com/eth", "node url")
	startPoint = flag.Int64("s", 34562327, "start point")
)

func main() {
	flag.Parse()
	p := new(Process)
	p.startPoint = uint64(*startPoint)
	p.Init()
	p.Start()

}
