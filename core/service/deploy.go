package service

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/shinzonetwork/view-creator/core/models"
	schemastore "github.com/shinzonetwork/view-creator/core/schema/store"
	viewstore "github.com/shinzonetwork/view-creator/core/view/store"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var SHINZO_HUB_PRECOMPILED_VIEW_REGISTRY_ADDRESS = "0x0000000000000000000000000000000000000210"
var DEFAULT_EVM_RPC = "http://127.0.0.1:8545/"

type ViewLite struct {
	Name      string           `json:"name"`
	Query     *string          `json:"query"`
	Sdl       *string          `json:"sdl"`
	Transform models.Transform `json:"transform"`
}

func StartLocalNodeTestAndDeploy(name string, viewstore viewstore.ViewStore, schemastore schemastore.SchemaStore, wallet Wallet) error {
	// err := StartLocalNodeAndTestView(name, viewstore, schemastore)
	// if err != nil {
	// 	return fmt.Errorf("view build and test failed locally")
	// }

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
		Name:      view.Name,
		Query:     view.Query,
		Sdl:       view.Sdl,
		Transform: view.Transform,
	}

	data, err := json.Marshal(viewLite)
	if err != nil {
		return fmt.Errorf("failed to marshal view: %w", err)
	}

	viewid := ComputeViewID(privateKey, data)
	hash, err := sendRegisterTx(DEFAULT_EVM_RPC, SHINZO_HUB_PRECOMPILED_VIEW_REGISTRY_ADDRESS, privateKey, data)
	if err != nil {
		return err
	}

	fmt.Println("✅ View deployment successful!")
	fmt.Println("----------------------------------------")
	fmt.Printf("🔑 View ID:           %s\n", viewid.Hex())
	fmt.Printf("📦 Transaction Hash:  %s\n", hash)
	fmt.Println("Blob size (bytes):", len(data))
	fmt.Println("----------------------------------------")
	fmt.Println("You can use the view ID and hash to query or verify registration on-chain.")

	return nil
}

func ComputeViewID(privateKey *ecdsa.PrivateKey, blob []byte) common.Hash {
	sender := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Concatenate sender address and blob (same as abi.encodePacked)
	combined := append(sender.Bytes(), blob...)

	// keccak256 hash
	viewID := crypto.Keccak256Hash(combined)

	return viewID
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
	tx := types.NewTransaction(nonce, to, big.NewInt(0), 5000000000, gasPrice, input)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign tx: %w", err)
	}

	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		return "", fmt.Errorf("failed to send tx: %w", err)
	}

	log.Printf("✅ View deployed! Tx hash: %s", signedTx.Hash().Hex())

	return signedTx.Hash().Hex(), nil
}
