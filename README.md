# Wacky Message. Translation and creative modification of clipboard text, and pasting into active input fields

Learning to use Go by roughly implementing crazy ideas

This project processes text from the clipboard by template, translates it into various languages, and inserts the translated text into an active input field in an open window.

## Description

The project allows users to automatically translate the text in the clipboard and paste the translated text into the corresponding input field in active applications.
## How to Use

1. **Run the Application**:
    - Build and run the application using the following command:
      ```bash
      go run main.go
      ```

2. **Clipboard Monitoring**:
    - The application continuously monitors the clipboard for text that ends with the trigger suffix `!!1`.
    - If a matching message is found, it modifies the message and sends it to the active window.

3. **Using the Hotkey**:
- The application supports the `F2` hotkey for processing text.
- To use the hotkey:
    1. Select the windows where the hotkey will be active by using the "Choose windows where hotkey will be active" option in the tray menu.
    2. Press `F1` while in one of the selected windows.
    3. The application will:
        - Copy the selected text to the clipboard.
        - Process the text (e.g., translation, adding emojis).
        - Paste the modified text back into the active input field.

4. **System Tray**:
    - The application runs in the system tray.
    - Right-click the tray icon to access the "Quit" option to exit the application.

---

## Compatibility

- **Operating System**: Windows only (relies on `user32.dll` for clipboard and keyboard interactions).
- **Go Version**: Requires Go 1.18 or later.
- **Dependencies**:
    - Ensure all dependencies in `go.mod` are installed by running:
      ```bash
      go mod tidy
      ```
- **Limitations**:
    - The application may not work properly in environments where clipboard or keyboard access is restricted (e.g., virtual machines or sandboxed environments).

