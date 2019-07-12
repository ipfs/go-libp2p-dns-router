package routinghelpers

import (
	"bytes"
	"context"
	"fmt"
	"testing"
)

func TestDnsValueStore(t *testing.T) {
	d := DNSValueStore{}

	ctx := context.Background()

	type pair struct {
		k string
		v string
	}

	for _, getPair := range []pair{
		{"bafybeigv6xgwkfhx3abfsuayuicb3l7xblpzrsdvjlt33oap43joacglyu", "/ipfs/QmZFLGKTiYvxhCAQGDDSCvRPYeQCTaQMBAPL1Cqb8S695p"},
	} {

		v, err := d.GetValue(ctx, getPair.k)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("Key:%q\nValue:%s\nExpected value:%s\n", getPair.k, v, getPair.v)
		if !bytes.Equal([]byte(getPair.v), v) {
			t.Fatal(fmt.Sprintf("Key %q has expected value\n%sbut received\n%s\n", getPair.k, getPair.v, v))
		}
	}

	/* TODO: put tests
	for _, putPair := range[]pair{
	    {"", ""},
	} {
	}
	*/

}
