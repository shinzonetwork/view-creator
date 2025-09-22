# view-creator

A tool for creating views in DefraDB using LensVM.

## Installation

```bash
git clone https://github.com/shinzonetwork/view-creator.git
cd view-creator
go build -o viewkit ./cmd/viewkit
```

## Quick Start

```bash
# Get help for all commands
./viewkit --help

# Get help for view commands
./viewkit view --help

# Initialize a new view
./viewkit view init my-view

# Add a query to the view
./viewkit view add query 'query { users { name email } }' --name my-view

# Deploy the view
./viewkit view deploy my-view --target local
```

## Command Reference

### Main Commands

#### `./viewkit view init [name] [flags]`
Initialize a new view

**Flags:**
- `--json`: Output the view in raw JSON format
- `--verbose`: Show full output including revision history
- `-h, --help`: Help for init

**Example:**
```bash
./viewkit view init my-new-view
./viewkit view init my-view --json --verbose
```

---

#### `./viewkit view add [command]`
Add components to an existing view

**Subcommands:**
- `lens`: Add lenses in a view
- `query`: Add or update the query of the view  
- `sdl`: Add or update the SDL of the view

**Global Flags:**
- `--name string`: Name of the view

##### `./viewkit view add lens [flags]`
Add lenses in a view

**Flags:**
- `--args string`: Arguments of the lens transform
- `--label string`: Name of lens
- `--path string`: Path to the WASM file (local)
- `--url string`: URL to download the WASM file from
- `--name string`: Name of the view

**Example:**
```bash
./viewkit view add lens --name my-view --label user-filter --path ./lens.wasm --args "filter=active"
./viewkit view add lens --name my-view --label remote-lens --url https://example.com/lens.wasm
```

##### `./viewkit view add query '<query>' [flags]`
Add or update the query of the view

**Global Flags:**
- `--name string`: Name of the view

**Example:**
```bash
./viewkit view add query 'query { users { name email } }' --name my-view
```

##### `./viewkit view add sdl '<sdl>' [flags]`
Add or update the SDL of the view

**Global Flags:**
- `--name string`: Name of the view

**Example:**
```bash
./viewkit view add sdl 'type User { name: String email: String }' --name my-view
```

---

#### `./viewkit view remove [command]`
Remove components from an existing view

**Subcommands:**
- `lens`: Remove a lens from the view
- `query`: Remove the query from the view
- `sdl`: Remove the SDL from the view

**Global Flags:**
- `--name string`: Name of the view

##### `./viewkit view remove lens [flags]`
Remove a lens from the view

**Flags:**
- `--label string`: Label of the lens to remove
- `--name string`: Name of the view

**Example:**
```bash
./viewkit view remove lens --name my-view --label user-filter
```

##### `./viewkit view remove query [flags]`
Remove the query from the view

**Global Flags:**
- `--name string`: Name of the view

**Example:**
```bash
./viewkit view remove query --name my-view
```

##### `./viewkit view remove sdl [flags]`
Remove the SDL from the view

**Global Flags:**
- `--name string`: Name of the view

**Example:**
```bash
./viewkit view remove sdl --name my-view
```

---

#### `./viewkit view deploy <name> [flags]`
Deploy a view to local, devnet, or mainnet

**Flags:**
- `--target string`: Where to deploy the view: local, devnet, or mainnet (required)
- `-h, --help`: Help for deploy

**Example:**
```bash
./viewkit view deploy my-view --target local
./viewkit view deploy my-view --target devnet
./viewkit view deploy my-view --target mainnet
```

---

#### `./viewkit view inspect [name] [flags]`
Inspect a saved view

**Flags:**
- `--json`: Output the view in raw JSON format
- `--verbose`: Show full output including revision history
- `-h, --help`: Help for inspect

**Example:**
```bash
./viewkit view inspect my-view
./viewkit view inspect my-view --json
./viewkit view inspect my-view --verbose
```

---

#### `./viewkit view rollback <viewName> [flags]`
Rollback a view to a specific version (or previous if not specified)

**Flags:**
- `--version int`: Target version to rollback to (default -1)
- `-h, --help`: Help for rollback

**Example:**
```bash
./viewkit view rollback my-view
./viewkit view rollback my-view --version 2
```

---

#### `./viewkit view delete [name] [flags]`
Delete a saved view

**Flags:**
- `-h, --help`: Help for delete

**Example:**
```bash
./viewkit view delete my-view
```

---

#### `./viewkit view test <name> [flags]`
Test if the view can build and compile successfully

**Flags:**
- `-h, --help`: Help for test

**Example:**
```bash
./viewkit view test my-view
```

---

### Developer Tools

#### `./viewkit tools schema [command]`
Manage custom models in the Viewkit schema

**Subcommands:**
- `add`: Add a custom schema type to the viewkit schema
- `inspect`: Show the full definition of a schema type
- `list`: List all schema types in the Viewkit schema (default and custom)
- `remove`: Remove a custom schema type from the viewkit schema
- `reset`: Clear all custom schema types (does not affect defaults)
- `update`: Update the default schemas from a remote source

##### `./viewkit tools schema add <schema> [flags]`
Add a custom schema type to the viewkit schema

**Flags:**
- `-h, --help`: Help for add

**Example:**
```bash
./viewkit tools schema add 'type CustomUser { id: ID! name: String! }'
```

##### `./viewkit tools schema inspect <type> [flags]`
Show the full definition of a schema type

**Flags:**
- `-h, --help`: Help for inspect

**Example:**
```bash
./viewkit tools schema inspect User
./viewkit tools schema inspect CustomUser
```

##### `./viewkit tools schema list [flags]`
List all schema types in the Viewkit schema (default and custom)

**Flags:**
- `-h, --help`: Help for list

**Example:**
```bash
./viewkit tools schema list
```

##### `./viewkit tools schema remove <type> [flags]`
Remove a custom schema type from the viewkit schema

**Flags:**
- `-h, --help`: Help for remove

**Example:**
```bash
./viewkit tools schema remove CustomUser
```

##### `./viewkit tools schema reset [flags]`
Clear all custom schema types (does not affect defaults)

**Flags:**
- `-h, --help`: Help for reset

**Example:**
```bash
./viewkit tools schema reset
```

##### `./viewkit tools schema update [flags]`
Update the default schemas from a remote source

**Flags:**
- `--version string`: Git branch or tag to fetch the default schema from (default: main)
- `-h, --help`: Help for update

**Example:**
```bash
./viewkit tools schema update
./viewkit tools schema update --version v1.2.0
./viewkit tools schema update --version development
```

---

### Shell Completion

#### `./viewkit completion [command]`
Generate autocompletion scripts for various shells

**Subcommands:**
- `bash`: Generate the autocompletion script for bash
- `fish`: Generate the autocompletion script for fish
- `powershell`: Generate the autocompletion script for powershell
- `zsh`: Generate the autocompletion script for zsh

##### `./viewkit completion bash [flags]`
Generate the autocompletion script for the bash shell

**Flags:**
- `--no-descriptions`: Disable completion descriptions
- `-h, --help`: Help for bash

**Setup Instructions:**
```bash
# Load completions in current session
source <(viewkit completion bash)

# Load completions for every new session (Linux)
./viewkit completion bash > /etc/bash_completion.d/viewkit

# Load completions for every new session (macOS)
./viewkit completion bash > $(brew --prefix)/etc/bash_completion.d/viewkit
```

##### `./viewkit completion zsh [flags]`
Generate the autocompletion script for the zsh shell

**Flags:**
- `--no-descriptions`: Disable completion descriptions
- `-h, --help`: Help for zsh

**Setup Instructions:**
```bash
# Enable shell completion (if not already enabled)
echo "autoload -U compinit; compinit" >> ~/.zshrc

# Load completions in current session
source <(viewkit completion zsh)

# Load completions for every new session (Linux)
./viewkit completion zsh > "${fpath[1]}/_viewkit"

# Load completions for every new session (macOS)
./viewkit completion zsh > $(brew --prefix)/share/zsh/site-functions/_viewkit
```

##### `./viewkit completion fish [flags]`
Generate the autocompletion script for the fish shell

**Flags:**
- `--no-descriptions`: Disable completion descriptions
- `-h, --help`: Help for fish

**Setup Instructions:**
```bash
# Load completions in current session
./viewkit completion fish | source

# Load completions for every new session
./viewkit completion fish > ~/.config/fish/completions/viewkit.fish
```

##### `./viewkit completion powershell [flags]`
Generate the autocompletion script for powershell

**Flags:**
- `--no-descriptions`: Disable completion descriptions
- `-h, --help`: Help for powershell

**Setup Instructions:**
```powershell
# Load completions in current session
viewkit completion powershell | Out-String | Invoke-Expression

# For persistent completions, add the above command to your PowerShell profile
```

---

### Additional Commands

#### `./viewkit help [command] [flags]`
Get help for any command in the application

**Flags:**
- `-h, --help`: Help for help

**Example:**
```bash
./viewkit help
./viewkit help view
./viewkit help view init
./viewkit help tools schema add
```

## Workflow Examples

### Creating a Complete View

```bash
# 1. Initialize a new view
./viewkit view init user-analytics

# 2. Add SDL schema
./viewkit view add sdl 'type User { id: ID! name: String! email: String! createdAt: DateTime! }' --name user-analytics

# 3. Add a query
./viewkit view add query 'query { users(filter: {createdAt: {gt: "2024-01-01"}}) { name email createdAt } }' --name user-analytics

# 4. Add a lens for data transformation
./viewkit view add lens --name user-analytics --label date-formatter --path ./formatters/date.wasm

# 5. Test the view
./viewkit view test user-analytics

# 6. Deploy to local environment
./viewkit view deploy user-analytics --target local

# 7. Inspect the deployed view
./viewkit view inspect user-analytics --verbose
```

### Managing View Versions

```bash
# View current state
./viewkit view inspect my-view

# Make changes
./viewkit view add query 'query { users { name } }' --name my-view

# Test changes
./viewkit view test my-view

# If something goes wrong, rollback
./viewkit view rollback my-view

# Or rollback to specific version
./viewkit view rollback my-view --version 1
```

## Global Flags

All commands support:
- `-h, --help`: Show help for the command
- `--json`: Output in JSON format (where applicable)
- `--verbose`: Show detailed output (where applicable)
