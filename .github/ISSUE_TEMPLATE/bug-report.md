```markdown
name: Bug report
description: Create a report to help us improve
title: "[BUG] "
labels: bug
assignees: ''

body:
  - type: markdown
    attributes:
      value: |
        ## Bug Report

  - type: input
    id: title
    attributes:
      label: Title
      description: A clear and concise title of what the bug is.
      placeholder: Bug title
    validations:
      required: true

  - type: textarea
    id: description
    attributes:
      label: Description
      description: A clear and concise description of what the bug is.
      placeholder: Describe the bug
    validations:
      required: true

  - type: textarea
    id: steps-to-reproduce
    attributes:
      label: Steps to reproduce
      description: |
        Steps to reproduce the behavior:
        1. Some step
        2. Another step
        3. And so on
        4. ...
      placeholder: Steps to reproduce the bug
    validations:
      required: true

  - type: textarea
    id: expected-behavior
    attributes:
      label: Expected behavior
      description: A clear and concise description of what you expected to happen.
      placeholder: Expected behavior
    validations:
      required: true

  - type: textarea
    id: screenshots
    attributes:
      label: Screenshots
      description: If applicable, add screenshots to help explain your problem.
      placeholder: Add screenshots or images

  - type: input
    id: environment
    attributes:
      label: Environment
      description: |
        Please complete the following information:
        - OS: [e.g. Windows, MacOS]
        - Browser [e.g. chrome, safari]
        - Version [e.g. 22]
      placeholder: Environment details
    validations:
      required: true

  - type: textarea
    id: additional-context
    attributes:
      label: Additional context
      description: Add any other context about the problem here.
      placeholder: Additional context
```