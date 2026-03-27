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
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/shinzonetwork/shinzo-view-creator/core/models"
	schemastore "github.com/shinzonetwork/shinzo-view-creator/core/schema/store"
	viewstore "github.com/shinzonetwork/shinzo-view-creator/core/view/store"
)

const defraVersion = "1.0.0-rc1"
const defraPort = "9181"

var defraCmd *exec.Cmd

type DefraViewPayload struct {
	Query        string  `json:"Query"`
	SDL          string  `json:"SDL"`
	TransformCID *string `json:"TransformCID,omitempty"`
}

type DefraLensModule struct {
	Path      string         `json:"Path"`
	Arguments map[string]any `json:"Arguments,omitempty"`
}

type DefraLensInner struct {
	Lenses []DefraLensModule `json:"Lenses"`
}

type DefraLensPayload struct {
	Lens DefraLensInner `json:"Lens"`
}

func defraURL() string {
	return "http://127.0.0.1:" + defraPort
}

func startDefraNode(debug bool) (string, error) {
	bin, err := EnsureDefraBinary(defraVersion)
	if err != nil {
		return "", fmt.Errorf("failed to ensure defradb binary: %w", err)
	}

	rootDir, err := os.MkdirTemp("", "defradb-root-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp rootdir: %w", err)
	}

	env := append(os.Environ(), "DEFRA_KEYRING_SECRET=1234")

	killExistingDefra()

	defraCmd = exec.Command(bin, "start", "--rootdir", rootDir)
	defraCmd.Env = env
	defraCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if debug {
		defraCmd.Stdout = os.Stdout
		defraCmd.Stderr = os.Stderr
	}

	if err := defraCmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start defradb: %w", err)
	}

	fmt.Println("🚀 DefraDB is running on port", defraPort)
	fmt.Println("⏳ Waiting for DefraDB to boot up...")
	time.Sleep(2 * time.Second)
	fmt.Println("✅ DefraDB booted up")

	return rootDir, nil
}

func registerLensAndCreateView(ctx context.Context, baseURL string, view models.View) (string, error) {
	var transformCID *string

	if len(view.Transform.Lenses) > 0 {
		lensPayload, err := BuildLensPayload(view)
		if err != nil {
			return "", fmt.Errorf("failed to build lens payload: %w", err)
		}

		fmt.Println("📦 Registering lens transform...")
		lensCID, err := RegisterLens(ctx, baseURL, lensPayload)
		if err != nil {
			return "", fmt.Errorf("failed to register lens: %w", err)
		}
		transformCID = &lensCID
	}

	viewPayload, err := BuildViewPayload(view, transformCID)
	if err != nil {
		return "", fmt.Errorf("failed to build view payload: %w", err)
	}

	return SendViewToDefra(ctx, baseURL, viewPayload)
}

func DeployViewToPlayground(name string, vs viewstore.ViewStore, ss schemastore.SchemaStore, url string) error {
	ctx := context.Background()

	fmt.Println("🔍 Loading view...")
	view, err := vs.Load(name)
	if err != nil {
		return fmt.Errorf("❌ Failed to load view: %w", err)
	}

	fmt.Println("📦 Applying schema...")
	schemaContent, err := ss.Load()
	if err != nil {
		return fmt.Errorf("❌ Failed to load schema: %w", err)
	}
	if err := ApplySchemaViaHTTP(ctx, url, schemaContent); err != nil {
		return fmt.Errorf("❌ Failed to apply schema: %w", err)
	}
	fmt.Println("✅ Schema applied")

	fmt.Println("🧠 Applying view...")
	result, err := registerLensAndCreateView(ctx, url, view)
	if err != nil {
		return fmt.Errorf("❌ Failed to apply view: %w", err)
	}
	fmt.Println("✅ View applied")

	collection, err := extractCollectionName(result)
	if err != nil {
		return fmt.Errorf("❌ Failed to extract collection name: %w", err)
	}

	fmt.Println("♻️  Refreshing view...")
	if err := RefreshView(ctx, url, collection); err != nil {
		return fmt.Errorf("❌ Failed to refresh view: %w", err)
	}
	fmt.Println("✅ View refreshed")

	fmt.Println("🎉 View deployed to playground at", url+"/")
	return nil
}

func StartLocalNodePlayground(schemastore schemastore.SchemaStore, debug bool) error {
	ctx := context.Background()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	if _, err := startDefraNode(debug); err != nil {
		return err
	}

	fmt.Println("⏳ Applying schema...")
	schemaContent, err := schemastore.Load()
	if err != nil {
		return cleanupDefra("failed to load schema", err)
	}

	if err := ApplySchemaViaHTTP(ctx, defraURL(), schemaContent); err != nil {
		return cleanupDefra("failed to apply schema", err)
	}
	fmt.Println("✅ Schema applied")

	if err := InsertMockData(ctx, defraURL()); err != nil {
		return cleanupDefra("failed to insert mock data", err)
	}

	fmt.Println("🧪 DefraDB playground ready at", defraURL()+"/")
	fmt.Println("📦 Press Ctrl+C to stop...")

	<-ctx.Done()
	return shutdownDefra()
}

func StartLocalNodeAndDeployView(name string, viewstore viewstore.ViewStore, schemastore schemastore.SchemaStore, debug bool) error {
	ctx := context.Background()

	view, err := viewstore.Load(name)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	if _, err := startDefraNode(debug); err != nil {
		return err
	}

	fmt.Println("⏳ Applying Schemas ...")
	schemaContent, err := schemastore.Load()
	if err != nil {
		return cleanupDefra("failed to load schema", err)
	}

	if err := ApplySchemaViaHTTP(ctx, defraURL(), schemaContent); err != nil {
		return cleanupDefra("failed to apply schema", err)
	}
	fmt.Println("✅ Schema Applied")

	if err := InsertMockData(ctx, defraURL()); err != nil {
		return cleanupDefra("failed to insert data", err)
	}

	fmt.Println("✅ Applying View ...")

	result, err := registerLensAndCreateView(ctx, defraURL(), view)
	if err != nil {
		return cleanupDefra("failed to send view", err)
	}

	collection, err := extractCollectionName(result)
	if err != nil {
		return cleanupDefra("failed to send view", err)
	}

	err = RefreshView(ctx, defraURL(), collection)
	if err != nil {
		return cleanupDefra("failed to send view", err)
	}

	fmt.Println("✅ View Successfully Applied")

	fmt.Println("🧪 Visit the DefraDB GraphQL Playground at", defraURL()+"/")
	fmt.Println("📦 Press Ctrl+C to stop...")

	<-ctx.Done()
	return shutdownDefra()
}

func StartLocalNodeAndTestView(name string, viewstore viewstore.ViewStore, schemastore schemastore.SchemaStore, debug bool) error {
	ctx := context.Background()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Println("🔍 Loading view...")
	view, err := viewstore.Load(name)
	if err != nil {
		return fmt.Errorf("❌ Failed to load view: %w", err)
	}

	if _, err := startDefraNode(debug); err != nil {
		return err
	}

	fmt.Println("📦 Applying schema...")
	schemaContent, err := schemastore.Load()
	if err != nil {
		return cleanupDefra("❌ Failed to load schema", err)
	}

	if err := ApplySchemaViaHTTP(ctx, defraURL(), schemaContent); err != nil {
		return cleanupDefra("❌ Failed to apply schema", err)
	}
	fmt.Println("✅ Schema applied")

	fmt.Println("📨 Inserting test data...")
	if err := InsertMockData(ctx, defraURL()); err != nil {
		return cleanupDefra("❌ Failed to insert data", err)
	}
	fmt.Println("✅ Data inserted")

	fmt.Println("🧠 Applying view...")
	result, err := registerLensAndCreateView(ctx, defraURL(), view)
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
	err = RefreshView(ctx, defraURL(), collection)
	if err != nil {
		return cleanupDefra("❌ Failed to refresh view", err)
	}
	fmt.Println("✅ View refreshed")

	fmt.Println("🔎 Querying view results...")
	fields := extractSDLFields(deref(view.Sdl))
	query := fmt.Sprintf(`{ %s { %s } }`, collection, strings.Join(fields, " "))
	queryResult, err := QueryDefra(ctx, defraURL(), query)
	if err != nil {
		return cleanupDefra("❌ Failed to query view", err)
	}
	fmt.Println("📊 View results:")
	fmt.Println(queryResult)

	fmt.Println("✅ Test flow completed successfully. Shutting down...")
	return shutdownDefra()
}

func BuildViewPayload(view models.View, transformCID *string) (string, error) {
	payload := DefraViewPayload{
		Query:        deref(view.Query),
		SDL:          deref(view.Sdl),
		TransformCID: transformCID,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func BuildLensPayload(view models.View) (string, error) {
	modules := make([]DefraLensModule, 0, len(view.Transform.Lenses))
	for _, lens := range view.Transform.Lenses {
		wasmBase64 := getPathInViewAssets(view.Name, lens.Path)
		args := make(map[string]any, len(lens.Arguments))
		for k, v := range lens.Arguments {
			args[k] = v
		}
		modules = append(modules, DefraLensModule{
			Path:      wasmBase64,
			Arguments: args,
		})
	}

	payload := DefraLensPayload{Lens: DefraLensInner{Lenses: modules}}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func RegisterLens(ctx context.Context, baseURL string, lensPayload string) (string, error) {
	url := baseURL + "/api/v0/lens"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(lensPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create lens request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send lens request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("lens registration failed (%d): %s", resp.StatusCode, string(body))
	}

	var lensResp struct {
		LensID string `json:"lensId"`
	}
	if err := json.Unmarshal(body, &lensResp); err != nil {
		return "", fmt.Errorf("failed to parse lens response: %w", err)
	}
	if lensResp.LensID == "" {
		return "", fmt.Errorf("no lensId in response: %s", string(body))
	}

	return lensResp.LensID, nil
}

func ApplySchemaViaHTTP(ctx context.Context, baseURL string, schema string) error {
	url := baseURL + "/api/v0/collections"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(schema))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusBadRequest && strings.Contains(string(body), "collection already exists") {
			return nil
		}
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func SendViewToDefra(ctx context.Context, baseURL string, jsonPayload string) (string, error) {
	url := baseURL + "/api/v0/view"

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
		if resp.StatusCode == http.StatusBadRequest {
			bodyStr := string(body)
			if idx := strings.Index(bodyStr, "collection already exists. Name: "); idx >= 0 {
				name := strings.TrimSpace(bodyStr[idx+len("collection already exists. Name: "):])
				name = strings.FieldsFunc(name, func(r rune) bool { return r == '"' || r == '}' || r == '\n' })[0]
				return fmt.Sprintf(`[{"Name":%q}]`, name), nil
			}
		}
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

func RefreshView(ctx context.Context, baseURL string, collection string) error {
	url := fmt.Sprintf("%s/api/v0/view/refresh?name=%s", baseURL, collection)

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

func extractCollectionName(result string) (string, error) {
	var parsed []map[string]any
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		return "", fmt.Errorf("failed to parse result: %w", err)
	}
	if len(parsed) == 0 {
		return "", fmt.Errorf("empty result")
	}

	if vRaw, ok := parsed[0]["version"]; ok && vRaw != nil {
		if v, ok := vRaw.(map[string]any); ok {
			if name, ok := v["Name"].(string); ok && name != "" {
				return name, nil
			}
			return "", fmt.Errorf("missing or invalid 'version.Name' field")
		}
		return "", fmt.Errorf("invalid 'version' type")
	}

	if name, ok := parsed[0]["Name"].(string); ok && name != "" {
		return name, nil
	}

	return "", fmt.Errorf("missing 'Name' in both formats")
}

func InsertMockData(ctx context.Context, baseURL string) error {
	fmt.Println("⏳ Inserting mock data (phase 1: blocks)...")
	if err := InsertDataToDefra(ctx, baseURL, GQL_BLOCKS); err != nil {
		return fmt.Errorf("failed to insert blocks: %w", err)
	}

	fmt.Println("⏳ Inserting mock data (phase 2: transactions)...")
	if err := InsertDataToDefra(ctx, baseURL, GQL_TRANSACTIONS); err != nil {
		return fmt.Errorf("failed to insert transactions: %w", err)
	}

	fmt.Println("⏳ Querying docIDs for relation linking...")
	blockDocIDs, err := queryDocIDsByHash(ctx, baseURL, "Ethereum__Mainnet__Block")
	if err != nil {
		return fmt.Errorf("failed to query block docIDs: %w", err)
	}
	txDocIDs, err := queryDocIDsByHash(ctx, baseURL, "Ethereum__Mainnet__Transaction")
	if err != nil {
		return fmt.Errorf("failed to query transaction docIDs: %w", err)
	}
	fmt.Printf("   blocks: %d, transactions: %d\n", len(blockDocIDs), len(txDocIDs))

	fmt.Println("⏳ Inserting mock data (phase 3: logs with relations)...")
	logsMutation := BuildLogsMutation(blockDocIDs, txDocIDs)
	if err := InsertDataToDefra(ctx, baseURL, logsMutation); err != nil {
		return fmt.Errorf("failed to insert logs: %w", err)
	}

	fmt.Println("✅ Mock data inserted")
	return nil
}

func queryDocIDsByHash(ctx context.Context, baseURL, typeName string) (map[string]string, error) {
	query := fmt.Sprintf(`{ %s { hash _docID } }`, typeName)
	result, err := QueryDefra(ctx, baseURL, query)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Data map[string][]struct {
			Hash  string `json:"hash"`
			DocID string `json:"_docID"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse query response: %w", err)
	}

	docIDs := make(map[string]string)
	for _, docs := range parsed.Data {
		for _, v := range docs {
			if v.Hash != "" && v.DocID != "" {
				docIDs[v.Hash] = v.DocID
			}
		}
	}
	return docIDs, nil
}

func InsertDataToDefra(ctx context.Context, baseURL string, data string) error {
	fmt.Println("⏳ Data Inserting...")

	reqBody := map[string]string{
		"query": data,
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(reqBody); err != nil {
		return fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/api/v0/graphql", buf)
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

	if bodyStr := string(body); strings.Contains(bodyStr, `"errors"`) {
		fmt.Println("⚠️  Insert response has errors:", bodyStr[:min(len(bodyStr), 500)])
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
	if defraCmd == nil || defraCmd.Process == nil {
		return nil
	}

	pid := defraCmd.Process.Pid
	pgid := -pid

	if err := syscall.Kill(pgid, syscall.SIGTERM); err != nil {
		fmt.Println("⚠️ Could not send SIGTERM to process group:", err)
	}

	done := make(chan error, 1)
	go func() { done <- defraCmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			fmt.Println("⚠️ DefraDB did not exit cleanly:", err)
		} else {
			fmt.Println("✅ DefraDB stopped.")
		}
	case <-time.After(5 * time.Second):
		fmt.Println("⚠️ DefraDB did not stop in time, force killing...")
		_ = syscall.Kill(pgid, syscall.SIGKILL)
		<-done
		fmt.Println("✅ DefraDB force killed.")
	}

	defraCmd = nil
	return nil
}

func killExistingDefra() {
	out, err := exec.Command("lsof", "-ti", "tcp:"+defraPort).Output()
	if err != nil || len(out) == 0 {
		return
	}
	for _, pidStr := range bytes.Split(bytes.TrimSpace(out), []byte("\n")) {
		fmt.Printf("⚠️ Killing leftover defradb process (pid %s) on port %s\n", pidStr, defraPort)
		_ = exec.Command("kill", "-9", string(pidStr)).Run()
	}
	time.Sleep(500 * time.Millisecond)
}

func cleanupDefra(reason string, err error) error {
	fmt.Println("❌", reason+":", err)
	_ = shutdownDefra()
	return fmt.Errorf("%s: %w", reason, err)
}

func defraDownloadURL(version string) string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	if arch == "amd64" {
		arch = "x86_64"
	}

	return fmt.Sprintf(
		"https://github.com/sourcenetwork/defradb/releases/download/v%s/defradb_%s_%s_%s",
		version, version, osName, arch,
	)
}

func QueryDefra(ctx context.Context, baseURL string, query string) (string, error) {
	reqBody := map[string]string{"query": query}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(reqBody); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/api/v0/graphql", buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
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

func extractSDLFields(sdl string) []string {
	braceIdx := strings.Index(sdl, "{")
	if braceIdx < 0 {
		return nil
	}
	body := sdl[braceIdx:]
	re := regexp.MustCompile(`(\w+)\s*:\s*(?:String|Int|Float|Boolean|ID|\[)`)
	matches := re.FindAllStringSubmatch(body, -1)
	var fields []string
	for _, m := range matches {
		fields = append(fields, m[1])
	}
	return fields
}
