# cricket-texts

Personal cron job: picks a random **cricket (insect)** fact from [`facts.json`](facts.json) and texts it via **Gmail SMTP** using each recipient’s **carrier email-to-SMS gateway** (free, US carriers).

Stack: **Go** (stdlib only), **GitHub Actions** on a schedule + manual runs.

## What you need

1. **Gmail + App Password** (not your normal Gmail password)
   - Turn on [2-Step Verification](https://myaccount.google.com/signinoptions/two-step-verification).
   - Create an [App Password](https://myaccount.google.com/apppasswords) for “Mail”.

2. **Each person’s phone number + carrier gateway** so we can build `5551234567@gateway`

   | Carrier (US) | Gateway domain | Notes |
   |--------------|----------------|--------|
   | Verizon | `vtext.com` | Still used by some scripts, but carriers are phasing free email-to-SMS gateways out industry-wide; delivery is not guaranteed. |
   | T-Mobile | `tmomail.net` | Often unreliable; carriers are winding these gateways down industry-wide. |
   | AT&T | ~~`txt.att.net`~~ | **Shut down June 17, 2025.** Email to `*@txt.att.net` / `*@mms.att.net` no longer works. [AT&T notice](https://www.att.com/support/article/wireless/KM1061254/). |
   | US Cellular | `email.uscc.net` | Confirm with your line; policies change. |

   Format in secrets: `NUMBER:gateway`, comma-separated, e.g.  
   `5551234567:vtext.com,5559876543:tmomail.net`

   **AT&T on your own line?** This repo’s free “email-to-text” path **cannot** deliver to you anymore. Use a **non-AT&T test number**, or switch the project to an **SMS API** (e.g. Twilio) for ~pennies per message.

3. **GitHub repo** with Actions enabled.

## GitHub secrets

In the repo: **Settings → Secrets and variables → Actions → New repository secret**

| Name | Value |
|------|--------|
| `GMAIL_ADDRESS` | Your Gmail address |
| `GMAIL_APP_PASSWORD` | 16-character app password (spaces optional) |
| `RECIPIENTS` | `5551234567:vtext.com,...` as above |

## Schedule

Default: **14:00 UTC daily** in [`.github/workflows/send.yml`](.github/workflows/send.yml). Edit the `cron` line to match your local time (UTC).

GitHub’s cron can run a few minutes late; fine for a daily fact.

## Try it locally

Install [Go 1.22+](https://go.dev/dl/). Then:

1. Copy [`.env.example`](.env.example) to **`.env`** and put your real Gmail, app password, and `RECIPIENTS` there (`.env` is gitignored).
2. From the `cricket-texts` folder:

```bash
cd cricket-texts
go run .
```

The program loads **`.env` automatically** if present. Variables already set in the shell (or in GitHub Actions) are **not** overwritten.

You can still set env vars manually instead of `.env` (`set` / `$env:...` / `export`).

### Windows: `go` not found in PowerShell or Git Bash

The real toolchain is normally `C:\Program Files\Go\bin\go.exe`. If `PATH` still points at `C:\Users\YOU\go\bin` (often only **`gopls.exe`** lives there), shells may never see `go`.

This repo doesn’t configure your OS, but profiles were added to prepend `C:\Program Files\Go\bin` when missing:

- `Documents\PowerShell\Microsoft.PowerShell_profile.ps1` (PowerShell 7 / `pwsh`)
- `Documents\WindowsPowerShell\Microsoft.PowerShell_profile.ps1` (Windows PowerShell 5)

**Open a new terminal** after that. If PowerShell blocks scripts, run once (your user account):  
`Set-ExecutionPolicy -Scope CurrentUser RemoteSigned`

**Git Bash:** `~/.bash_profile` prepends `/c/Program Files/Go/bin`. New Git Bash window after saving.

**Always works (any shell):**  
`& "C:\Program Files\Go\bin\go.exe" version`

## Edit facts

All strings live in [`facts.json`](facts.json). Keep lines short (under ~160 characters) so they fit one SMS on most carriers.

## Gmail says `535 5.7.8 Username and Password not accepted`

The code is hitting Gmail’s SMTP; Google is rejecting auth. Typical causes:

1. **Use an App Password, not your normal Gmail password.** Enable [2‑Step Verification](https://myaccount.google.com/signinoptions/two-step-verification), then create an [App Password](https://myaccount.google.com/apppasswords) for “Mail”.
2. **Address must match the Google account** that generated the App Password (`GMAIL_ADDRESS` = that account’s Gmail, including `@gmail.com` unless it’s Workspace).
3. **Google Workspace / school / work**: admins can **disable** App Passwords or SMTP — then consumer instructions don’t apply.
4. **Advanced Protection Program**: App Passwords are **not** available.
5. **Copy/paste quirks**: regenerate the App Password and paste the 16 characters again; stray spaces/quotes in `.env` are stripped by the program, but wrong characters still fail.

After fixing, run `go run .` again.

## Notes

- Gateways are **best-effort**; carriers can throttle or change behavior.
- This sends **one email** with multiple `To:` addresses (all carrier addresses). If you prefer one email per number, we can change that later.
- `facts.json` is **embedded** in the binary with `go:embed`; no extra files needed at runtime.
