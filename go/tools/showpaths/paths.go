// Copyright 2018 ETH Zurich, Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Simple show paths application for SCION.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/env"
	"github.com/scionproto/scion/go/lib/log"
	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/sciond/pathprobe"
	"github.com/scionproto/scion/go/lib/snet"
)

var (
	dstIAStr     = flag.String("dstIA", "", "Destination IA address: ISD-AS")
	srcIAStr     = flag.String("srcIA", "", "Source IA address: ISD-AS")
	sciondPath   = flag.String("sciond", "", "SCIOND socket path")
	timeout      = flag.Duration("timeout", 5*time.Second, "Timeout in seconds")
	maxPaths     = flag.Int("maxpaths", 10, "Maximum number of paths")
	sciondFromIA = flag.Bool("sciondFromIA", false, "SCIOND socket path from IA address:ISD-AS")
	expiration   = flag.Bool("expiration", false, "Show path expiration timestamps")
	refresh      = flag.Bool("refresh", false, "Set refresh flag for SCIOND path request")
	status       = flag.Bool("p", false, "Probe the paths and print out the statuses")
	version      = flag.Bool("version", false, "Output version information and exit.")
)

var (
	dstIA addr.IA
	srcIA addr.IA
	local snet.Addr
)

func init() {
	flag.Var((*snet.Addr)(&local), "local", "Local address to use for health checks")
	flag.Usage = flagUsage
}

func main() {
	log.AddLogConsFlags()
	validateFlags()
	if err := log.SetupFromFlags(""); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s", err)
		flag.Usage()
		os.Exit(1)
	}
	defer log.LogPanicAndExit()

	sd := sciond.NewService(*sciondPath, false)
	var err error
	sdConn, err := sd.ConnectTimeout(*timeout)
	if err != nil {
		LogFatal("Failed to connect to SCIOND", "err", err)
	}
	reply, err := sdConn.Paths(context.Background(), dstIA, srcIA, uint16(*maxPaths),
		sciond.PathReqFlags{Refresh: *refresh})
	if err != nil {
		LogFatal("Failed to retrieve paths from SCIOND", "err", err)
	}
	if reply.ErrorCode != sciond.ErrorOk {
		LogFatal("SCIOND unable to retrieve paths", "ErrorCode", reply.ErrorCode)
	}

	fmt.Println("Available paths to", dstIA)
	var pathStatuses map[string]pathprobe.Status
	if *status {
		ctx, cancelF := context.WithTimeout(context.Background(), *timeout)
		pathStatuses, err = pathprobe.Prober{
			Local: local,
			DstIA: dstIA,
		}.GetStatuses(ctx, reply.Entries)
		cancelF()
		if err != nil {
			LogFatal("Failed to get status", "err", err)
		}
	}
	for i, path := range reply.Entries {
		fmt.Printf("[%2d] %s", i, path.Path.String())
		if *expiration {
			fmt.Printf(" Expires: %s (%s)", path.Path.Expiry(),
				time.Until(path.Path.Expiry()).Truncate(time.Second))
		}
		if *status {
			fmt.Printf(" Status: %s", pathStatuses[pathprobe.PathKey(path)])
		}
		fmt.Printf("\n")
	}
}

func validateFlags() {
	flag.Parse()
	var err error
	if *version {
		fmt.Print(env.VersionInfo())
		os.Exit(0)
	}
	if *dstIAStr == "" {
		LogFatal("Missing destination IA")
	} else {
		dstIA, err = addr.IAFromString(*dstIAStr)
		if err != nil {
			LogFatal("Unable to parse destination IA", "err", err)
		}
	}

	if *srcIAStr != "" {
		if srcIA, err = addr.IAFromString(*srcIAStr); err != nil {
			LogFatal("Unable to parse source IA", "err", err)
		}
	}

	if *sciondFromIA {
		if *sciondPath != "" {
			LogFatal("Only one of -sciond or -sciondFromIA can be specified")
		}
		if srcIA.IsZero() {
			LogFatal("-srcIA flag is missing")
		}
		*sciondPath = sciond.GetDefaultSCIONDPath(&srcIA)
	} else if *sciondPath == "" {
		*sciondPath = sciond.GetDefaultSCIONDPath(nil)
	}

	if *status && (local.IA.IsZero() || local.Host == nil) {
		LogFatal("Local address is required for health checks")
	}
}

func flagUsage() {
	fmt.Fprintf(os.Stderr, `
Usage: showpaths [flags]

Lists available paths between SCION ASes. Paths might be retrieved from a local cache, and they
might not forward traffic successfully (for example, if a network link went down). To probe if the
paths are healthy, use -p.

flags:
`)
	flag.PrintDefaults()
}

func LogFatal(msg string, a ...interface{}) {
	log.Crit(msg, a...)
	os.Exit(1)
}
