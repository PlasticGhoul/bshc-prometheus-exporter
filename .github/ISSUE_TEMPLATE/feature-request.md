```markdown
name: Feature request
description: Suggest an idea for this project
title: "[FEATURE] "
labels: enhancement
assignees: ''

body:
  - type: markdown
    attributes:
      value: |
        ## Feature Request

  - type: input
    id: title
    attributes:
      label: Title
      description: A clear and concise title for your feature request.
      placeholder: Add a title for your feature request
    validations:
      required: true

  - type: textarea
    id: description
    attributes:
      label: Description
      description: A clear and concise description of what the problem is. Ex. I'm always frustrated when [...]
      placeholder: Describe your feature request
    validations:
      required: true

  - type: textarea
    id: solution
    attributes:
      label: Proposed Solution
      description: Describe the solution you'd like.
      placeholder: Describe the solution you'd like
    validations:
      required: true

  - type: textarea
    id: alternatives
    attributes:
      label: Alternatives
      description: Describe any alternative solutions or features you've considered.
      placeholder: Describe any alternative solutions or features you've considered

  - type: input
    id: additional-context
    attributes:
      label: Additional context
      description: Add any other context or screenshots about the feature request here.
      placeholder: Add any other context or screenshots
```