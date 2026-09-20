# Result V5 - Structured Security Alternatives

M27 replaces the flattened operation credential list with
`authoring.prompt-context.v2`. Each operation now retains ordered alternative
binding sets: alternatives are OR, bindings inside one set are AND, and an
empty set explicitly permits anonymous access. The version break prevents v1
JSON from silently losing authentication meaning. OpenUdon and Ramen own the
product-specific selection, deferral, and artifact policy.
