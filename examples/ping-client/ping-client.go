package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"time"

	ndn "github.com/paulbellamy/go-ndn"
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

var usage = `Usage: ping-client [options] ndn:/name/prefix

Ping a NDN name prefix using Interests with name ndn:/name/prefix/ping/number.
The numbers in the Interests are randomly generated unless specified.
`

// TODO: Just hardcode to where nfd runs for now. We should to more in-depth
// stuff to configure this, like:
// https://github.com/named-data/ndn-cxx/blob/master/src/transport/unix-transport.cpp
var defaultSocketName = "/var/run/nfd.sock"
var defaultAddress = "192.168.59.103:6363"
var minimumPingInterval = 1000 * time.Millisecond
var pingTimeoutThreshold = 4000 * time.Millisecond

var flags = flag.NewFlagSet(usage, flag.PanicOnError)
var pingInterval = flags.Duration("i", minimumPingInterval, fmt.Sprintf("set ping interval in seconds (minimum %v)", minimumPingInterval))
var totalPings = flags.Int("c", -1, "set total number of pings")
var startPingNumber = flags.Int("n", -1, "set the starting number, the number is incremented by 1 after each Interest")
var clientIdentifier = flags.Int("p", 0, "add identifier to the Interest names before the numbers to avoid conflict")
var allowCaching = flags.Bool("a", false, "allow routers to return stale Data from cache")
var printTimestamp = flags.Bool("t", false, "print timestamp")
var quietMode = flags.Bool("q", false, "quiet output")

func parseArgs() (prefix string) {
	defer func() {
		if r := recover(); r != nil {
			printUsage()
		}
	}()
	flags.Parse(os.Args[1:])

	if *pingInterval < minimumPingInterval {
		printUsage()
	}

	if *startPingNumber < 0 {
		*startPingNumber = rand.Int()
	}

	return address(flags)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	prefix := parseArgs()

	var out io.Writer = os.Stdout
	if *quietMode {
		out = ioutil.Discard
	}
	errOut := os.Stderr

	fmt.Fprintln(out, "=== Pinging prefix ===")

	transport, err := net.Dial("tcp", defaultAddress)
	if err != nil {
		fmt.Fprintln(errOut, err)
		os.Exit(1)
	}
	face := ndn.NewFace(transport)
	// TODO: Catch sigterm and sigint to print final summary
	// or can we just defer this?
	// finish(face, out, prefix)

	timer := time.NewTicker(*pingInterval)
	defer timer.Stop()

	var pingsReceived int

	for pingsSent := 0; *totalPings < 0 || pingsSent < *totalPings; pingsSent++ {
		select {
		case <-timer.C:
			go performPing(out, errOut, face, prefix, startPingNumber, &pingsReceived)
		}
	}
}

func pingPacketName(prefix string, nextPingNumber int) (name.Name, error) {
	n, err := name.ParseURI(prefix)
	if err != nil {
		return nil, err
	}

	n = n.AppendString("ping")
	if *clientIdentifier != 0 {
		n = n.AppendString(fmt.Sprint(*clientIdentifier))
	}
	n = n.AppendString(fmt.Sprint(nextPingNumber))
	return n, nil
}

func interestPacket(prefix string, nextPingNumber int) (*packets.Interest, error) {
	n, err := pingPacketName(prefix, nextPingNumber)
	interest := &packets.Interest{}
	interest.SetName(n)
	interest.SetMustBeFresh(!*allowCaching)
	interest.SetInterestLifetime(pingTimeoutThreshold)
	return interest, err
}

func finish(face *ndn.Face, out io.Writer, prefix string, sentPings, receivedPings int) {
	face.Close()
	printPingStatistics(out, prefix)
	if sentPings != receivedPings {
		os.Exit(1)
	}
}

func performPing(out, errOut io.Writer, face *ndn.Face, prefix string, nextPingNumber *int, pingsReceived *int) {
	interest, err := interestPacket(prefix, *nextPingNumber)
	if err != nil {
		fmt.Fprintln(errOut, "Error:", err)
		os.Exit(1)
	}
	(*nextPingNumber)++

	pendingInterest, err := face.ExpressInterest(interest)
	if err != nil {
		fmt.Fprintln(errOut, "Error:", err)
		os.Exit(1)
	}
	start := time.Now()
	select {
	case <-pendingInterest.Data:
		now := time.Now()
		roundTripTime := now.Sub(start)
		(*pingsReceived)++

		outputTimestamp := ""
		if *printTimestamp {
			outputTimestamp = fmt.Sprint(now.Format(time.RFC3339Nano), " - ")
		}
		fmt.Fprintf(
			out,
			"%sContent From %s - Ping Reference = %d  \t- Round Trip Time = %v\n",
			outputTimestamp,
			prefix,
			interest.GetName().GetSuffix(1).ToURI(),
			roundTripTime,
		)
	case <-pendingInterest.Timeout:
		now := time.Now()

		outputTimestamp := ""
		if *printTimestamp {
			outputTimestamp = fmt.Sprint(now.Format(time.RFC3339Nano), " - ")
		}

		fmt.Fprintf(
			out,
			"%sTimeout From %s - Ping Reference = %d\n",
			outputTimestamp,
			prefix,
			interest.GetName().GetSuffix(1).ToURI(),
		)
	}
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
