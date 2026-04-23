# CRN Typer v0.1.0

A desktop helper that automatically a saved list of CRNs into the active window, pressing **Tab** after each CRN.

It’s built for the LAU registration page that provide multiple “CRN” input fields.

Check `tester.html` for a local test page that you can test the program on.

## Platform support

| OS | Status | Trigger key | CRN file |
| --- | --- | --- | --- |
| macOS | Supported | Backspace | `~/Documents/crns.txt` |
| Windows | Supported | Backspace | `%USERPROFILE%\Documents\crns.txt` |

Linux and other operating systems are not supported.

## How it works

- You enter/paste your CRNs in the app UI.
- Click **Load** to “arm” typing.
- Put your cursor in the first CRN field (e.g., in your browser).
- Press **Backspace**.
- The app types each CRN and then presses Tab to move to the next field.

CRNs are persisted as plain text (one per line) in:

- macOS: `~/Documents/crns.txt`
- Windows: `%USERPROFILE%\Documents\crns.txt`
- If the home directory can’t be resolved: `./crns.txt`

## Requirements

- Go **1.23+**
- A desktop environment (this is a GUI app built with Fyne)
- Windows builds from source need a working cgo-capable C compiler, such as MinGW-w64.

## Permissions

### macOS permissions

This app uses global keyboard hooks + synthetic keyboard events (via `robotgo`/`gohook`). On macOS you’ll typically need to allow the app (or the terminal you run it from) under:

- **System Settings → Privacy & Security → Accessibility**
- **System Settings → Privacy & Security → Input Monitoring** if your macOS version prompts for it

If typing or the Backspace trigger doesn’t work, permissions are the first thing to check.

### Windows permissions

No extra permission is normally required when typing into regular desktop apps. If browser is running as administrator, run CRN Typer as administrator too;

## Run

From the repo root:

```bash
go run .
```

Or build a binary:

```bash
go build -o crntyper .
./crntyper
```

## Usage

1. Launch the app.
2. Enter your CRNs (the UI automatically adds a new empty row as you type).
3. Click **Load**.
4. Focus the first CRN input field in the target application/webpage.
5. Press **Backspace** to type the full list.

Tip: If you edit CRNs while the app is “loaded”, it auto-saves to `crns.txt`.

## Test locally

Open `tester.html` in your browser (double-click it or use “Open File…”), click into the first CRN box, then follow the steps above.

## Disclosure
This project was vibe-coded, but I personally reviewed the entire codebase and tested it thoroughly.

## License

Apache-2.0. See `LICENSE`.
