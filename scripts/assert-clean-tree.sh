#!/bin/sh
# Fail when the working tree moved, including files the command created but
# never added to the index. `git diff --exit-code` only looks at tracked paths,
# so a generator that starts writing a new document would slip past it.
set -eu

what="${1:-the previous step}"

# --porcelain lists modified tracked paths and untracked paths alike.
status="$(git status --porcelain --untracked-files=all)"
if [ -n "$status" ]; then
	echo "$what changed the working tree; commit the result." >&2
	echo "$status" >&2
	git --no-pager diff >&2
	exit 1
fi
