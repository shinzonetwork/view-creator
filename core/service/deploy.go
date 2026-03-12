package service

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"time"

	"github.com/shinzonetwork/shinzo-view-creator/core/models"
	schemastore "github.com/shinzonetwork/shinzo-view-creator/core/schema/store"
	viewstore "github.com/shinzonetwork/shinzo-view-creator/core/view/store"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/shinzonetwork/viewbundle-go"
)

var SHINZO_HUB_PRECOMPILED_VIEW_REGISTRY_ADDRESS = "0x0000000000000000000000000000000000000210"

// ViewLite stays for local testing / debug printing if you want,
// but we do NOT send this JSON on-chain anymore.
type ViewLite struct {
	Query     *string          `json:"query"`
	Sdl       *string          `json:"sdl"`
	Transform models.Transform `json:"transform"`
}

// ---------------------------
// Deploy entry
// ---------------------------

func StartLocalNodeTestAndDeploy(
	name string,
	viewstore viewstore.ViewStore,
	schemastore schemastore.SchemaStore,
	wallet Wallet,
	rpc string,
	debug bool,
) error {
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

	privateKey, err := DerivePrivateKey(wallet)
	if err != nil {
		return err
	}

	// Load view (small metadata, lens references)
	view, err := viewstore.Load(name)
	if err != nil {
		return err
	}
	if view.Query == nil || view.Sdl == nil {
		return fmt.Errorf("view missing query or sdl")
	}

	// Build viewbundle.View inline (no helper func)
	vb := viewbundle.View{
		Query: *view.Query,
		Sdl:   *view.Sdl,
		Transform: viewbundle.Transform{
			Lenses: make([]viewbundle.Lens, 0, len(view.Transform.Lenses)),
		},
	}

	for i := range view.Transform.Lenses {
		l := view.Transform.Lenses[i]

		if l.Path == "" {
			return fmt.Errorf("lens[%d] missing Path (expected base64 wasm)", i)
		}

		// Read the wasm file and base64-encode it for the bundler
		wasmBase64, err := viewstore.GetAssetBlob(name, l.Label)
		if err != nil {
			return fmt.Errorf("lens[%d] failed to load wasm blob: %w", i, err)
		}

		// If it's a map/struct, marshal to JSON string.
		args := ""
		switch v := any(l.Arguments).(type) {
		case string:
			args = v
		default:
			bz, err := json.Marshal(l.Arguments)
			if err != nil {
				return fmt.Errorf("marshal lens[%d] arguments: %w", i, err)
			}
			args = string(bz)
		}

		vb.Transform.Lenses = append(vb.Transform.Lenses, viewbundle.Lens{
			Path:      wasmBase64,
			Arguments: args,
		})
	}

	bd := viewbundle.NewBundler()
	wireBytes, err := bd.BundleView(vb)
	if err != nil {
		return fmt.Errorf("viewbundle.BundleView: %w", err)
	}

	// Gather lens blobs (raw bytes) in the same order as view.Transform.Lenses
	lensBlobs := make([][]byte, 0, len(view.Transform.Lenses))
	for i := range view.Transform.Lenses {
		l := &view.Transform.Lenses[i]

		blob, err := viewstore.GetAssetBlob(name, l.Label)
		if err != nil {
			return fmt.Errorf("failed to get blob for lens %q: %w", l.Label, err)
		}

		// Your GetAssetBlob might return []byte or string depending on your store.
		// Handle both safely.
		switch v := any(blob).(type) {
		case []byte:
			lensBlobs = append(lensBlobs, v)
		case string:
			lensBlobs = append(lensBlobs, []byte(v))
		default:
			// last resort: JSON it
			bz, _ := json.Marshal(v)
			lensBlobs = append(lensBlobs, bz)
		}
	}

	// Compute view key + id based on (sender, wireBytes)
	viewHash, viewID, err := ComputeViewIDFromWire(privateKey, *view.Sdl, wireBytes)
	if err != nil {
		return err
	}

	txHash, err := sendRegisterTx(rpc, SHINZO_HUB_PRECOMPILED_VIEW_REGISTRY_ADDRESS, privateKey, wireBytes)
	if err != nil {
		return err
	}

	fmt.Println("✅ View deployment successful!")
	fmt.Println("----------------------------------------")
	fmt.Printf("🔑 View ID:           %s\n", viewID)
	fmt.Printf("🔑 View Key:          %s\n", viewHash.Hex())
	fmt.Printf("📦 Transaction Hash:  %s\n", txHash)
	fmt.Printf("Wire bytes (tx data):%d\n", len(wireBytes))
	fmt.Println("----------------------------------------")

	return nil
}

// ---------------------------
// View ID (matches precompile key logic)
// ---------------------------

func ComputeViewIDFromWire(privateKey *ecdsa.PrivateKey, sdl string, wireBytes []byte) (common.Hash, string, error) {
	sender := crypto.PubkeyToAddress(privateKey.PublicKey)

	re := regexp.MustCompile(`\btype\s+([A-Za-z0-9_]+)\b`)
	matches := re.FindStringSubmatch(sdl)
	if len(matches) < 2 {
		return common.Hash{}, "", fmt.Errorf("invalid SDL, could not get resource name")
	}
	resourceName := matches[1]

	// Must match chain: keccak(sender, encodedValue)
	viewHash := crypto.Keccak256Hash(sender.Bytes(), wireBytes)

	// Keep your original format (includes 0x). If you want cleaner type names, use viewHash.Hex()[2:].
	viewID := fmt.Sprintf("%s_%s", resourceName, viewHash.Hex())

	return viewHash, viewID, nil
}

// ---------------------------
// Tx send (register(bytes))
// ---------------------------

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

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Encode register(bytes) calldata
	input := EncodeRegisterBytesCalldata(payload)

	to := common.HexToAddress(contractAddr)

	// Estimate gas (better than hardcoding)
	msg := ethereum.CallMsg{
		From: fromAddress,
		To:   &to,
		Data: input,
	}
	gasLimit, err := client.EstimateGas(ctx, msg)
	if err != nil {
		// fallback (keep your old style if estimate fails)
		gasLimit = 100000000
	} else {
		// add padding
		gasLimit = uint64(float64(gasLimit) * 1.20)
	}

	tx := types.NewTransaction(nonce, to, big.NewInt(0), gasLimit, gasPrice, input)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign tx: %w", err)
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return "", fmt.Errorf("failed to send tx: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}

func EncodeRegisterBytesCalldata(payload []byte) []byte {
	// methodID = keccak256("register(bytes)")[:4]
	methodID := crypto.Keccak256([]byte("register(bytes)"))[:4]

	// ABI encoding for single dynamic bytes:
	// 0x00..0x20 offset
	// length
	// data (padded)
	offset := make([]byte, 32)
	offset[31] = 32

	length := make([]byte, 32)
	bigLen := new(big.Int).SetInt64(int64(len(payload))).Bytes()
	copy(length[32-len(bigLen):], bigLen)

	padding := 32 - (len(payload) % 32)
	if padding == 32 {
		padding = 0
	}
	padded := append(payload, make([]byte, padding)...)

	out := make([]byte, 0, 4+32+32+len(padded))
	out = append(out, methodID...)
	out = append(out, offset...)
	out = append(out, length...)
	out = append(out, padded...)
	return out
}
