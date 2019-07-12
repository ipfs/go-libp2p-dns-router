package routinghelpers

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"strings"

	cid "github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/routing"
	record_pb "github.com/libp2p/go-libp2p-record/pb"
)

//FIXME: TODO: specs and implementation 👌
//var _ routing.Routing = (*DNSRouter)(nil)
var _ routing.ValueStore = (*DNSValueStore)(nil)

// DNS uses domains for record storage
type DNSValueStore struct {
}

// Get value from FQDN's TXT
//TODO: opts handling
func (dnsRouter *DNSValueStore) GetValue(ctx context.Context, key string, opts ...routing.Option) ([]byte, error) {
	txtFragments, err := net.LookupTXT(key + ".dns.ipns.dev")
	if err != nil {
		return nil, err
	}

	dec, err := base64.StdEncoding.DecodeString(strings.Join(txtFragments, ""))
	if err != nil {
		return nil, err
	}

	record := record_pb.Record{}
	// defer func() { fmt.Printf("Record: %s\n", record.String()) }() // DBG
	if err = record.Unmarshal(dec); err != nil {
		//FIXME: error is unhandled
		// test data from nameserver does not contain TimeReceived
		fmt.Println("err:", err)
	}

	//FIXME: for some reason key value seems to be reversed?
	//return record.Value, nil
	return record.Key, nil
}

func (dnsRouter *DNSValueStore) SearchValue(ctx context.Context, key string, opts ...routing.Option) (<-chan []byte, error) {
	nameChan := make(chan []byte)
	go func() {
		for range ctx.Done() {
			value, err := dnsRouter.GetValue(ctx, key, opts...)
			if err != nil {
				return
			}

			nameChan <- value
		}
	}()

	return nameChan, nil
}

// Put value to FQDN/nameserver for storage
func (dnsRouter *DNSValueStore) PutValue(ctx context.Context, key string, value []byte, opts ...routing.Option) error {
	builder := cid.V1Builder{}
	refCid, err := builder.Sum([]byte(key))
	if err != nil {
		return err
	}

	jsonStr := fmt.Sprintf(`{
		key: %s,
		record: %s,
		subdomain: true,
	}
	`, refCid.String(),
		base64.StdEncoding.EncodeToString(value))

	req, err := http.NewRequest("PUT", "https://ipns.dev", bytes.NewBufferString(jsonStr))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("DBG: %#v\n", resp.Status)
	/* TODO: what status is expected?
	if resp.Status != http.StatusOK {
	    return someError
	}
	*/

	return nil
}
