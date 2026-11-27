name: Bug report
about: Report behavior that does not match the documented service
title: ""
labels: bug
assignees: ""
body:
  - type: markdown
    attributes:
      value: >-
        Do not post vulnerability details publicly. Report them privately at
        https://github.com/apimgr/zipcodes/security/advisories/new.
  - type: input
    id: version
    attributes:
      label: Version or commit
    validations:
      required: true
  - type: input
    id: platform
    attributes:
      label: Operating system and architecture
    validations:
      required: true
  - type: textarea
    id: expected
    attributes:
      label: Expected behavior
    validations:
      required: true
  - type: textarea
    id: actual
    attributes:
      label: Actual behavior and reproduction steps
      description: Redact credentials, personal data, hostnames, and other sensitive values.
    validations:
      required: true
  - type: checkboxes
    id: security
    attributes:
      label: Security impact
      options:
        - label: This may have security impact; I will report it privately instead of posting details here.
