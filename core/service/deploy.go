package service

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"regexp"

	"github.com/shinzonetwork/shinzo-view-creator/core/models"
	schemastore "github.com/shinzonetwork/shinzo-view-creator/core/schema/store"
	viewstore "github.com/shinzonetwork/shinzo-view-creator/core/view/store"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var SHINZO_HUB_PRECOMPILED_VIEW_REGISTRY_ADDRESS = "0x0000000000000000000000000000000000000210"

type ViewLite struct {
	Query     *string          `json:"query"`
	Sdl       *string          `json:"sdl"`
	Transform models.Transform `json:"transform"`
}

func StartLocalNodeTestAndDeploy(name string, viewstore viewstore.ViewStore, schemastore schemastore.SchemaStore, wallet Wallet, rpc string, debug bool) error {
	fmt.Println("🔧 Building and testing view before deployment...")

	var stdout, stderr *os.File
	var null *os.File

	if !debug {
		null, _ = os.Open(os.DevNull)
		stdout = os.Stdout
		stderr = os.Stderr
		os.Stdout = null
		os.Stderr = null
	}

	err := StartLocalNodeAndTestView(name, viewstore, schemastore, debug)

	if !debug {
		os.Stdout = stdout
		os.Stderr = stderr
		_ = null.Close()
	}

	if err != nil {
		return fmt.Errorf("❌ View failed to build or pass tests: %w", err)
	}

	fmt.Println("⏳ View built and tested successfully. Deploying...")

	// get the view, and parse it to complete string blob
	privateKey, err := DerivePrivateKey(wallet)
	if err != nil {
		return err
	}

	view, err := viewstore.Load(name)

	for i := range view.Transform.Lenses {
		lens := &view.Transform.Lenses[i]
		blob, err := viewstore.GetAssetBlob(name, lens.Label)
		if err != nil {
			return fmt.Errorf("failed to get blob for lens %q: %w", lens.Label, err)
		}
		lens.Path = blob
	}

	viewLite := ViewLite{
		Query:     view.Query,
		Sdl:       view.Sdl,
		Transform: view.Transform,
	}

	data, err := json.Marshal(viewLite)
	if err != nil {
		return fmt.Errorf("failed to marshal view: %w", err)
	}

	viewHash, viewid, err := ComputeViewID(privateKey, data)
	if err != nil {
		return err
	}

	hash, err := sendRegisterTx(rpc, SHINZO_HUB_PRECOMPILED_VIEW_REGISTRY_ADDRESS, privateKey, data)
	if err != nil {
		return err
	}

	fmt.Println("✅ View deployment successful!")
	fmt.Println("----------------------------------------")
	fmt.Printf("🔑 View ID:           %s\n", viewid)
	fmt.Printf("🔑 View Key:          %s\n", viewHash)
	fmt.Printf("📦 Transaction Hash:  %s\n", hash)
	fmt.Println("Blob size (bytes):", len(data))
	fmt.Println("----------------------------------------")
	fmt.Println("You can use the view ID and hash to query or verify registration on-chain.")

	return nil
}

func ComputeViewID(privateKey *ecdsa.PrivateKey, blob []byte) (common.Hash, string, error) {
	sender := crypto.PubkeyToAddress(privateKey.PublicKey)

	re := regexp.MustCompile(`type\s+([A-Za-z0-9_]+)`) // TODO: more robost regex to get sdl type
	matches := re.FindStringSubmatch(string(blob))
	if len(matches) < 1 {
		return common.Hash{}, "", fmt.Errorf("invalid SDL, could not get resource name")
	}

	ResourceName := matches[1]

	// Concatenate sender address and blob (same as abi.encodePacked)
	combined := append(sender.Bytes(), blob...)

	// keccak256 hash
	viewHash := crypto.Keccak256Hash(combined)

	return viewHash, fmt.Sprintf("%s_%s", ResourceName, viewHash.Hex()), nil
}

func sendRegisterTx(
	rpcURL string,
	contractAddr string,
	privateKey *ecdsa.PrivateKey,
	payload []byte,
) (string, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to RPC: %w", err)
	}
	defer client.Close()

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID: %w", err)
	}

	// ---------------------------------------
	// Encode input for register(bytes)
	// methodID = keccak256("register(bytes)")[:4]
	// data = [methodID][offset][length][payload]
	// ---------------------------------------

	methodSig := []byte("register(bytes)")
	methodID := crypto.Keccak256(methodSig)[:4]

	// dynamic offset starts at 0x20
	offset := make([]byte, 32)
	copy(offset[31:], []byte{32})

	// actual payload (view blob)
	payloadLength := big.NewInt(int64(len(payload))).Bytes()
	length := make([]byte, 32)
	copy(length[32-len(payloadLength):], payloadLength)

	// pad the payload to 32-byte alignment
	padding := 32 - (len(payload) % 32)
	if padding == 32 {
		padding = 0
	}
	paddedPayload := append(payload, make([]byte, padding)...)

	// build final calldata
	input := append(methodID, append(offset, append(length, paddedPayload...)...)...)

	// Build and sign transaction
	to := common.HexToAddress(contractAddr)
	tx := types.NewTransaction(nonce, to, big.NewInt(0), 99000000, gasPrice, input)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign tx: %w", err)
	}

	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		return "", fmt.Errorf("failed to send tx: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}
