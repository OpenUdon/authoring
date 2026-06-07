package operationlifecycle

import (
	"path"
	"slices"
	"strings"
	"unicode"

	"github.com/OpenUdon/authoring/promptcontext"
)

type Options struct {
	Goal         string
	DesiredState bool
}

type Expansion struct {
	SeedOperationID string          `json:"seed_operation_id,omitempty"`
	FamilyKey       string          `json:"family_key,omitempty"`
	Roles           []RoleCandidate `json:"roles,omitempty"`
	Diagnostics     []Diagnostic    `json:"diagnostics,omitempty"`
}

type RoleCandidate struct {
	Role       string                           `json:"role"`
	Operation  promptcontext.OperationCandidate `json:"operation"`
	Confidence string                           `json:"confidence,omitempty"`
	Reason     string                           `json:"reason,omitempty"`
}

type Diagnostic struct {
	Code     string `json:"code,omitempty"`
	Severity string `json:"severity,omitempty"`
	Message  string `json:"message,omitempty"`
}

type candidateScore struct {
	operation promptcontext.OperationCandidate
	score     int
	reason    string
}

// Expand returns a conservative lifecycle expansion for seed. It only selects
// sibling operations when one same-source candidate is clearly stronger for a
// role; ambiguous role matches are diagnosed and omitted.
func Expand(ctx promptcontext.Context, seed promptcontext.OperationCandidate, opts Options) Expansion {
	ctx = promptcontext.Normalize(ctx)
	seed = normalizeSeed(ctx, seed)
	seedID := operationID(seed)
	out := Expansion{
		SeedOperationID: seedID,
		FamilyKey:       familyKey(seed),
	}
	if seedID == "" {
		out.Diagnostics = append(out.Diagnostics, Diagnostic{Code: "operation_lifecycle.seed_missing", Severity: "warning", Message: "seed operation is empty"})
		return out
	}
	siblings := map[string]candidateScore{}
	for _, role := range []string{"read", "update", "delete"} {
		if role == "update" && !goalWantsUpdate(opts.Goal) {
			continue
		}
		match, ok, diag := bestSibling(ctx.Operations, seed, role)
		if diag.Code != "" {
			out.Diagnostics = append(out.Diagnostics, diag)
		}
		if !ok {
			continue
		}
		siblings[role] = match
	}
	primaryRole := primaryRole(seed, opts, len(siblings) > 0)
	out.Roles = append(out.Roles, RoleCandidate{
		Role:       primaryRole,
		Operation:  seed,
		Confidence: "high",
		Reason:     "selected seed operation",
	})
	if primaryRole == "read" {
		return out
	}
	for _, role := range []string{"read", "update", "delete"} {
		if role == primaryRole {
			continue
		}
		match, ok := siblings[role]
		if !ok {
			continue
		}
		out.Roles = append(out.Roles, RoleCandidate{
			Role:       role,
			Operation:  match.operation,
			Confidence: confidence(match.score),
			Reason:     match.reason,
		})
	}
	return out
}

func normalizeSeed(ctx promptcontext.Context, seed promptcontext.OperationCandidate) promptcontext.OperationCandidate {
	if operationID(seed) == "" {
		return seed
	}
	for _, op := range ctx.Operations {
		if sameOperation(op, seed) {
			return op
		}
	}
	return seed
}

func primaryRole(seed promptcontext.OperationCandidate, opts Options, expanded bool) string {
	if opts.DesiredState && expanded {
		switch strings.ToUpper(strings.TrimSpace(seed.Verb)) {
		case "POST":
			return "create"
		case "PUT":
			if operationHasAny(seed, "create", "createorupdate", "insert") {
				return "create"
			}
			return "update"
		case "PATCH":
			return "update"
		case "GET", "HEAD":
			return "read"
		case "DELETE":
			return "delete"
		}
	}
	return methodRole(seed)
}

func methodRole(op promptcontext.OperationCandidate) string {
	switch strings.ToUpper(strings.TrimSpace(op.Verb)) {
	case "GET", "HEAD":
		return "read"
	case "DELETE":
		return "delete"
	case "POST":
		return "post"
	case "PUT":
		return "put"
	case "PATCH":
		return "patch"
	default:
		return "post"
	}
}

func bestSibling(ops []promptcontext.OperationCandidate, seed promptcontext.OperationCandidate, role string) (candidateScore, bool, Diagnostic) {
	var matches []candidateScore
	for _, op := range ops {
		if sameOperation(op, seed) || !sameSource(op, seed) {
			continue
		}
		score, reason := siblingScore(op, seed, role)
		if score < 60 {
			continue
		}
		matches = append(matches, candidateScore{operation: op, score: score, reason: reason})
	}
	if len(matches) == 0 {
		return candidateScore{}, false, Diagnostic{}
	}
	slices.SortStableFunc(matches, func(a, b candidateScore) int {
		if a.score != b.score {
			return b.score - a.score
		}
		return strings.Compare(operationID(a.operation), operationID(b.operation))
	})
	if len(matches) > 1 && matches[0].score == matches[1].score {
		return candidateScore{}, false, Diagnostic{
			Code:     "operation_lifecycle.ambiguous_" + role,
			Severity: "warning",
			Message:  "multiple same-source operations match lifecycle role " + role,
		}
	}
	return matches[0], true, Diagnostic{}
}

func siblingScore(op, seed promptcontext.OperationCandidate, role string) (int, string) {
	if !lifecyclePathsMatch(seed.Path, op.Path, role) {
		return 0, ""
	}
	if role == "read" && operationHasAny(op, "list") {
		return 0, ""
	}
	if (role == "update" || role == "delete") && operationHasAny(op, "collection") {
		return 0, ""
	}
	score := 0
	var reasons []string
	if verbMatchesRole(op.Verb, role) {
		score += 45
		reasons = append(reasons, "HTTP verb matches "+role)
	}
	if operationNameMatchesRole(op, role) {
		score += 35
		reasons = append(reasons, "operation id/name matches "+role)
	}
	if role == "update" && strings.EqualFold(op.Verb, "PATCH") {
		score += 8
		reasons = append(reasons, "PATCH update preferred")
	}
	if role == "update" && operationHasAny(op, "patch") {
		score += 6
		reasons = append(reasons, "patch operation preferred")
	}
	if sameFamily(seed, op) {
		score += 30
		reasons = append(reasons, "operation family matches seed")
	}
	score += 25
	reasons = append(reasons, "collection/item path matches seed")
	return score, strings.Join(reasons, "; ")
}

func verbMatchesRole(verb, role string) bool {
	switch role {
	case "read":
		return strings.EqualFold(verb, "GET") || strings.EqualFold(verb, "HEAD")
	case "update":
		return strings.EqualFold(verb, "PUT") || strings.EqualFold(verb, "PATCH")
	case "delete":
		return strings.EqualFold(verb, "DELETE")
	default:
		return false
	}
}

func operationNameMatchesRole(op promptcontext.OperationCandidate, role string) bool {
	switch role {
	case "read":
		return operationHasAny(op, "get", "read", "show", "describe")
	case "update":
		return operationHasAny(op, "update", "patch", "replace", "createorupdate", "put")
	case "delete":
		return operationHasAny(op, "delete", "remove")
	default:
		return false
	}
}

func operationHasAny(op promptcontext.OperationCandidate, terms ...string) bool {
	tokens := operationTokens(op)
	for _, term := range terms {
		if tokens[strings.ToLower(term)] {
			return true
		}
	}
	return false
}

func sameSource(a, b promptcontext.OperationCandidate) bool {
	return strings.TrimSpace(a.SourceID) != "" && strings.TrimSpace(a.SourceID) == strings.TrimSpace(b.SourceID)
}

func sameOperation(a, b promptcontext.OperationCandidate) bool {
	aID := operationID(a)
	bID := operationID(b)
	return aID != "" && aID == bID && strings.TrimSpace(a.SourceID) == strings.TrimSpace(b.SourceID)
}

func operationID(op promptcontext.OperationCandidate) string {
	return strings.TrimSpace(firstNonEmpty(op.OperationID, op.ID))
}

func familyKey(op promptcontext.OperationCandidate) string {
	tokens := familyTokens(op)
	return strings.Join(tokens, ".")
}

func sameFamily(a, b promptcontext.OperationCandidate) bool {
	aTokens := familyTokens(a)
	bTokens := familyTokens(b)
	if len(aTokens) == 0 || len(bTokens) == 0 {
		return false
	}
	common := 0
	bSet := map[string]bool{}
	for _, token := range bTokens {
		bSet[token] = true
	}
	for _, token := range aTokens {
		if bSet[token] {
			common++
		}
	}
	return common >= min(2, min(len(aTokens), len(bTokens)))
}

func familyTokens(op promptcontext.OperationCandidate) []string {
	raw := operationTokens(op)
	var out []string
	for token := range raw {
		if lifecycleWord(token) || token == "v1" || token == "v2" || token == "v3" || token == "api" {
			continue
		}
		out = append(out, token)
	}
	slices.Sort(out)
	return out
}

func operationTokens(op promptcontext.OperationCandidate) map[string]bool {
	text := operationID(op) + " " + op.Name + " " + strings.Join(op.Tags, " ")
	parts := strings.FieldsFunc(text, func(r rune) bool {
		return r == '_' || r == '-' || r == '.' || r == '/' || r == ':' || unicode.IsSpace(r)
	})
	out := map[string]bool{}
	for _, part := range parts {
		for _, token := range splitCamel(part) {
			token = strings.ToLower(strings.TrimSpace(token))
			if token != "" {
				out[token] = true
			}
		}
	}
	joined := strings.ToLower(strings.NewReplacer("_", "", "-", "", ".", "").Replace(operationID(op)))
	if strings.Contains(joined, "createorupdate") {
		out["createorupdate"] = true
		out["create"] = true
		out["update"] = true
	}
	return out
}

func splitCamel(value string) []string {
	var out []string
	start := 0
	runes := []rune(value)
	for i := 1; i < len(runes); i++ {
		if unicode.IsUpper(runes[i]) && (unicode.IsLower(runes[i-1]) || (i+1 < len(runes) && unicode.IsLower(runes[i+1]))) {
			out = append(out, string(runes[start:i]))
			start = i
		}
	}
	out = append(out, string(runes[start:]))
	return out
}

func lifecycleWord(token string) bool {
	switch token {
	case "create", "insert", "post", "put", "get", "read", "show", "describe", "list", "update", "patch", "replace", "delete", "remove", "createorupdate":
		return true
	default:
		return false
	}
}

func lifecyclePathsMatch(seedPath, opPath, role string) bool {
	seed := normalizePath(seedPath)
	op := normalizePath(opPath)
	if seed == "" || op == "" {
		return false
	}
	if seed == op {
		if role == "read" || role == "update" || role == "delete" {
			return hasPathParameter(op)
		}
		return true
	}
	seedBase := collectionPath(seed)
	opBase := collectionPath(op)
	if seedBase == "" || opBase == "" || seedBase != opBase {
		return false
	}
	switch role {
	case "read", "update", "delete":
		return hasPathParameter(op)
	default:
		return true
	}
}

func normalizePath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	value = "/" + strings.Trim(path.Clean("/"+value), "/")
	if value == "/upload" {
		return "/"
	}
	value = strings.TrimPrefix(value, "/upload/")
	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}
	return value
}

func collectionPath(value string) string {
	parts := strings.Split(strings.Trim(value, "/"), "/")
	var out []string
	for _, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			continue
		}
		out = append(out, part)
	}
	if len(out) == 0 {
		return "/"
	}
	return "/" + strings.Join(out, "/")
}

func hasPathParameter(value string) bool {
	for _, part := range strings.Split(value, "/") {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			return true
		}
	}
	return false
}

func confidence(score int) string {
	switch {
	case score >= 90:
		return "high"
	case score >= 70:
		return "medium"
	default:
		return "low"
	}
}

func goalWantsUpdate(goal string) bool {
	goal = strings.ToLower(goal)
	for _, word := range []string{"update", "updates", "updated", "patch", "patches", "modify", "modifies", "replace", "supports update"} {
		if strings.Contains(goal, word) {
			return true
		}
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
