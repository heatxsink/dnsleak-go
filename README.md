# dnsleak-go

Found this program at https://github.com/macvk/dnsleaktest/blob/master/dnsleaktest.go figured it could use some love. 

## Description 

The test shows DNS leaks and your external IP. If you use the same ASN for DNS and connection - you have no leak, otherwise here might be a problem.


## Prerequisites

- golang 1.15 or higher.


## Installation

    $ go get github.com/heatxsink/dnsleak-go/...
    $ dnsleak

