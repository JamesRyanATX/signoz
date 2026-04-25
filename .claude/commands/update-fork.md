---
description: Safely rebase both custom branches to a new upstream SigNoz version. Run as /update-fork <version> (e.g. /update-fork v0.98.0).
argument-hint: <version>
allowed-tools: Bash Read Grep
---

Safely update this fork to upstream SigNoz version `$ARGUMENTS`.

This repo is a personal fork of `SigNoz/signoz`. It carries two custom branches:

- **`multiple-database-names`** — exactly one commit on top of `main` that makes all 6 ClickHouse database names runtime-configurable via env vars and config fields.
- **`ldap`** — 7 commits on top of `multiple-database-names` that add a proprietary LDAP authentication overlay.

This command rebases both onto a new upstream version tag, producing `multiple-database-names-vX` and `ldap-vX`. The `utility` branch (home of CLAUDE.md and these commands) is not touched.

## Files our `multiple-database-names` commit touches

These are the key files to watch for upstream conflicts and for the impact assessment:

```
pkg/analytics/tables.go
pkg/telemetrytraces/tables.go
pkg/telemetrylogs/tables.go
pkg/telemetrymetrics/tables.go
pkg/telemetrymeter/tables.go
pkg/telemetrymetadata/tables.go
pkg/telemetrystore/config.go
pkg/telemetrystore/defaults.go
pkg/telemetrystore/clickhousetelemetrystore/provider.go
pkg/signoz/config.go
pkg/query-service/constants/constants.go
pkg/query-service/app/clickhouseReader/reader.go
pkg/query-service/app/clickhouseReader/options.go
pkg/querybuilder/resourcefilter/statement_builder.go
pkg/telemetrytraces/statement_builder.go
pkg/telemetrylogs/statement_builder.go
pkg/telemetrymetrics/statement_builder.go
pkg/telemetrymeter/statement_builder.go
```

## Files our `ldap` commits touch

```
pkg/ldap/client.go
pkg/ldap/client_test.go
pkg/types/ssotypes/ldap.go
pkg/types/domain.go
pkg/types/licensetypes/plan.go
pkg/modules/user/impluser/module.go
pkg/modules/user/impluser/handler.go
pkg/modules/user/user.go
pkg/sqlmigration/050_add_ldap_support.go
pkg/signoz/provider.go
go.mod
go.sum
frontend/src/container/OrganizationSettings/AuthDomains/Edit/EditLDAP.tsx
frontend/src/types/api/SAML/listDomain.ts
```

## Step 1 — Normalize version and derive names

Parse `$ARGUMENTS`:
- Strip leading `v` if present, then re-add it to normalize (e.g. both `0.98.0` and `v0.98.0` → `v0.98.0`)
- Replace dots with dashes and prepend branch prefixes:
  - `VERSIONED_MDB=multiple-database-names-v0-98-0`
  - `VERSIONED_LDAP=ldap-v0-98-0`

Store `VERSION`, `VERSIONED_MDB`, and `VERSIONED_LDAP` for the rest of the steps.

## Step 2 — Pre-flight checks

Run all of these before touching anything:

```bash
# Must be clean
git status --porcelain
```
Stop if there are any uncommitted changes. Ask the user to stash or commit first.

```bash
# Neither versioned branch must already exist
git branch --list "$VERSIONED_MDB"
git branch --list "$VERSIONED_LDAP"
```
Stop if either exists — the update for this version may have already been done, or there's a naming collision.

```bash
# Confirm we're in the right repo (signoz, not signoz-otel-collector)
git remote get-url origin
```
Should contain `signoz` but NOT `signoz-otel-collector`. Stop if wrong.

## Step 3 — Add upstream remote if missing

```bash
git remote get-url upstream 2>/dev/null || echo "MISSING"
```

If missing:
```bash
git remote add upstream git@github.com:SigNoz/signoz.git
```

## Step 4 — Fetch upstream and verify the tag

```bash
git fetch upstream --tags
```

```bash
git tag --list "$VERSION"
```

Stop with a clear error if the tag doesn't exist. List the 5 most recent tags to help the user find the right name:
```bash
git tag --sort=-version:refname | head -5
```

## Step 5 — Check if main needs updating

```bash
CURRENT_MAIN=$(git rev-parse main)
TARGET_SHA=$(git rev-parse "$VERSION")
echo "main is at: $CURRENT_MAIN"
echo "target tag is at: $TARGET_SHA"
git merge-base --is-ancestor "$TARGET_SHA" main && echo "ALREADY_AT_OR_AHEAD" || echo "NEEDS_UPDATE"
```

If main is already at or ahead of the target, skip steps 6–7 and go straight to step 8.

## Step 6 — Update main

Show the user what is about to happen:
```bash
echo "Resetting main: $(git rev-parse --short main) → $(git rev-parse --short $VERSION)"
echo "Commits being pulled in:"
git log --oneline main..$VERSION
```

**Pause here and confirm with the user before proceeding.** This is a force operation.

```bash
git checkout main
git reset --hard "$VERSION"
```

## Step 7 — Force push main

**Confirm with the user before running this.** This rewrites `origin/main`.

```bash
git push origin main --force
```

## Step 8 — Create `multiple-database-names-vX` and rebase

```bash
git checkout -b "$VERSIONED_MDB" multiple-database-names
```

```bash
git rebase main
```

### If rebase conflicts occur in the mdn commit

For each conflicted file, read it carefully and resolve by preserving these invariants:

- **`pkg/telemetry*/tables.go` and `pkg/analytics/tables.go`**: Each file must retain the `var dbName string`, `func Init(databaseName string)` (sets dbName if non-empty), and `func DBName() string` (returns dbName, falling back to the hardcoded default). If upstream added or reorganized constants in these files, keep our three functions intact.
- **`pkg/telemetrystore/config.go`**: Must retain the 6 `*Database` fields (`TraceDatabase`, `MetricsDatabase`, `LogsDatabase`, `MeterDatabase`, `MetadataDatabase`, `AnalyticsDatabase`) with their `mapstructure` tags in the `ClickhouseConfig` struct.
- **`pkg/telemetrystore/defaults.go`**: Must retain all 6 `Default*Database` constants.
- **`pkg/telemetrystore/clickhousetelemetrystore/provider.go`**: Must retain all 6 `pkg.Init(config.Clickhouse.*Database)` calls near the top of the provider constructor.
- **`pkg/signoz/config.go`**: Must retain the backward-compat env-var block that reads `CLICKHOUSE_TRACE_DATABASE`, `CLICKHOUSE_DATABASE`, `CLICKHOUSE_LOG_DATABASE`, `CLICKHOUSE_METER_DATABASE`, and `CLICKHOUSE_ANALYTICS_DATABASE`.
- **Statement builder files (`pkg/telemetry*/statement_builder.go`, etc.)**: These call `DBName()` instead of using a hardcoded string. If upstream refactored these files, ensure the calls still use `DBName()`.

If upstream added new SQL queries with hardcoded database name strings in any of the watched files, add the same `DBName()` treatment our commit applies to the existing queries.

Stage and continue for each resolved file:
```bash
git add <file> && git rebase --continue
```

If the rebase becomes too complex, run `git rebase --abort`, report exactly which files conflicted and why, and ask the user how to proceed.

## Step 9 — Create `ldap-vX` rebased onto `multiple-database-names-vX`

```bash
git checkout -b "$VERSIONED_LDAP" ldap
git rebase --onto "$VERSIONED_MDB" multiple-database-names
```

This replays the 7 LDAP commits (`multiple-database-names..ldap`) onto `$VERSIONED_MDB` as the new base.

### If rebase conflicts occur in the ldap commits

The core invariants to preserve across all 7 commits:

- **`pkg/modules/user/impluser/module.go`**: Must retain the `AuthenticateWithLDAP` method and the LDAP-check-before-standard-login block. The flow is: if domain has `SsoEnabled && SsoType == types.LDAP`, call `AuthenticateWithLDAP` and return early; otherwise fall through to the normal password check.
- **`pkg/types/licensetypes/plan.go`**: Must retain `SSO` feature as `Active: true` in both `BasicPlan` and `DefaultFeatureSet`. If upstream changed the plan definitions, preserve our SSO activation.
- **`pkg/types/domain.go`**: Must retain any LDAP-specific fields we added. If upstream changed the domain struct, merge carefully.
- **`pkg/signoz/provider.go`**: Must retain `sqlmigration.NewAddLdapSupportFactory(sqlstore)` in the migration registration list.

Stage and continue:
```bash
git add <file> && git rebase --continue
```

If the rebase becomes too complex, abort, report the conflicted files with context, and ask the user how to proceed.

## Step 10 — Assess upstream impact

After both rebases succeed, run this analysis to determine if our features need follow-up work.

Capture the old main SHA from step 5 as `OLD_MAIN_SHA`.

**For `multiple-database-names-vX`:**

```bash
# Did upstream change any of our watched files?
git log --oneline "$OLD_MAIN_SHA".."$VERSION" -- \
  pkg/analytics/ \
  pkg/telemetrytraces/ \
  pkg/telemetrylogs/ \
  pkg/telemetrymetrics/ \
  pkg/telemetrymeter/ \
  pkg/telemetrymetadata/ \
  pkg/telemetrystore/ \
  pkg/signoz/config.go \
  pkg/query-service/constants/constants.go \
  pkg/query-service/app/clickhouseReader/ \
  pkg/querybuilder/
```

```bash
# Did upstream introduce new hardcoded database name strings?
git diff "$OLD_MAIN_SHA".."$VERSION" -- \
  pkg/telemetry*/ pkg/analytics/ pkg/query-service/ pkg/querybuilder/ \
  | grep -E '^\+.*\b(signoz_logs|signoz_traces|signoz_metrics|signoz_metadata|signoz_analytics|signoz_meter)\b' \
  | grep -v '^+++' | grep -v '^\+\s*//' | grep -v 'DefaultDatabase\|default\|// '
```

**For `ldap-vX`:**

```bash
# Did upstream change any of our LDAP watched files?
git log --oneline "$OLD_MAIN_SHA".."$VERSION" -- \
  pkg/modules/user/ \
  pkg/types/domain.go \
  pkg/types/licensetypes/ \
  pkg/sqlmigration/ \
  pkg/signoz/provider.go
```

Report to the user:
- Which watched files upstream also changed
- Whether any new hardcoded database name strings were introduced (flag as requiring extension of our feature)
- Overall verdict: **no action needed** vs **specific follow-up required**

## Step 11 — Build check

Switch to the ldap-vX branch (which inherits all mdn changes) and build:

```bash
git checkout "$VERSIONED_LDAP"
go build ./...
```

Fix any compilation errors before proceeding. If errors are in our custom code, resolve them. If errors are in unrelated upstream code, report them separately.

## Step 12 — Push both branches and hand off

```bash
git push origin "$VERSIONED_MDB"
git push origin "$VERSIONED_LDAP"
```

Then print a summary and stop. Do **not** touch `multiple-database-names`, `ldap`, or `utility`.

Summary to print:
- `$VERSIONED_MDB` ready (pushed to origin)
- `$VERSIONED_LDAP` ready (pushed to origin), based on `$VERSIONED_MDB`
- `main` updated: old SHA → `$(git rev-parse --short $VERSION)`
- Number of upstream commits pulled in
- Any follow-up items from the impact assessment (step 10)
- Verification commands:

```bash
git log --oneline "$VERSIONED_MDB" -5
git log --oneline "$VERSIONED_LDAP" -5
go test -race \
  ./pkg/telemetrystore/... \
  ./pkg/telemetrytraces/... \
  ./pkg/telemetrylogs/... \
  ./pkg/telemetrymetrics/... \
  ./pkg/telemetrymeter/... \
  ./pkg/telemetrymetadata/... \
  ./pkg/analytics/... \
  ./pkg/ldap/...
```
