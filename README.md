# atlas-cli

Is an open-source solution designed to simplify working with Atlassian products through a set of convenient command-line utilities. This toolset supports Confluence and Jira, enabling efficient data retrieval, export, and automation.

## Features

### 1. Upload Worklogs to Jira from Local File
_Upload daily worklogs into Jira from YAML files. Only logs with the `issue_key` are uploaded._

`go run main.py uw {path_to_daily_worklogs.yml} [-p --push]`
- `-p --push`: If True, uploads worklogs to Jira. If False, only displays the logs (default: False).

**Supported Formats:**
- YAML/YML

**YAML Example:**
```yml
---
date: '2025-05-02'
worklogs:
- started_at: '07:00'
  ended_at: '07:30'
  comment: I had breakfast 
  tags:
    - meal
- started_at: '07:30'
  ended_at: '10:20'
  comment: Worked on the task - integration with an external system.
  issue_key: KAN-1
  tags:
    - work
```
- `issue_key`: The Jira issue key (optional). Only logs with this key will be uploaded.
- Other fields (`started_at`, `ended_at`, `comment`) are required.

### 2. Retrieve Information About Confluence Spaces

_Fetch all spaces in your Confluence instance, displaying their - names and space keys._

`go run main.py sk [-b --batchSize=123]`

- `-b --batchSize`: Number of spaces to retrieve per page (default: 100).

## Installation

1. Get your API token: [Atlassian API Tokens](https://id.atlassian.com/manage-profile/security/api-tokens)
2. Rename `.env.example` to `.env`:
3. Update `.env` with your credentials.
4. Install dependencies: ```go get .```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or new features.