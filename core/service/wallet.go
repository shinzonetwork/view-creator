package service

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bip39 "github.com/cosmos/go-bip39"
	bip32 "github.com/tyler-smith/go-bip32"

	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	Mnemonic string `json:"mnemonic"`
	Address  string `json:"address"`
}

func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", fmt.Errorf("failed to generate entropy: %w", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}
	return mnemonic, nil
}

func DeriveEvmAddress(mnemonic string) (string, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", fmt.Errorf("invalid mnemonic phrase")
	}

	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", fmt.Errorf("failed to create master key: %w", err)
	}

	purpose, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		return "", fmt.Errorf("failed to derive purpose key: %w", err)
	}

	coinType, err := purpose.NewChildKey(bip32.FirstHardenedChild + 60)
	if err != nil {
		return "", fmt.Errorf("failed to derive coin type key: %w", err)
	}

	account, err := coinType.NewChildKey(bip32.FirstHardenedChild + 0)
	if err != nil {
		return "", fmt.Errorf("failed to derive account key: %w", err)
	}

	change, err := account.NewChildKey(0)
	if err != nil {
		return "", fmt.Errorf("failed to derive change key: %w", err)
	}

	addressIndexKey, err := change.NewChildKey(0)
	if err != nil {
		return "", fmt.Errorf("failed to derive address index key: %w", err)
	}

	privateKey, err := crypto.ToECDSA(addressIndexKey.Key)
	if err != nil {
		return "", fmt.Errorf("failed to convert to ECDSA: %w", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return address.Hex(), nil
}

func DerivePrivateKey(wallet Wallet) (*ecdsa.PrivateKey, error) {
	if !bip39.IsMnemonicValid(wallet.Mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic phrase")
	}

	seed := bip39.NewSeed(wallet.Mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key: %w", err)
	}

	purpose, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		return nil, fmt.Errorf("failed to derive purpose key: %w", err)
	}

	coinType, err := purpose.NewChildKey(bip32.FirstHardenedChild + 60)
	if err != nil {
		return nil, fmt.Errorf("failed to derive coin type key: %w", err)
	}

	account, err := coinType.NewChildKey(bip32.FirstHardenedChild + 0)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account key: %w", err)
	}

	change, err := account.NewChildKey(0)
	if err != nil {
		return nil, fmt.Errorf("failed to derive change key: %w", err)
	}

	addressIndexKey, err := change.NewChildKey(0)
	if err != nil {
		return nil, fmt.Errorf("failed to derive address index key: %w", err)
	}

	privateKey, err := crypto.ToECDSA(addressIndexKey.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to ECDSA: %w", err)
	}

	return privateKey, nil
}

func DeriveShinzoAddress(mnemonic string) (string, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", fmt.Errorf("invalid mnemonic phrase")
	}

	// Set Bech32 prefix to "shinzo"
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("shinzo", "shinzohub")

	hdPath := hd.CreateHDPath(118, 0, 0).String()
	privKeyBytes, err := hd.Secp256k1.Derive()(mnemonic, "", hdPath)
	if err != nil {
		return "", fmt.Errorf("failed to derive private key: %w", err)
	}

	privKey := secp256k1.PrivKey{Key: privKeyBytes}
	addr := sdk.AccAddress(privKey.PubKey().Address())

	return addr.String(), nil
}

func SaveWallet(mnemonic, address string) error {
	wallet := Wallet{Mnemonic: mnemonic, Address: address}

	data, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal wallet: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	dir := filepath.Join(home, ".shinzo", "wallet")

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create wallet dir: %w", err)
	}

	path := filepath.Join(dir, "wallet.json")
	return os.WriteFile(path, data, 0600)
}

func ImportMnemonic(mnemonic string) error {
	address, err := DeriveEvmAddress(mnemonic)
	if err != nil {
		return fmt.Errorf("failed to derive address: %w", err)
	}

	return SaveWallet(mnemonic, address)
}

func LoadWallet() (Wallet, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Wallet{}, fmt.Errorf("failed to get user home directory: %w", err)
	}

	dir := filepath.Join(home, ".shinzo", "wallet")
	path := filepath.Join(dir, "wallet.json")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Wallet{}, fmt.Errorf("wallet not found. Please generate one using `viewkit wallet generate`")
		}
		return Wallet{}, fmt.Errorf("failed to read wallet file: %w", err)
	}

	var wallet Wallet
	if err := json.Unmarshal(data, &wallet); err != nil {
		return Wallet{}, err
	}
	return wallet, nil
}

func GenerateWallet() (string, string, error) {
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		return "", "", err
	}

	address, err := DeriveEvmAddress(mnemonic)
	if err != nil {
		return "", "", err
	}

	if err := SaveWallet(mnemonic, address); err != nil {
		return "", "", err
	}

	return mnemonic, address, nil
}
