name: Feature request
about: Suggest a change to the ZIP lookup service
title: ""
labels: enhancement
assignees: ""
body:
  - type: textarea
    id: problem
    attributes:
      label: Problem to solve
    validations:
      required: true
  - type: textarea
    id: behavior
    attributes:
      label: Desired behavior and use case
    validations:
      required: true
  - type: textarea
    id: constraints
    attributes:
      label: Constraints and non-goals
  - type: textarea
    id: compatibility
    attributes:
      label: Compatibility expectations
      description: Feature parity is the default; route/API parity is only included when explicitly requested.
