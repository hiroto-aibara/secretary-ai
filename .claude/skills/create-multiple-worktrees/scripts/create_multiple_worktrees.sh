#!/bin/bash
set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

log_dry() {
    echo -e "${CYAN}[DRY-RUN]${NC} $1"
}

# Default options
DRY_RUN=false
NO_SETUP=false
BASE_BRANCH=""
TASK_FILES=()

# Usage
usage() {
    cat << EOF
Usage: $0 --base <branch> [OPTIONS] <task-file.md> [task-file.md...]

Creates git worktrees from TASK.md files for parallel feature development.

Required:
  --base <branch>   Base branch to create worktrees from (e.g., main, develop)

Arguments:
  task-file.md      Path to TASK.md files (feature name is extracted from filename)

Options:
  --no-setup        Skip environment setup (mise install, make setup)
  --dry-run         Show what would be created without actually creating
  -h, --help        Show this help message

Examples:
  $0 --base main tasks/user-auth.md tasks/dashboard.md
  $0 --base develop tasks/*.md
  $0 --base main --dry-run tasks/user-auth.md
  $0 --base main --no-setup tasks/api-v2.md

EOF
    exit 1
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --base)
            if [[ -z "${2:-}" ]]; then
                log_error "--base requires a branch name"
                exit 1
            fi
            BASE_BRANCH="$2"
            shift 2
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --no-setup)
            NO_SETUP=true
            shift
            ;;
        -h|--help)
            usage
            ;;
        -*)
            log_error "Unknown option: $1"
            usage
            ;;
        *)
            # Check if file exists
            if [[ -f "$1" ]]; then
                TASK_FILES+=("$1")
            else
                log_error "File not found: $1"
                exit 1
            fi
            shift
            ;;
    esac
done

# Check if --base is provided (required)
if [[ -z "$BASE_BRANCH" ]]; then
    log_error "--base option is required. Specify the base branch (e.g., --base main)"
    echo ""
    usage
fi

# Check if we have any task files
if [[ ${#TASK_FILES[@]} -eq 0 ]]; then
    log_error "No TASK.md files provided"
    usage
fi

# Get the root directory of the repository
REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || echo "")
if [[ -z "$REPO_ROOT" ]]; then
    log_error "Not a git repository"
    exit 1
fi
cd "$REPO_ROOT"

# Check if we're in a git repository
if [[ ! -d ".git" ]]; then
    log_error "Not a git repository"
    exit 1
fi

# Verify base branch exists
if ! git show-ref --verify --quiet "refs/heads/${BASE_BRANCH}" && \
   ! git show-ref --verify --quiet "refs/remotes/origin/${BASE_BRANCH}"; then
    log_error "Base branch '${BASE_BRANCH}' does not exist"
    exit 1
fi

# Function to extract feature name from filename
extract_feature_name() {
    local filepath="$1"
    local filename=$(basename "$filepath")
    # Remove .md extension
    echo "${filename%.md}"
}

# Function to generate random port (range: 10000-60000)
generate_random_port() {
    echo $((RANDOM % 50000 + 10000))
}

# Function to copy file if it exists
copy_if_exists() {
    local src="$1"
    local dest="$2"
    if [[ -f "${src}" ]]; then
        mkdir -p "$(dirname "$dest")"
        cp "${src}" "${dest}"
        return 0
    fi
    return 1
}

# Function to copy environment files using glob pattern
copy_env_files() {
    local worktree_dir="$1"
    local copied_count=0

    log_info "Detecting environment files..."

    # Find all .env* files up to 3 levels deep, excluding node_modules and .venv
    while IFS= read -r -d '' env_file; do
        # Skip if file is in excluded directories
        if [[ "$env_file" == *"node_modules"* ]] || \
           [[ "$env_file" == *".venv"* ]] || \
           [[ "$env_file" == *".worktrees"* ]] || \
           [[ "$env_file" == *".git"* ]]; then
            continue
        fi

        # Get relative path
        local rel_path="${env_file#./}"
        local dest_path="${worktree_dir}/${rel_path}"

        # Create directory and copy file
        mkdir -p "$(dirname "$dest_path")"

        # For root .env file, randomize ports
        if [[ "$rel_path" == ".env" ]]; then
            local random_frontend_port=$(generate_random_port)
            local random_backend_port=$(generate_random_port)
            local random_agent_port=$(generate_random_port)

            sed -e "s/^FRONTEND_PORT=.*/FRONTEND_PORT=${random_frontend_port}/" \
                -e "s/^BACKEND_PORT=.*/BACKEND_PORT=${random_backend_port}/" \
                -e "s/^AGENT_PORT=.*/AGENT_PORT=${random_agent_port}/" \
                -e "s/^PORT=.*/PORT=${random_frontend_port}/" \
                "$env_file" > "$dest_path"
        else
            cp "$env_file" "$dest_path"
        fi

        ((copied_count++))
        log_info "  Copied: ${rel_path}"
    done < <(find . -maxdepth 3 -name ".env*" -type f -print0 2>/dev/null)

    # Also copy .envrc if exists
    if [[ -f ".envrc" ]]; then
        cp ".envrc" "${worktree_dir}/.envrc"
        ((copied_count++))
        log_info "  Copied: .envrc"
    fi

    if [[ $copied_count -eq 0 ]]; then
        log_warn "No environment files found"
    else
        log_info "Copied ${copied_count} environment file(s)"
    fi
}

# Function to run environment setup
run_setup() {
    local worktree_dir="$1"
    local feature_name="$2"

    if [[ "$NO_SETUP" == true ]]; then
        log_info "Skipping setup (--no-setup specified)"
        return 0
    fi

    log_info "Running environment setup..."

    # Check for mise.toml or .mise.toml
    if [[ -f "${worktree_dir}/mise.toml" ]] || [[ -f "${worktree_dir}/.mise.toml" ]]; then
        log_info "  Running mise install..."
        (cd "${worktree_dir}" && mise install 2>&1) || log_warn "mise install completed with warnings for ${feature_name}"
    fi

    # Run make setup if Makefile with setup target exists
    if [[ -f "${worktree_dir}/Makefile" ]]; then
        if grep -q "^setup:" "${worktree_dir}/Makefile" 2>/dev/null; then
            log_info "  Running make setup..."
            (cd "${worktree_dir}" && make setup 2>&1) || log_warn "make setup completed with warnings for ${feature_name}"
        fi
    fi
}

# Function to create a single worktree
create_single_worktree() {
    local task_file="$1"
    local feature_name=$(extract_feature_name "$task_file")
    local branch_name="feature/${feature_name}"
    local worktree_dir=".worktrees/${feature_name}"

    # Check if worktree already exists
    if [[ -d "${worktree_dir}" ]]; then
        log_warn "Worktree already exists: ${worktree_dir} (skipping)"
        return 1
    fi

    # Create worktree
    if git show-ref --verify --quiet "refs/heads/${branch_name}"; then
        log_warn "Branch ${branch_name} already exists, using it"
        git worktree add "${worktree_dir}" "${branch_name}" 2>/dev/null || {
            log_error "Failed to create worktree for ${feature_name}"
            return 1
        }
    else
        git worktree add -b "${branch_name}" "${worktree_dir}" "${BASE_BRANCH}" 2>/dev/null || {
            log_error "Failed to create worktree for ${feature_name}"
            return 1
        }
    fi

    # Copy TASK.md to worktree
    cp "$task_file" "${worktree_dir}/TASK.md"
    log_info "TASK.md copied to ${worktree_dir}"

    # Copy environment files using glob detection
    copy_env_files "${worktree_dir}"

    # Run environment setup
    run_setup "${worktree_dir}" "${feature_name}"

    log_success "${feature_name} created (${worktree_dir})"
    return 0
}

# Dry run mode
if [[ "$DRY_RUN" == true ]]; then
    echo ""
    log_dry "Would create the following worktrees:"
    log_dry "Base branch: ${BASE_BRANCH}"
    echo ""

    # Show detected env files
    log_dry "Detected environment files:"
    while IFS= read -r -d '' env_file; do
        if [[ "$env_file" != *"node_modules"* ]] && \
           [[ "$env_file" != *".venv"* ]] && \
           [[ "$env_file" != *".worktrees"* ]] && \
           [[ "$env_file" != *".git"* ]]; then
            echo "    ${env_file#./}"
        fi
    done < <(find . -maxdepth 3 -name ".env*" -type f -print0 2>/dev/null)
    echo ""

    for task_file in "${TASK_FILES[@]}"; do
        feature_name=$(extract_feature_name "$task_file")
        echo "  - .worktrees/${feature_name}/"
        echo "    ├── TASK.md (from ${task_file})"
        echo "    ├── .env* (auto-detected)"
        echo "    └── (branch: feature/${feature_name} from ${BASE_BRANCH})"
    done
    echo ""
    log_dry "Total: ${#TASK_FILES[@]} worktrees"
    echo ""
    exit 0
fi

# Main execution
echo ""
log_info "Creating ${#TASK_FILES[@]} worktrees from TASK.md files..."
log_info "Base branch: ${BASE_BRANCH}"
echo ""

SUCCESS_COUNT=0
FAILED_COUNT=0
CREATED_WORKTREES=()

for task_file in "${TASK_FILES[@]}"; do
    feature_name=$(extract_feature_name "$task_file")
    log_step "Creating worktree: ${feature_name}"
    if create_single_worktree "$task_file"; then
        ((SUCCESS_COUNT++))
        CREATED_WORKTREES+=("$feature_name")
    else
        ((FAILED_COUNT++))
    fi
    echo ""
done

# Print summary
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN} ${SUCCESS_COUNT} worktrees created successfully!${NC}"
if [[ $FAILED_COUNT -gt 0 ]]; then
    echo -e "${YELLOW} ${FAILED_COUNT} worktrees failed or skipped${NC}"
fi
echo -e "${GREEN}========================================${NC}"
echo ""

for feature in "${CREATED_WORKTREES[@]}"; do
    echo ".worktrees/${feature}/"
    echo "  ├── TASK.md"
    echo "  ├── .env* (auto-detected)"
    echo "  └── feature/${feature} (from ${BASE_BRANCH})"
done

echo ""
echo "To start working on a feature:"
echo "  cd .worktrees/<feature-name>"
echo "  cat TASK.md    # Review task details"
echo "  claude         # Start Claude Code session"
echo ""
echo "To remove all worktrees when done:"
echo "  for wt in .worktrees/*/; do git worktree remove \"\$wt\"; done"
echo ""
