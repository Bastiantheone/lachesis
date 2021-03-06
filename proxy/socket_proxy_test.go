package proxy

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/andrecronje/lachesis/common"
	crypto "github.com/andrecronje/lachesis/crypto"
	"github.com/andrecronje/lachesis/poset"
	aproxy "github.com/andrecronje/lachesis/proxy/app"
)

func TestSocketProxyServer(t *testing.T) {
	clientAddr := "127.0.0.1:9990"
	proxyAddr := "127.0.0.1:9991"
	proxy := aproxy.NewSocketAppProxy(clientAddr, proxyAddr, 1*time.Second, common.NewTestLogger(t))
	submitCh := proxy.SubmitCh()

	tx := []byte("the test transaction")

	// Listen for a request
	go func() {
		select {
		case st := <-submitCh:
			// Verify the command
			if !reflect.DeepEqual(st, tx) {
				t.Fatalf("tx mismatch: %#v %#v", tx, st)
			}
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("timeout")
		}
	}()

	// now client part connecting to RPC service
	// and calling methods
	dummyClient, err := NewDummySocketClient(clientAddr, proxyAddr, common.NewTestLogger(t))
	if err != nil {
		t.Fatal(err)
	}
	err = dummyClient.SubmitTx(tx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSocketProxyClient(t *testing.T) {
	clientAddr := "127.0.0.1:9992"
	proxyAddr := "127.0.0.1:9993"

	//launch dummy application
	dummyClient, err := NewDummySocketClient(clientAddr, proxyAddr, common.NewTestLogger(t))
	if err != nil {
		t.Fatal(err)
	}
<<<<<<< HEAD
	clientCh := dummyClient.lachesisProxy.CommitCh()
=======
	initialStateHash := dummyClient.state.stateHash
>>>>>>> 1ed81963270c25f7cc7eaa8bfdef464338bbd00b

	//create client proxy
	proxy := aproxy.NewSocketAppProxy(clientAddr, proxyAddr, 1*time.Second, common.NewTestLogger(t))

	//create a few blocks
	blocks := [5]poset.Block{}
	for i := 0; i < 5; i++ {
		blocks[i] = poset.NewBlock(i, i+1, []byte{}, [][]byte{[]byte(fmt.Sprintf("block %d transaction", i))})
	}

	//commit first block and check that the client's statehash is correct
	stateHash, err := proxy.CommitBlock(blocks[0])
	if err != nil {
		t.Fatal(err)
	}

	expectedStateHash := initialStateHash
	for _, t := range blocks[0].Transactions() {
		tHash := crypto.SHA256(t)
		expectedStateHash = crypto.SimpleHashFromTwoHashes(expectedStateHash, tHash)
	}

	if !reflect.DeepEqual(stateHash, expectedStateHash) {
		t.Fatalf("StateHash should be %v, not %v", expectedStateHash, stateHash)
	}

	snapshot, err := proxy.GetSnapshot(blocks[0].Index())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(snapshot, expectedStateHash) {
		t.Fatalf("Snapshot should be %v, not %v", expectedStateHash, snapshot)
	}

	//commit a few more blocks, then attempt to restore back to block 0 state
	for i := 1; i < 5; i++ {
		_, err := proxy.CommitBlock(blocks[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	err = proxy.Restore(snapshot)
	if err != nil {
		t.Fatalf("Error restoring snapshot: %v", err)
	}

}
