# 0001. Record architecture decisions

- **Status:** Accepted
- **Date:** 2026-11-06

## Context

We need a lightweight, discoverable way to record significant architecture
decisions and the reasoning behind them.

## Decision

Adopt Architecture Decision Records (ADRs) in this folder. Each ADR is a short,
numbered markdown document capturing the context, decision, and consequences of a
significant choice, following the form of this template.

## Consequences

- Decisions and their rationale become searchable and traceable over time.
- One new file per decision, named `NNNN-short-title.md`.
- Superseding an ADR does not delete it; mark the old one `Superseded by NNNN`.

## Alternatives

- A single changelog / log file — rejected: harder to link decisions to code.
- A full ADR tool — rejected: overkill for the size of this project.
