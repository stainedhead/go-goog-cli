# Setup Guide: Google OAuth2 Configuration

This guide walks you through setting up OAuth2 credentials to use `goog` with your Google account.

## Overview

```
+------------------+     +-------------------+     +------------------+
|  Google Cloud    |     |    goog CLI       |     |   Your Google    |
|  Console         | --> |    Application    | --> |   Account        |
|  (OAuth Setup)   |     |                   |     |   (Gmail/Cal)    |
+------------------+     +-------------------+     +------------------+
        |                         |
        v                         v
  Client ID/Secret          Access Tokens
  (one-time setup)          (per account)
```

## Prerequisites

- A Google account (personal Gmail or Google Workspace)
- Access to [Google Cloud Console](https://console.cloud.google.com/)
- `goog` CLI installed (`go install ./cmd/goog` or `go build -o bin/goog ./cmd/goog`)

## Step 1: Create a Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)

2. Click the project dropdown at the top of the page:

   ```
   +------------------------------------------------------+
   | Google Cloud    [Select a project ▼]    🔔  ?  [=]   |
   +------------------------------------------------------+
   ```

3. Click **"New Project"** in the dialog:

   ```
   +----------------------------------+
   |  Select a project                |
   |  +-----------+  +-------------+  |
   |  | RECENT    |  | ALL         |  |
   |  +-----------+  +-------------+  |
   |                                  |
   |  No projects yet                 |
   |                                  |
   |  [NEW PROJECT]  [OPEN]           |
   +----------------------------------+
   ```

4. Enter project details:
   - **Project name**: `goog-cli` (or any name you prefer)
   - **Organization**: Leave as default or select your organization
   - Click **"Create"**

5. Wait for the project to be created (notification will appear), then select it.

## Step 2: Enable Gmail and Calendar APIs

1. In the Cloud Console, go to **"APIs & Services"** > **"Library"**

   ```
   Navigation menu (☰)
   ├── APIs & Services
   │   ├── Enabled APIs & services
   │   ├── Library            <-- Click here
   │   ├── Credentials
   │   └── OAuth consent screen
   ```

2. Search for and enable these APIs:

   | API | Search Term | Enable Button |
   |-----|-------------|---------------|
   | Gmail API | `gmail api` | Click result → **"Enable"** |
   | Google Calendar API | `calendar api` | Click result → **"Enable"** |

   ```
   +------------------------------------------+
   |  Gmail API                    [ENABLE]   |
   |  ----------------------------------------|
   |  Gmail API - Read, send, and manage      |
   |  email messages and labels               |
   +------------------------------------------+
   ```

## Step 3: Configure OAuth Consent Screen

1. Go to **"APIs & Services"** > **"OAuth consent screen"**

2. Select **User Type**:

   | Type | Use Case |
   |------|----------|
   | **Internal** | Google Workspace users only (no verification needed) |
   | **External** | Personal Gmail accounts (requires test users or verification) |

   For personal use, select **External** and click **"Create"**.

3. Fill in the **App Information**:

   ```
   +------------------------------------------+
   |  OAuth consent screen                    |
   |  ----------------------------------------|
   |  App name:        [goog-cli            ] |
   |  User support:    [your.email@gmail.com] |
   |  ----------------------------------------|
   |  Developer contact:                      |
   |  Email:           [your.email@gmail.com] |
   +------------------------------------------+
   ```

   Click **"Save and Continue"**.

4. **Scopes** - Click **"Add or Remove Scopes"** and add:

   ```
   +------------------------------------------------------------------+
   | Selected scopes:                                                  |
   |------------------------------------------------------------------|
   | [x] .../auth/gmail.modify        Read/write access to Gmail      |
   | [x] .../auth/gmail.send          Send emails                     |
   | [x] .../auth/calendar            Full calendar access            |
   | [x] .../auth/userinfo.email      View email address              |
   +------------------------------------------------------------------+
   ```

   Or manually enter these scope URIs:
   ```
   https://www.googleapis.com/auth/gmail.modify
   https://www.googleapis.com/auth/gmail.send
   https://www.googleapis.com/auth/calendar
   https://www.googleapis.com/auth/userinfo.email
   ```

   Click **"Update"**, then **"Save and Continue"**.

5. **Test Users** (External apps only):

   ```
   +------------------------------------------+
   |  Test users                              |
   |  ----------------------------------------|
   |  + ADD USERS                             |
   |                                          |
   |  Add the Google accounts that can        |
   |  access this app while in testing mode   |
   |                                          |
   |  Email: [your.email@gmail.com    ] [Add] |
   +------------------------------------------+
   ```

   Add your email address and click **"Save and Continue"**.

6. Review the summary and click **"Back to Dashboard"**.

## Step 4: Create OAuth Credentials

1. Go to **"APIs & Services"** > **"Credentials"**

2. Click **"+ CREATE CREDENTIALS"** > **"OAuth client ID"**:

   ```
   +------------------------------------------+
   |  + CREATE CREDENTIALS ▼                  |
   |  ----------------------------------------|
   |  > API key                               |
   |  > OAuth client ID        <-- Select    |
   |  > Service account                       |
   +------------------------------------------+
   ```

3. Configure the OAuth client:

   ```
   +------------------------------------------+
   |  Create OAuth client ID                  |
   |  ----------------------------------------|
   |  Application type:                       |
   |  [Desktop app                        ▼]  |
   |                                          |
   |  Name:                                   |
   |  [goog-cli-desktop                    ]  |
   +------------------------------------------+
   ```

   - **Application type**: `Desktop app`
   - **Name**: `goog-cli-desktop`
   - Click **"Create"**

4. **Copy your credentials** from the dialog:

   ```
   +--------------------------------------------------+
   |  OAuth client created                            |
   |  ------------------------------------------------|
   |  Your Client ID:                                 |
   |  [123456789-abc123.apps.googleusercontent.com]   |
   |                                            [📋]  |
   |                                                  |
   |  Your Client Secret:                             |
   |  [GOCSPX-abcdef123456...]                        |
   |                                            [📋]  |
   |                                                  |
   |  [DOWNLOAD JSON]              [OK]               |
   +--------------------------------------------------+
   ```

   **IMPORTANT**: Save these credentials immediately! As of 2025, Google only shows the client secret once at creation time.

## Step 5: Configure Environment Variables

Set the OAuth credentials as environment variables:

### macOS / Linux (bash/zsh)

Add to your `~/.bashrc`, `~/.zshrc`, or `~/.profile`:

```bash
export GOOG_CLIENT_ID="your-client-id.apps.googleusercontent.com"
export GOOG_CLIENT_SECRET="GOCSPX-your-client-secret"
```

Then reload:
```bash
source ~/.bashrc  # or ~/.zshrc
```

### Windows (PowerShell)

```powershell
# Temporary (current session)
$env:GOOG_CLIENT_ID = "your-client-id.apps.googleusercontent.com"
$env:GOOG_CLIENT_SECRET = "GOCSPX-your-client-secret"

# Permanent (user environment)
[Environment]::SetEnvironmentVariable("GOOG_CLIENT_ID", "your-client-id.apps.googleusercontent.com", "User")
[Environment]::SetEnvironmentVariable("GOOG_CLIENT_SECRET", "GOCSPX-your-client-secret", "User")
```

### Windows (Command Prompt)

```cmd
setx GOOG_CLIENT_ID "your-client-id.apps.googleusercontent.com"
setx GOOG_CLIENT_SECRET "GOCSPX-your-client-secret"
```

### Verify Configuration

```bash
echo $GOOG_CLIENT_ID
# Should output: your-client-id.apps.googleusercontent.com
```

## Step 6: Authenticate with goog

Run the login command:

```bash
goog auth login
```

This will:

1. **Open your browser** to Google's consent screen:

   ```
   +--------------------------------------------------+
   |  Google                                          |
   |  ------------------------------------------------|
   |  Sign in to continue to goog-cli                 |
   |                                                  |
   |  [your.email@gmail.com                        ]  |
   |                                                  |
   |  [Continue]                                      |
   +--------------------------------------------------+
   ```

2. **Request permissions** (you'll see the scopes you configured):

   ```
   +--------------------------------------------------+
   |  goog-cli wants to access your Google Account    |
   |  ------------------------------------------------|
   |  This will allow goog-cli to:                    |
   |                                                  |
   |  [x] Read, compose, send, and permanently        |
   |      delete all your email from Gmail            |
   |                                                  |
   |  [x] See, edit, share, and permanently delete    |
   |      all the calendars you can access            |
   |                                                  |
   |  [x] See your primary Google Account email       |
   |      address                                     |
   |                                                  |
   |  [Cancel]                    [Allow]             |
   +--------------------------------------------------+
   ```

3. **Redirect to localhost** after you click "Allow":

   ```
   +--------------------------------------------------+
   |  Authentication Successful!                      |
   |  ------------------------------------------------|
   |  You have successfully authenticated with Google.|
   |  You can close this window and return to the     |
   |  terminal.                                       |
   +--------------------------------------------------+
   ```

4. **Store tokens** securely in your system keyring.

### Terminal Output

```
$ goog auth login
Opening browser for Google authentication...
Waiting for authorization...
Successfully authenticated as your.email@gmail.com
Account added: default (your.email@gmail.com)
```

## Step 7: Verify Authentication

Check your authentication status:

```bash
goog auth status
```

Output:
```
Account: default
Email: your.email@gmail.com
Token Status: Valid
Expires: 2024-01-15 15:30:00 UTC
Scopes:
  - https://www.googleapis.com/auth/gmail.modify
  - https://www.googleapis.com/auth/calendar
  - https://www.googleapis.com/auth/userinfo.email
```

Test with a simple command:

```bash
# List recent emails
goog mail list --max-results 5

# Show today's calendar
goog cal today
```

## Multi-Account Setup

Add additional accounts with aliases:

```bash
# Add work account
goog auth login --account work

# Add personal account with specific scopes
goog auth login --account personal --scopes gmail.readonly,calendar

# List all accounts
goog account list

# Switch default account
goog account switch work

# Use specific account for a command
goog mail list --account personal
```

## Scope Reference

The `--scopes` flag accepts these shorthand values:

| Shorthand | Full Scope | Access Level |
|-----------|------------|--------------|
| `gmail` | gmail.readonly | Read-only mail access |
| `gmail.readonly` | gmail.readonly | Read-only mail access |
| `gmail.send` | gmail.send | Send emails only |
| `gmail.modify` | gmail.modify | Read/write (no permanent delete) |
| `gmail.compose` | gmail.compose | Create and modify drafts |
| `gmail.labels` | gmail.labels | Manage labels only |
| `calendar` | calendar.readonly | Read-only calendar |
| `calendar.readonly` | calendar.readonly | Read-only calendar |
| `calendar.events` | calendar.events | Manage events only |
| `calendar.full` | calendar | Full calendar access |

### Recommended Scope Combinations

| Use Case | Scopes |
|----------|--------|
| Full access | `gmail.modify,calendar.full` (default) |
| Read-only | `gmail.readonly,calendar.readonly` |
| Email only | `gmail.modify` |
| Calendar only | `calendar.full` |
| AI agent (minimal) | `gmail.readonly,calendar.readonly` |

## Troubleshooting

### "Access blocked: This app's request is invalid"

**Cause**: OAuth consent screen not configured or credentials mismatch.

**Fix**: Verify your OAuth consent screen is set up and your Client ID matches.

### "Error 403: access_denied"

**Cause**: Your email is not in the test users list.

**Fix**: Add your email to test users in OAuth consent screen settings.

### "GOOG_CLIENT_ID environment variable is not set"

**Cause**: Environment variables not configured.

**Fix**: Set the environment variables as shown in Step 5 and restart your terminal.

### "OAuth error: redirect_uri_mismatch"

**Cause**: The redirect URI doesn't match what's configured in Google Cloud.

**Fix**: Desktop apps should auto-configure correctly. If issues persist, check that your OAuth client type is "Desktop app".

### Token Expired / Refresh Failed

**Cause**: Tokens expire after 7 days in testing mode, or refresh token was revoked.

**Fix**: Re-authenticate:
```bash
goog auth logout
goog auth login
```

### Port 8085 Already in Use

**Cause**: Another application is using the default callback port.

**Fix**: Set a different port:
```bash
export GOOG_REDIRECT_PORT=8086
goog auth login
```

## Security Notes

1. **Never share** your Client Secret or commit it to version control
2. **Tokens are stored** in your system's secure keyring (Keychain on macOS, Credential Manager on Windows)
3. **Refresh tokens** allow long-term access - revoke them at [myaccount.google.com/permissions](https://myaccount.google.com/permissions) if compromised
4. **Test mode** limits tokens to 7 days - publish your app for production use

## Publishing Your App (Optional)

For long-lived tokens without the 7-day expiration:

1. Go to **OAuth consent screen** in Google Cloud Console
2. Click **"Publish App"**
3. For sensitive scopes (gmail.modify, calendar), you'll need to:
   - Verify domain ownership
   - Submit for Google review
   - Provide privacy policy

For personal use, re-authenticating every 7 days is simpler than the verification process.

## References

- [Google OAuth 2.0 for Desktop Apps](https://developers.google.com/identity/protocols/oauth2/native-app)
- [Configure OAuth Consent Screen](https://developers.google.com/workspace/guides/configure-oauth-consent)
- [Gmail API Scopes](https://developers.google.com/workspace/gmail/api/auth/scopes)
- [Calendar API Scopes](https://developers.google.com/workspace/calendar/api/auth)
- [Google Cloud Console](https://console.cloud.google.com/)
