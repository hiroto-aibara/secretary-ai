#!/bin/bash
set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

# Default values
FORCE=false
DELETE_BRANCH=false
RETURN_BRANCH="main"

# Variables set during validation
CURRENT_BRANCH=""
WORKTREE_PATH=""
WORKTREE_NAME=""
REPO_ROOT=""

# Usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Removes the current worktree and returns to the specified branch."
    echo ""
    echo "Options:"
    echo "  --branch <branch>      Branch to return to (default: main)"
    echo "  --force                Force remove even with uncommitted changes"
    echo "  --delete-branch        Also delete the local branch"
    echo "  --help                 Show this help message"
    echo ""
    echo "Example:"
    echo "  $0"
    echo "  $0 --branch develop"
    echo "  $0 --delete-branch"
    exit 1
}

# Parse arguments
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --branch)
                RETURN_BRANCH="$2"
                shift 2
                ;;
            --force)
                FORCE=true
                shift
                ;;
            --delete-branch)
                DELETE_BRANCH=true
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

    # Check if we're in a worktree directory
    GIT_DIR=$(git rev-parse --git-dir)
    if [[ ! "$GIT_DIR" =~ \.git/worktrees/ ]]; then
        log_error "Not in a worktree directory"
        log_info "This script must be run from within .worktrees/<feature-name>/"
        exit 1
    fi

    # Get current branch
    CURRENT_BRANCH=$(git branch --show-current)
    if [ -z "$CURRENT_BRANCH" ]; then
        log_error "Cannot determine current branch"
        exit 1
    fi

    # Get worktree path
    WORKTREE_PATH=$(pwd)
    WORKTREE_NAME=$(basename "$WORKTREE_PATH")

    # Get repository root (parent of .worktrees)
    # Navigate up to find the main repository
    REPO_ROOT=$(git worktree list | head -1 | awk '{print $1}')

    log_info "Current branch: $CURRENT_BRANCH"
    log_info "Worktree path: $WORKTREE_PATH"
    log_info "Repository root: $REPO_ROOT"
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
            log_info "Please commit or stash your changes before cleanup"
            log_info "Or use --force to skip this check (changes will be lost!)"
            exit 1
        else
            log_warn "Uncommitted changes detected (--force specified, changes will be lost!)"
        fi
    else
        log_info "No uncommitted changes detected"
    fi
}

# Cleanup worktree
cleanup_worktree() {
    log_step "Removing worktree..."

    # Move to repository root
    cd "$REPO_ROOT"

    # Remove worktree
    if [ "$FORCE" = true ]; then
        if git worktree remove --force "$WORKTREE_PATH"; then
            log_info "Worktree force-removed: $WORKTREE_PATH"
        else
            log_error "Failed to remove worktree"
            log_info ""
            log_info "You can manually remove it with:"
            log_info "  git worktree remove --force $WORKTREE_PATH"
            exit 1
        fi
    else
        if git worktree remove "$WORKTREE_PATH"; then
            log_info "Worktree removed: $WORKTREE_PATH"
        else
            log_error "Failed to remove worktree"
            log_info ""
            log_info "You can manually remove it with:"
            log_info "  git worktree remove $WORKTREE_PATH"
            log_info ""
            log_info "Or force remove:"
            log_info "  git worktree remove --force $WORKTREE_PATH"
            exit 1
        fi
    fi
}

# Delete local branch
delete_local_branch() {
    if [ "$DELETE_BRANCH" = false ]; then
        return 0
    fi

    log_step "Deleting local branch..."

    if git branch -d "$CURRENT_BRANCH" 2>/dev/null; then
        log_info "Local branch deleted: $CURRENT_BRANCH"
    else
        log_warn "Could not delete branch (may not be fully merged)"
        log_info "To force delete: git branch -D $CURRENT_BRANCH"
    fi
}

# Return to main branch
return_to_branch() {
    log_step "Returning to $RETURN_BRANCH branch..."

    if git checkout "$RETURN_BRANCH"; then
        log_info "Now on branch: $RETURN_BRANCH"
    else
        log_warn "Failed to checkout $RETURN_BRANCH"
        log_info "Current directory: $(pwd)"
    fi
}

# Print summary
print_summary() {
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN} Worktree Cleanup Completed!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo "Removed worktree: $WORKTREE_NAME"
    echo "Current branch: $RETURN_BRANCH"

    if [ "$DELETE_BRANCH" = false ]; then
        echo ""
        echo "Note:"
        echo "  - Local branch ($CURRENT_BRANCH) is preserved"
        echo "  - Remote branch is preserved"
        echo ""
        echo "To delete local branch:"
        echo "  git branch -d $CURRENT_BRANCH"
    fi
    echo ""
}

# Main
main() {
    parse_arguments "$@"
    validate_environment
    check_uncommitted_changes
    cleanup_worktree
    delete_local_branch
    return_to_branch
    print_summary
}

main "$@"
