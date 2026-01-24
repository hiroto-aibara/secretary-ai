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

log_section() {
    echo -e "${CYAN}$1${NC}"
}

# Default values
FORCE=false
TITLE=""
BODY=""
BASE_BRANCH="main"
DRAFT_FLAG=""
PR_URL=""
CURRENT_BRANCH=""
REPO_ROOT=""

# Usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Creates a PR from the current branch (worktree is NOT deleted)."
    echo ""
    echo "Options:"
    echo "  --title <title>        PR title (interactive if not specified)"
    echo "  --body <body>          PR body (interactive if not specified)"
    echo "  --base <branch>        Base branch (default: main)"
    echo "  --draft                Create as draft PR"
    echo "  --force                Continue even with uncommitted changes (not recommended)"
    echo "  --help                 Show this help message"
    echo ""
    echo "Example:"
    echo "  $0"
    echo "  $0 --title \"feat: Add new feature\" --body \"Description\""
    echo "  $0 --draft"
    exit 1
}

# Parse arguments
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --title)
                TITLE="$2"
                shift 2
                ;;
            --body)
                BODY="$2"
                shift 2
                ;;
            --base)
                BASE_BRANCH="$2"
                shift 2
                ;;
            --draft)
                DRAFT_FLAG="--draft"
                shift
                ;;
            --force)
                FORCE=true
                shift
                ;;
            --help)
                usage
                ;;
            *)
                log_error "Unknown option: $1"
                usage
                ;;
        esac
    done
}

# Validate environment
validate_environment() {
    log_step "Validating environment..."

    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        log_error "Not a git repository"
        exit 1
    fi

    # Check if gh CLI is installed
    if ! command -v gh &> /dev/null; then
        log_error "gh CLI is not installed"
        log_info "Install from: https://cli.github.com/"
        exit 1
    fi

    # Check if we're in a worktree directory (optional warning)
    GIT_DIR=$(git rev-parse --git-dir)
    if [[ ! "$GIT_DIR" =~ \.git/worktrees/ ]]; then
        log_warn "Not in a worktree directory (running from main repository)"
    fi

    # Get current branch
    CURRENT_BRANCH=$(git branch --show-current)
    if [ -z "$CURRENT_BRANCH" ]; then
        log_error "Cannot determine current branch"
        exit 1
    fi

    # Get repository root
    REPO_ROOT=$(git rev-parse --show-toplevel)

    log_info "Current branch: $CURRENT_BRANCH"
    log_info "Base branch: $BASE_BRANCH"
}

# Check for uncommitted changes
check_uncommitted_changes() {
    log_step "Checking for uncommitted changes..."

    if [ "$(git status --porcelain)" != "" ]; then
        if [ "$FORCE" = false ]; then
            log_error "You have uncommitted changes"
            echo ""
            git status --short
            echo ""
            log_info "Please commit or stash your changes before creating PR"
            log_info "Or use --force to skip this check (not recommended)"
            exit 1
        else
            log_warn "Uncommitted changes detected (--force specified, continuing anyway)"
        fi
    else
        log_info "No uncommitted changes detected"
    fi
}

# Check unpushed commits and push if needed
check_unpushed_commits() {
    log_step "Checking push status..."

    # Check if remote branch exists
    if ! git rev-parse --verify --quiet "origin/$CURRENT_BRANCH" > /dev/null 2>&1; then
        log_warn "Remote branch does not exist"
        log_info "Pushing branch to origin..."
        git push -u origin "$CURRENT_BRANCH"
        return
    fi

    # Check for unpushed commits
    UNPUSHED=$(git rev-list "origin/$CURRENT_BRANCH..HEAD" --count 2>/dev/null || echo "0")
    if [ "$UNPUSHED" -gt 0 ]; then
        log_warn "$UNPUSHED unpushed commit(s) detected"
        log_info "Pushing to origin..."
        git push origin "$CURRENT_BRANCH"
    else
        log_info "Branch is up to date with origin"
    fi
}

# Analyze changes and output structured information
analyze_changes() {
    log_step "Analyzing changes..."
    echo ""

    # Get changed files
    local changed_files
    changed_files=$(git diff --name-only "$BASE_BRANCH"...HEAD 2>/dev/null || git diff --name-only "$BASE_BRANCH..HEAD" 2>/dev/null || echo "")

    if [ -z "$changed_files" ]; then
        log_warn "No changes detected compared to $BASE_BRANCH"
        return
    fi

    log_section "### 変更ファイル分析"
    echo ""

    # Backend layers
    local domain_files=$(echo "$changed_files" | grep "app/domain/" || true)
    local app_files=$(echo "$changed_files" | grep "app/application/" || true)
    local infra_files=$(echo "$changed_files" | grep "app/infrastructure/" || true)
    local api_files=$(echo "$changed_files" | grep "app/api/" || true)
    local schema_files=$(echo "$changed_files" | grep "app/schemas/" || true)

    # Frontend
    local web_files=$(echo "$changed_files" | grep "apps/web/" || true)

    # Tests
    local test_files=$(echo "$changed_files" | grep "tests/" || true)

    # Backend section
    if [[ -n "$domain_files" || -n "$app_files" || -n "$infra_files" || -n "$api_files" || -n "$schema_files" ]]; then
        echo "**Backend:**"
        if [ -n "$domain_files" ]; then
            local count=$(echo "$domain_files" | wc -l | tr -d ' ')
            local names=$(echo "$domain_files" | xargs -I {} basename {} | head -3 | tr '\n' ', ' | sed 's/,$//')
            echo "- Domain層: ${count}ファイル (${names})"
        fi
        if [ -n "$app_files" ]; then
            local count=$(echo "$app_files" | wc -l | tr -d ' ')
            local names=$(echo "$app_files" | xargs -I {} basename {} | head -3 | tr '\n' ', ' | sed 's/,$//')
            echo "- Application層: ${count}ファイル (${names})"
        fi
        if [ -n "$infra_files" ]; then
            local count=$(echo "$infra_files" | wc -l | tr -d ' ')
            local names=$(echo "$infra_files" | xargs -I {} basename {} | head -3 | tr '\n' ', ' | sed 's/,$//')
            echo "- Infrastructure層: ${count}ファイル (${names})"
        fi
        if [ -n "$api_files" ]; then
            local count=$(echo "$api_files" | wc -l | tr -d ' ')
            local names=$(echo "$api_files" | xargs -I {} basename {} | head -3 | tr '\n' ', ' | sed 's/,$//')
            echo "- API層: ${count}ファイル (${names})"
        fi
        if [ -n "$schema_files" ]; then
            local count=$(echo "$schema_files" | wc -l | tr -d ' ')
            local names=$(echo "$schema_files" | xargs -I {} basename {} | head -3 | tr '\n' ', ' | sed 's/,$//')
            echo "- Schemas: ${count}ファイル (${names})"
        fi
        echo ""
    fi

    # Frontend section
    if [ -n "$web_files" ]; then
        local count=$(echo "$web_files" | wc -l | tr -d ' ')
        echo "**Frontend:**"
        echo "- ${count}ファイル変更"
        echo ""
    fi

    # Tests section
    if [ -n "$test_files" ]; then
        local count=$(echo "$test_files" | wc -l | tr -d ' ')
        echo "**テスト:**"
        echo "- ${count}ファイル"
        echo ""
    fi

    # Statistics
    local stats=$(git diff --shortstat "$BASE_BRANCH"...HEAD 2>/dev/null || git diff --shortstat "$BASE_BRANCH..HEAD" 2>/dev/null || echo "")
    if [ -n "$stats" ]; then
        echo "**統計:**"
        echo "- $stats"
        echo ""
    fi

    # Commit history
    local commits=$(git log --oneline "$BASE_BRANCH..HEAD" 2>/dev/null | head -10)
    if [ -n "$commits" ]; then
        echo "**コミット履歴:**"
        echo "$commits" | while read -r line; do
            echo "- $line"
        done
        echo ""
    fi

    # TASK.md extraction
    extract_task_info
}

# Extract information from TASK.md
extract_task_info() {
    local task_file="$REPO_ROOT/TASK.md"

    if [ ! -f "$task_file" ]; then
        # Try worktree root
        task_file="$(pwd)/TASK.md"
    fi

    if [ ! -f "$task_file" ]; then
        # TASK.md not found, skip extraction
        return
    fi

    log_section "### TASK.md情報"
    echo ""

    # Extract summary section (概要) - macOS compatible
    local summary=$(awk '/^## 概要/{flag=1; next} /^## /{flag=0} flag' "$task_file" 2>/dev/null | head -5)
    if [ -n "$summary" ]; then
        echo "**Summary候補:**"
        echo "$summary"
        echo ""
    fi

    # Extract acceptance criteria (受け入れ基準) - macOS compatible
    local criteria=$(awk '/^## 受け入れ基準/{flag=1; next} /^## /{flag=0} flag' "$task_file" 2>/dev/null)
    if [ -n "$criteria" ]; then
        echo "**Test plan候補:**"
        echo "$criteria"
        echo ""
    fi
}

# Create pull request
create_pull_request() {
    log_step "Creating pull request..."

    # Build gh pr create command
    GH_CMD="gh pr create --base $BASE_BRANCH"

    if [ -n "$DRAFT_FLAG" ]; then
        GH_CMD="$GH_CMD $DRAFT_FLAG"
    fi

    if [ -n "$TITLE" ]; then
        GH_CMD="$GH_CMD --title \"$TITLE\""
    fi

    if [ -n "$BODY" ]; then
        GH_CMD="$GH_CMD --body \"$BODY\""
    fi

    # Execute PR creation
    if eval "$GH_CMD"; then
        # Get PR URL
        PR_URL=$(gh pr view --json url -q .url 2>/dev/null || echo "")

        log_info "Pull request created successfully"
        if [ -n "$PR_URL" ]; then
            log_info "PR URL: $PR_URL"
        fi
        return 0
    else
        log_error "Failed to create pull request"
        exit 1
    fi
}

# Print summary
print_summary() {
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN} PR Created!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""

    if [ -n "$PR_URL" ]; then
        echo "PR URL: $PR_URL"
    fi
    echo ""
    echo "Note:"
    echo "  - Worktree is NOT deleted"
    echo "  - You can continue making changes locally"
    echo ""
    echo "After PR is merged, run:"
    echo "  /cleanup-worktree"
    echo ""
}

# Main
main() {
    parse_arguments "$@"
    validate_environment
    check_uncommitted_changes
    check_unpushed_commits
    analyze_changes
    create_pull_request
    print_summary
}

main "$@"
