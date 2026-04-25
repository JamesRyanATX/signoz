---
description: Audit both custom branches at a versioned release. Run as /audit-fork <version> (e.g. /audit-fork v0.98.0).
argument-hint: <version>
allowed-tools: Bash Read
---

Audit the fork at version `$ARGUMENTS`. Run **every check before reporting** — do not stop at first failure. Collect all results, then emit the final report at the end.

## Argument parsing

Normalize the version:
- Strip leading `v` if absent, re-add it (so `0.98.0` and `v0.98.0` both → `v0.98.0`)
- Replace `.` with `-` and prepend branch prefixes:
  - `VERSIONED_MDB=multiple-database-names-v0-98-0`
  - `VERSIONED_LDAP=ldap-v0-98-0`
- Path suffix (for worktree dirs): e.g. `v0-98-0`

Store as `VERSION`, `VERSIONED_MDB`, `VERSIONED_LDAP`, and:
- `WORKTREE_MDB=/tmp/signoz-audit-mdn-<path-suffix>`
- `WORKTREE_LDAP=/tmp/signoz-audit-ldap-<path-suffix>`

---

## Step 1 — Create worktrees

```bash
git worktree add "$WORKTREE_MDB" "$VERSIONED_MDB"
```

**FAIL (stop):** branch not found. List available versioned branches and stop — do not run further checks:
```bash
git branch --list "multiple-database-names-*"
```

```bash
git worktree add "$WORKTREE_LDAP" "$VERSIONED_LDAP"
```

**FAIL (stop):** branch not found. List available ldap versioned branches:
```bash
git branch --list "ldap-*"
```

The cleanup commands run at the end **regardless of results**:
```bash
git worktree remove --force "$WORKTREE_MDB"
git worktree remove --force "$WORKTREE_LDAP"
```

---

## Part 1 — Audit `multiple-database-names-vX`

All checks in this part use `$WORKTREE_MDB`.

### Checks A1–A6: DBName() in all 6 telemetry packages

For each package (`telemetrytraces`, `telemetrylogs`, `telemetrymetrics`, `telemetrymeter`, `telemetrymetadata`) run:
```bash
grep -n "^func DBName()" "$WORKTREE_MDB/pkg/<pkg>/tables.go"
```

For analytics:
```bash
grep -n "^func DBName()" "$WORKTREE_MDB/pkg/analytics/tables.go"
```

**PASS:** function declaration found. **FAIL:** missing — report the package.

### Checks B1–B6: Init() in all 6 packages

For each package (`telemetrytraces`, `telemetrylogs`, `telemetrymetrics`, `telemetrymeter`, `telemetrymetadata`) run:
```bash
grep -n "^func Init(databaseName string)" "$WORKTREE_MDB/pkg/<pkg>/tables.go"
```

For analytics:
```bash
grep -n "^func Init(databaseName string)" "$WORKTREE_MDB/pkg/analytics/tables.go"
```

**PASS:** function declaration found. **FAIL:** missing.

### Check C: telemetrystore config fields

```bash
grep -n "TraceDatabase\|MetricsDatabase\|LogsDatabase\|MeterDatabase\|MetadataDatabase\|AnalyticsDatabase" \
  "$WORKTREE_MDB/pkg/telemetrystore/config.go"
```

**PASS:** all 6 field names found. **FAIL:** report which are missing.

### Check D: default constants

```bash
grep -n "DefaultTraceDatabase\|DefaultMetricsDatabase\|DefaultLogsDatabase\|DefaultMeterDatabase\|DefaultMetadataDatabase\|DefaultAnalyticsDatabase" \
  "$WORKTREE_MDB/pkg/telemetrystore/defaults.go"
```

**PASS:** all 6 constant names found. **FAIL:** report which are missing.

### Check E: provider Init() wiring

```bash
grep -n "telemetrytraces\.Init\|telemetrymetrics\.Init\|telemetrylogs\.Init\|telemetrymeter\.Init\|telemetrymetadata\.Init\|analytics\.Init" \
  "$WORKTREE_MDB/pkg/telemetrystore/clickhousetelemetrystore/provider.go"
```

**PASS:** all 6 Init() calls found (one per line). **FAIL:** report which are missing.

### Check F: env-var backward compat

```bash
grep -n "CLICKHOUSE_TRACE_DATABASE\|CLICKHOUSE_DATABASE\|CLICKHOUSE_LOG_DATABASE\|CLICKHOUSE_METER_DATABASE\|CLICKHOUSE_ANALYTICS_DATABASE" \
  "$WORKTREE_MDB/pkg/signoz/config.go"
```

**PASS:** all 5 env-var names found. **FAIL:** report which are missing. Note: there is no dedicated `CLICKHOUSE_METADATA_DATABASE` env var — this is expected and is not a failure.

### Check G: no rogue hardcoded database strings in SQL paths

```bash
grep -rn \
  '"signoz_traces\|"signoz_logs\|"signoz_metrics\|"signoz_meter\|"signoz_metadata\|"signoz_analytics' \
  "$WORKTREE_MDB/pkg/" --include="*.go" | grep -v "_test.go"
```

For each hit, classify it:
- **Acceptable:** inside a `const` or `var` block used only as a default config value (e.g. in `defaults.go`, `createDefaultConfig`, or a `Default*` constant that feeds config initialization), or in a comment.
- **Flag:** if the const/var it assigns is referenced in SQL execution paths — `PrepareBatch`, `fmt.Sprintf` building SQL, or direct query execution.

**PASS:** no flagged hits (only acceptable defaults). **FAIL:** report file:line and the SQL usage location.

### Check H: Go build (mdn)

```bash
cd "$WORKTREE_MDB" && go build ./... 2>&1
```

**PASS:** exits 0. **FAIL:** compiler errors — report them. If errors are in files touched by our commit, flag as a rebase regression.

### Check I: Go tests (mdn)

```bash
cd "$WORKTREE_MDB" && go test -count=1 \
  ./pkg/telemetrystore/... \
  ./pkg/telemetrytraces/... \
  ./pkg/telemetrylogs/... \
  ./pkg/telemetrymetrics/... \
  ./pkg/telemetrymeter/... \
  ./pkg/telemetrymetadata/... \
  ./pkg/analytics/... 2>&1
```

Allow up to 5 minutes.

**PASS:** exits 0, all packages print `ok`.
**SKIP:** exits non-zero but all failures are `connection refused`, `context deadline exceeded`, or `dial tcp ... connect: connection refused` — ClickHouse not available.
**FAIL:** any `FAIL` line not attributable to missing ClickHouse.

---

## Part 2 — Audit `ldap-vX`

All checks in this part use `$WORKTREE_LDAP`.

### Check J: go-ldap dependency

```bash
grep "go-ldap" "$WORKTREE_LDAP/go.mod"
```

**PASS:** `github.com/go-ldap/ldap` found. **FAIL:** dependency missing.

### Check K: LdapConfig type

```bash
grep -n "^type LdapConfig struct" "$WORKTREE_LDAP/pkg/types/ssotypes/ldap.go"
```

**PASS:** struct declaration found. **FAIL:** file missing or struct not declared.

### Check L: LDAP client entry points

```bash
grep -n "^func NewClient\|^func (c \*Client) Authenticate" \
  "$WORKTREE_LDAP/pkg/ldap/client.go"
```

**PASS:** both `NewClient` and `Authenticate` found. **FAIL:** report which are missing.

### Check M: migration factory registered

```bash
grep -n "NewAddLdapSupportFactory" "$WORKTREE_LDAP/pkg/signoz/provider.go"
```

**PASS:** factory call found. **FAIL:** migration not wired into provider.

### Check N: AuthenticateWithLDAP wired in user module

```bash
grep -n "AuthenticateWithLDAP\|SsoType.*LDAP\|types\.LDAP" \
  "$WORKTREE_LDAP/pkg/modules/user/impluser/module.go"
```

Verify two things are present: the method definition (`func.*AuthenticateWithLDAP`) and the LDAP-before-standard-login branch (`SsoEnabled.*SsoType.*LDAP` or equivalent conditional).

**PASS:** both found. **FAIL:** report which is missing.

### Check O: SSO enabled in plans

```bash
awk '/BasicPlan/,/^}/' "$WORKTREE_LDAP/pkg/types/licensetypes/plan.go" \
  | grep -A3 "SSO" | grep "Active"
```

```bash
awk '/DefaultFeatureSet/,/^}/' "$WORKTREE_LDAP/pkg/types/licensetypes/plan.go" \
  | grep -A3 "SSO" | grep "Active"
```

**PASS:** `Active: true` found for SSO in both `BasicPlan` and `DefaultFeatureSet`. **FAIL:** SSO missing or `Active: false` in either.

### Check P: frontend EditLDAP component

```bash
test -f "$WORKTREE_LDAP/frontend/src/container/OrganizationSettings/AuthDomains/Edit/EditLDAP.tsx" \
  && echo "FOUND" || echo "MISSING"
```

**PASS:** file exists. **FAIL:** missing.

### Check Q: Go build (full, includes ldap)

```bash
cd "$WORKTREE_LDAP" && go build ./... 2>&1
```

**PASS:** exits 0. **FAIL:** compiler errors. If errors are in `pkg/ldap/`, `pkg/modules/user/`, or `pkg/types/`, flag as a rebase regression.

### Check R: Go tests for ldap package

```bash
cd "$WORKTREE_LDAP" && go test -count=1 ./pkg/ldap/... 2>&1
```

**PASS:** exits 0. **FAIL:** test failures — report test names. Note: `pkg/ldap/client_test.go` tests require an LDAP server; if failures are all `connection refused` or dial errors, mark **SKIP** rather than FAIL.

---

## Step 2 — Cleanup worktrees

Always run, even if earlier steps failed:

```bash
git worktree remove --force "$WORKTREE_MDB"
git worktree remove --force "$WORKTREE_LDAP"
```

---

## Step 3 — Final report

Print the structured report using the exact format below.

```
AUDIT REPORT — <VERSION>
================================================

MULTIPLE-DATABASE-NAMES  (<VERSIONED_MDB>)
─────────────────────────────────────────────────  ──────
A1 telemetrytraces    DBName()                     PASS/FAIL
A2 telemetrylogs      DBName()                     PASS/FAIL
A3 telemetrymetrics   DBName()                     PASS/FAIL
A4 telemetrymeter     DBName()                     PASS/FAIL
A5 telemetrymetadata  DBName()                     PASS/FAIL
A6 analytics          DBName()                     PASS/FAIL
B1 telemetrytraces    Init()                       PASS/FAIL
B2 telemetrylogs      Init()                       PASS/FAIL
B3 telemetrymetrics   Init()                       PASS/FAIL
B4 telemetrymeter     Init()                       PASS/FAIL
B5 telemetrymetadata  Init()                       PASS/FAIL
B6 analytics          Init()                       PASS/FAIL
C  telemetrystore: 6 config fields                 PASS/FAIL
D  defaults.go: 6 default constants                PASS/FAIL
E  provider.go: 6 Init() calls wired               PASS/FAIL
F  signoz/config.go: env-var compat (5 vars)       PASS/FAIL
G  No rogue hardcoded database strings             PASS/FAIL
H  go build ./...                                  PASS/FAIL
I  go test ./pkg/telemetry*/ ./pkg/analytics/      PASS/FAIL/SKIP

LDAP  (<VERSIONED_LDAP>)
─────────────────────────────────────────────────  ──────
J  go.mod: go-ldap dependency                      PASS/FAIL
K  ssotypes/ldap.go: LdapConfig struct             PASS/FAIL
L  pkg/ldap/client.go: NewClient + Authenticate    PASS/FAIL
M  provider.go: migration factory registered       PASS/FAIL
N  user module: AuthenticateWithLDAP wired         PASS/FAIL
O  licensetypes: SSO active in BasicPlan+Default   PASS/FAIL
P  frontend: EditLDAP.tsx exists                   PASS/FAIL
Q  go build ./... (full)                           PASS/FAIL
R  go test ./pkg/ldap/...                          PASS/FAIL/SKIP
─────────────────────────────────────────────────  ──────
OVERALL                                            PASS/FAIL

FAILURES
────────
<For each FAIL: check code — file:line — finding — recommended fix>
```

**OVERALL is PASS** if all checks pass or only SKIP. **OVERALL is FAIL** if any check fails.
