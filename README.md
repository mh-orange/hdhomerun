# hdhomerun

[![Build Status](https://travis-ci.org/mh-orange/hdhomerun.svg?branch=master)](https://travis-ci.org/mh-orange/hdhomerun) 
[![GoDoc](https://godoc.org/github.com/mh-orange/hdhomerun?status.png)](https://godoc.org/github.com/mh-orange/hdhomerun) 
[![Coverage Status](https://coveralls.io/repos/github/mh-orange/hdhomerun/badge.svg?branch=master)](https://coveralls.io/github/mh-orange/hdhomerun) 

====== 

## Overview

Use this package to connect to and interact with HD HomeRun devices.

## Installation

```bash
go get github.com/mh-orange/hdhomerun
```

## Examples

###### Discover devices:

```go

package main

import(
  "fmt"
  "github.com/mh-orange/hdhomerun"
)

func main() {
  for discoverResult := range hdhomerun.Discover(nil, time.Millisecond*200) {
    if discoverResult.Err != nil {
      fmt.Printf("Error during discovery: %v\n", discoverResult.Err)
    } else {
      fmt.Printf("hdhomerun device %s found at %s\n", discoverResult.Device.ID(), discoverResult.Device.Addr())
    }
  }
}
```

###### Tune to a channel:

```go
package main

import(
  "github.com/mh-orange/hdhomerun"
  "net"
)

func main() {
  device, _ := hdhomerun.ConnectTCP(&net.TCPAddr{net.IP{192,168,1,100}, 65001, ""})
  tuner := device.Tuner(0)
  tuner.Tune("auto", 177000000)
}
```
###### Scan available channels:

```go
package main

import(
  "fmt"
  "github.com/mh-orange/hdhomerun"
  "net"
)

func main() {
  device, _ := hdhomerun.ConnectTCP(&net.TCPAddr{net.IP{192,168,1,100}, 65001, ""})
  tuner := device.Tuner(0)
  for result := range tuner.Scan() {
    if result.Err != nil {
      fmt.Printf("Error scanning for channels: %v\n", err)
      continue
    }
    fmt.Printf("Found channel %s\n", result.Channel.Name)
  }
}
```

