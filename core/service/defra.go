package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/shinzonetwork/shinzo-view-creator/core/models"
	schemastore "github.com/shinzonetwork/shinzo-view-creator/core/schema/store"
	viewstore "github.com/shinzonetwork/shinzo-view-creator/core/view/store"
)

var defraCmd *exec.Cmd

type DefraViewPayload struct {
	Query     string         `json:"Query"`
	SDL       string         `json:"SDL"`
	Transform map[string]any `json:"Transform"`
}

func StartLocalNodeAndDeployView(name string, viewstore viewstore.ViewStore, schemastore schemastore.SchemaStore, debug bool) error {
	ctx := context.Background()

	view, err := viewstore.Load(name)
	if err != nil {
		return err
	}

	viewJson, err := ConvertViewToDefraJson(view)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := "9181"

	bin, err := EnsureDefraBinary("0.18.0")
	if err != nil {
		return fmt.Errorf("failed to ensure defradb binary: %w", err)
	}

	rootDir, err := os.MkdirTemp("", "defradb-root-*")
	if err != nil {
		return fmt.Errorf("failed to create temp rootdir: %w", err)
	}

	env := append(os.Environ(),
		"DEFRA_KEYRING_SECRET=1234",
	)

	defraCmd = exec.Command(bin, "start", "--rootdir", rootDir)
	defraCmd.Env = env

	if debug {
		defraCmd.Stdout = os.Stdout
		defraCmd.Stderr = os.Stderr
	}

	if err := defraCmd.Start(); err != nil {
		return fmt.Errorf("failed to start defradb: %w", err)
	}

	fmt.Println("🚀 DefraDB is running on port", port)
	fmt.Println("⏳ Waiting for DefraDB to boot up...")
	time.Sleep(2 * time.Second)
	fmt.Println("✅ DefraDB booted up")

	fmt.Println("⏳ Applying Schemas ...")
	schemaContent, err := schemastore.Load()
	if err != nil {
		return cleanupDefra("failed to load schema", err)
	}

	schemaCmd := exec.Command(bin, "client", "schema", "add", schemaContent, "--rootdir", rootDir)
	schemaCmd.Env = env

	if err := schemaCmd.Run(); err != nil {
		return cleanupDefra("failed to apply schema", err)
	}
	fmt.Println("✅ Schema Applied")

	if err := InsertDataToDefra(ctx, GQL); err != nil {
		return cleanupDefra("failed to insert data", err)
	}

	fmt.Println("✅ Applying View ...")

	result, err := SendViewToDefra(ctx, "http://127.0.0.1:9181", viewJson)
	if err != nil {
		return cleanupDefra("failed to send view", err)
	}

	collection, err := extractCollectionName(result)
	if err != nil {
		return cleanupDefra("failed to send view", err)
	}

	err = RefreshView(ctx, "http://127.0.0.1:9181", collection)
	if err != nil {
		return cleanupDefra("failed to send view", err)
	}

	fmt.Println("✅ View Successfully Applied")

	fmt.Println("🧪 Visit the DefraDB GraphQL Playground at http://127.0.0.1:9181/")
	fmt.Println("📦 Press Ctrl+C to stop...")

	<-ctx.Done()
	return shutdownDefra()
}

func StartLocalNodeAndTestView(name string, viewstore viewstore.ViewStore, schemastore schemastore.SchemaStore, debug bool) error {
	ctx := context.Background()

	fmt.Println("🔍 Loading view...")
	view, err := viewstore.Load(name)
	if err != nil {
		return fmt.Errorf("❌ Failed to load view: %w", err)
	}

	viewJson, err := ConvertViewToDefraJson(view)
	if err != nil {
		return fmt.Errorf("❌ Failed to convert view to JSON: %w", err)
	}

	fmt.Println("⚙️  Ensuring DefraDB binary...")
	bin, err := EnsureDefraBinary("0.18.0")
	if err != nil {
		return fmt.Errorf("❌ Failed to ensure DefraDB binary: %w", err)
	}

	fmt.Println("📁 Creating temporary root directory...")
	rootDir, err := os.MkdirTemp("", "defradb-root-*")
	if err != nil {
		return fmt.Errorf("❌ Failed to create temp rootdir: %w", err)
	}

	env := append(os.Environ(),
		"DEFRA_KEYRING_SECRET=1234",
	)

	fmt.Println("🚀 Starting DefraDB...")
	defraCmd = exec.Command(bin, "start", "--rootdir", rootDir)

	// This is here for debug purposes; Show command output as it happens
	if debug {
		defraCmd.Stdout = os.Stdout
		defraCmd.Stderr = os.Stderr
	}

	defraCmd.Env = env

	if err := defraCmd.Start(); err != nil {
		return fmt.Errorf("❌ Failed to start DefraDB: %w", err)
	}

	fmt.Println("⏳ Waiting for DefraDB to boot...")
	time.Sleep(2 * time.Second)
	fmt.Println("✅ DefraDB booted")

	fmt.Println("📦 Applying schema...")
	schemaContent, err := schemastore.Load()
	if err != nil {
		return cleanupDefra("❌ Failed to load schema", err)
	}

	schemaCmd := exec.Command(bin, "client", "schema", "add", schemaContent, "--rootdir", rootDir)
	schemaCmd.Env = env

	// This is here for debug purposes; Show command output as it happens
	// schemaCmd.Stdout = os.Stdout
	// schemaCmd.Stderr = os.Stderr

	if err := schemaCmd.Run(); err != nil {
		return cleanupDefra("❌ Failed to apply schema", err)
	}
	fmt.Println("✅ Schema applied")

	fmt.Println("📨 Inserting test data...")
	if err := InsertDataToDefra(ctx, GQL); err != nil {
		return cleanupDefra("❌ Failed to insert data", err)
	}
	fmt.Println("✅ Data inserted")

	fmt.Println("🧠 Applying view...")
	result, err := SendViewToDefra(ctx, "http://127.0.0.1:9181", viewJson)
	if err != nil {
		return cleanupDefra("❌ Failed to apply view", err)
	}
	fmt.Println("✅ View applied")

	fmt.Println("🔎 Extracting collection name...")
	collection, err := extractCollectionName(result)
	if err != nil {
		return cleanupDefra("❌ Failed to extract collection name", err)
	}

	fmt.Println("♻️  Refreshing view...")
	err = RefreshView(ctx, "http://127.0.0.1:9181", collection)
	if err != nil {
		return cleanupDefra("❌ Failed to refresh view", err)
	}
	fmt.Println("✅ View refreshed")

	fmt.Println("✅ Test flow completed successfully. Shutting down...")
	return shutdownDefra()
}

func ConvertViewToDefraJson(view models.View) (string, error) {
	transform := map[string]any{
		"lenses": []map[string]any{},
	}

	for _, lens := range view.Transform.Lenses {
		lensMap := map[string]any{
			"path":      getPathInViewAssets(view.Name, lens.Path),
			"arguments": lens.Arguments,
		}
		transform["lenses"] = append(transform["lenses"].([]map[string]any), lensMap)
	}

	payload := DefraViewPayload{
		Query:     deref(view.Query),
		SDL:       deref(view.Sdl),
		Transform: transform,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func SendViewToDefra(ctx context.Context, defraURL string, jsonPayload string) (string, error) {
	url := defraURL + "/api/v0/view"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

func RefreshView(ctx context.Context, defraURL string, collection string) error {
	url := fmt.Sprintf("%s/api/v0/view/refresh?name=%s", defraURL, collection)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.Header.Set("Accept", "*/*")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send refresh request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected refresh status %d", resp.StatusCode)
	}

	return nil
}

// func extractCollectionName(result string) (string, error) {
// 	var parsed []map[string]interface{}
// 	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
// 		return "", fmt.Errorf("failed to parse result: %w", err)
// 	}

// 	if len(parsed) == 0 {
// 		return "", fmt.Errorf("empty result")
// 	}

// 	version, ok := parsed[0]["version"].(map[string]interface{})
// 	if !ok {
// 		return "", fmt.Errorf("missing or invalid 'version' field")
// 	}

// 	name, ok := version["Name"].(string)
// 	if !ok {
// 		return "", fmt.Errorf("missing or invalid 'Name' field")
// 	}

// 	return name, nil
// }

func extractCollectionName(result string) (string, error) {
	var parsed []map[string]any
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		return "", fmt.Errorf("failed to parse result: %w", err)
	}
	if len(parsed) == 0 {
		return "", fmt.Errorf("empty result")
	}

	// Format A: { "version": { "Name": ... }, "schema": ... }
	if vRaw, ok := parsed[0]["version"]; ok && vRaw != nil {
		if v, ok := vRaw.(map[string]any); ok {
			if name, ok := v["Name"].(string); ok && name != "" {
				return name, nil
			}
			return "", fmt.Errorf("missing or invalid 'version.Name' field")
		}
		return "", fmt.Errorf("invalid 'version' type")
	}

	// Format B: { "Name": "...", "VersionID": "...", ... }
	if name, ok := parsed[0]["Name"].(string); ok && name != "" {
		return name, nil
	}

	return "", fmt.Errorf("missing 'Name' in both formats")
}

func InsertDataToDefra(ctx context.Context, data string) error {
	fmt.Println("⏳ Data Inserting...")

	reqBody := map[string]string{
		"query": data,
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(reqBody); err != nil {
		return fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "http://127.0.0.1:9181/api/v0/graphql", buf)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("✅ Data Inserted Successfully")
	return nil
}

func EnsureDefraBinary(version string, dir ...string) (string, error) {
	var base string
	if len(dir) == 0 || dir[0] == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("unable to get home directory: %w", err)
		}
		base = filepath.Join(home, ".shinzo", "defra")
	} else {
		base = filepath.Join(dir[0], ".shinzo", "defra")
	}

	binaryPath := filepath.Join(base, "defradb")

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		if err := DownloadDefraDB(version, dir...); err != nil {
			return "", fmt.Errorf("failed to download defradb: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("error checking defradb binary: %w", err)
	}

	return binaryPath, nil
}

func DownloadDefraDB(version string, dir ...string) error {
	var base string

	if len(dir) == 0 || dir[0] == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("unable to get home directory: %w", err)
		}
		base = filepath.Join(home, ".shinzo", "defra")
	} else {
		base = filepath.Join(dir[0], ".shinzo", "defra")
	}

	if err := os.MkdirAll(base, 0755); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	url := defraDownloadURL(version)
	fmt.Println("Downloading DefraDB from:", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
	    body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
	    return fmt.Errorf("failed to download defradb (%s): status %d: %s", url, resp.StatusCode, string(body))
	}


	binary := filepath.Join(base, "defradb")
	out, err := os.Create(binary)
	if err != nil {
		return fmt.Errorf("failed to create binary file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write binary: %w", err)
	}

	if err := os.Chmod(binary, 0755); err != nil {
		return fmt.Errorf("failed to chmod binary: %w", err)
	}

	fmt.Println("DefraDB downloaded to:", binary)
	return nil
}

func DeleteDefraDB(dir ...string) error {
	var base string

	if len(dir) == 0 || dir[0] == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("unable to get home directory: %w", err)
		}
		base = filepath.Join(home, ".shinzo", "defra")
	} else {
		base = filepath.Join(dir[0], ".shinzo", "defra")
	}

	binary := filepath.Join(base, "defradb")

	if _, err := os.Stat(binary); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(binary)
}

func shutdownDefra() error {
	if defraCmd != nil && defraCmd.Process != nil {
		if err := defraCmd.Process.Signal(syscall.SIGTERM); err != nil {
			fmt.Println("⚠️ Could not send SIGTERM:", err)
		} else if err := defraCmd.Wait(); err != nil {
			fmt.Println("⚠️ DefraDB did not exit cleanly:", err)
		} else {
			fmt.Println("✅ DefraDB stopped.")
		}
	}
	return nil
}

func cleanupDefra(reason string, err error) error {
	fmt.Println("❌", reason+":", err)
	_ = shutdownDefra()
	return fmt.Errorf("%s: %w", reason, err)
}

func defraDownloadURL(version string) string {
    osName := runtime.GOOS
    arch := runtime.GOARCH

    // DefraDB v0.18.0 release assets use x86_64 (not amd64) for Linux/Windows filenames.
    if arch == "amd64" {
        arch = "x86_64"
    }

    return fmt.Sprintf(
        "https://github.com/sourcenetwork/defradb/releases/download/v%s/defradb_%s_%s_%s",
        version, version, osName, arch,
    )
}


func getPathInViewAssets(viewName, relativePath string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	fullPath := filepath.Join(home, ".shinzo", "views", viewName, relativePath)
	return "file://" + fullPath
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
