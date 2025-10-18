## Summary
- What does this PR change?
- Why is it needed?

## Checklist (SDL + CI Gates)
- [ ] Lint passes (golangci-lint)
- [ ] Unit tests added/updated
- [ ] Coverage meets threshold (current: ${{ vars.COVERAGE_MIN || 0 }}%)
- [ ] Govulncheck passes
- [ ] Fuzz smoke passes (if targets exist)
- [ ] Chaos smoke (if applicable)
- [ ] Structured logs follow OBSERVABILITY_SPEC (fields included)
- [ ] Threat model updated if this PR changes trust boundaries or inputs
- [ ] Plan updated (docs/PLAN.md & docs/plan.json) if scope/timeline changes

## Testing Notes
- Commands executed / how to verify locally:

## Risk & Rollback
- Risk level and mitigation
- Rollback plan
