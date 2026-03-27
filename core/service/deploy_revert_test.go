package service

import (
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TestSendRegisterTx_RevertsOnInvalidPayload verifies that sendRegisterTx
// returns an error when the on-chain transaction reverts.
//
// This is an integration test that requires a running local ShinzoHub node.
// It sends garbage bytes (not valid VWL format) to the ViewRegistry
// precompile, which causes viewbundle.DecodeHeader() to fail and the
// precompile to revert. The test confirms that sendRegisterTx detects
// the revert via receipt status instead of silently reporting success.
//
// Prerequisites:
//   - Local ShinzoHub testnet: cd shinzohub && make sh-testnet
//   - Export key: export TEST_PRIVATE_KEY=$(./build/shinzohubd keys unsafe-export-eth-key acc0 --keyring-backend test --home ~/.shinzohub)
//
// Run:
//
//	TEST_PRIVATE_KEY=<hex> go test ./core/service/ -run TestSendRegisterTx_RevertsOnInvalidPayload -v -count=1
func TestSendRegisterTx_RevertsOnInvalidPayload(t *testing.T) {
	rpcURL := "http://localhost:8545"

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		t.Skip("skipping: local chain not reachable at", rpcURL)
	}
	client.Close()

	acc0Key := os.Getenv("TEST_PRIVATE_KEY")
	if acc0Key == "" {
		t.Skip("skipping: set TEST_PRIVATE_KEY env var")
	}

	privateKey, err := crypto.HexToECDSA(acc0Key)
	if err != nil {
		t.Fatalf("invalid private key: %v", err)
	}

	_, err = sendRegisterTx(
		rpcURL,
		SHINZO_HUB_PRECOMPILED_VIEW_REGISTRY_ADDRESS,
		privateKey,
		[]byte("INVALID_VWL_DATA"),
	)

	if err == nil {
		t.Fatal("expected sendRegisterTx to return error for invalid payload, got nil")
	}

	if !strings.Contains(err.Error(), "reverted") {
		t.Errorf("expected error to mention 'reverted', got: %s", err.Error())
	}

	t.Logf("correctly returned error: %v", err)
}
