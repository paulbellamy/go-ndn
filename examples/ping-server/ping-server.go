package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"time"

	ndn "github.com/paulbellamy/go-ndn"
	"github.com/paulbellamy/go-ndn/encoding/tlv"
	"github.com/paulbellamy/go-ndn/name"
	"github.com/paulbellamy/go-ndn/packets"
)

type argSource interface {
	Arg(int) string
	PrintDefaults()
}

func address(flags argSource) string {
	a := flags.Arg(0)
	if a == "" {
		printUsage()
	}
	return a
}

func printUsage() {
	fmt.Println(usage)
	flags.PrintDefaults()
	os.Exit(1)
}

var usage = `Usage: ping-server [options] ndn:/name/prefix

Ping a NDN name prefix using Interests with name ndn:/name/prefix/ping/number.
The numbers in the Interests are randomly generated unless specified.
`

// TODO: Just hardcode to where nfd runs for now. We should to more in-depth
// stuff to configure this, like:
// https://github.com/named-data/ndn-cxx/blob/master/src/transport/unix-transport.cpp
var defaultSocketName = "/var/run/nfd.sock"
var defaultAddress = "192.168.59.103:6363"
var minimumFreshnessPeriod = 1000 * time.Millisecond

var flags = flag.NewFlagSet(usage, flag.PanicOnError)
var freshnessPeriod = flags.Duration("x", minimumFreshnessPeriod, "set FreshnessSeconds")
var maximumPings = flags.Int("p", -1, "specify number of pings to be satisfied (>=1)")
var printTimestamp = flags.Bool("t", false, "print timestamp")
var quietMode = flags.Bool("q", false, "quiet output")

func parseArgs() (prefix string) {
	defer func() {
		if r := recover(); r != nil {
			printUsage()
		}
	}()
	flags.Parse(os.Args[1:])

	if *freshnessPeriod < minimumFreshnessPeriod {
		printUsage()
	}

	return address(flags)
}

func main() {
	prefix := parseArgs()

	var out io.Writer = os.Stdout
	if *quietMode {
		out = ioutil.Discard
	}
	errOut := os.Stderr

	fmt.Fprintln(out, "=== Ping Server ===")

	face, err := newFace()
	if err != nil {
		fmt.Fprintln(errOut, err)
		os.Exit(1)
	}

	n, err := pingPacketName(prefix)
	if err != nil {
		fmt.Fprintln(errOut, err)
		os.Exit(1)
	}

	interests, err := face.RegisterPrefix(n.AppendString("ping"))
	if err != nil {
		fmt.Fprintln(errOut, err)
		os.Exit(1)
	}

	var pingsReceived int
	keyChain := ndn.NewKeyChain()

	for pingsReceived := 0; *maximumPings < 0 || pingsReceived < *maximumPings; pingsReceived++ {
		select {
		case interest, ok := <-interests:
			if !ok {
				break
			}
			logInterestPacket(out, interest)
			respondToPing(out, errOut, face, keyChain, interest)
		}
	}

	if *maximumPings > 0 && *maximumPings <= pingsReceived {
		fmt.Println("\n\nTotal Ping Interests Processed = ", pingsReceived)
		fmt.Printf("Shutting Down Ping Server (%s).\n", prefix)
		face.Close()
	}
}

func newFace() (*ndn.Face, error) {
	transport, err := net.Dial("tcp", defaultAddress)
	if err != nil {
		return nil, err
	}
	return ndn.NewFace(transport), nil
}

func pingPacketName(prefix string) (name.Name, error) {
	n, err := name.ParseURI(prefix)
	if err != nil {
		return nil, err
	}

	return n.AppendString("ping"), nil
}

func logInterestPacket(out io.Writer, interest *packets.Interest) {
	if *printTimestamp {
		fmt.Fprint(out, time.Now().Format(time.RFC3339Nano), " - ")
	}
	fmt.Fprintln(out, "Interest Received - Ping Reference = ", interest.GetName().Get(-1).ToURI())
}

func respondToPing(out, errOut io.Writer, face *ndn.Face, keyChain *ndn.KeyChain, interest *packets.Interest) {
	face.Put(dataPacket(keyChain, interest))
}

func dataPacket(keyChain *ndn.KeyChain, interest *packets.Interest) *packets.Data {
	data := &packets.Data{}
	data.SetName(interest.GetName())
	data.SetFreshnessPeriod(*freshnessPeriod)
	data.SetContent([]byte("NDN TLV Ping Response"))
	keyChain.Sign(data, name.New(name.Component{"What do I put for certificate name here?"}), tlv.NewEncoder)
	return data
}

func printPingStatistics(out io.Writer, prefix string) {
	/* TODO: Stats output
	fmt.Fprintf(
		out,
		"\n\n===  Ping Statistics For %s ===\nSent=%v, Received=%v, Packet Loss=%v,Total Time=%v\nRound Trip Time (Min/Max/Avg/MDev) = (%v/%v/%v/%v)\n",
		prefix,
		sent,
		received,
		packetsLost,
		totalTime,
		minRoundTrip,
		maxRoundTrip,
		avgRoundTrip,
		mdevRoundTrip,
	)
	*/
}
