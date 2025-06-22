# Where should we write the GIF?
Output ./examples/demo.gif

# Set up a terminal with font size
Set Shell bash
Set FontSize 16
Set Width 1920
Set Height 1080
Set TypingSpeed 150ms
#Set Theme "Catppuccin Macchiato"

# Sleep 500ms
Hide
Type@10ms "DEBUG=TRUE go run . ./config/60pct-keyboard.json && clear"
Enter
Show
Sleep 500ms
Sleep 1s
Type "jjkk"
Sleep 500ms
Enter
Sleep 2s
Type "Pack mu bix dd"
Backspace 8
Type "y box with fov"
Screenshot ./examples/screenshot.png
Backspace 2
Type "ive dozen liquor jugs"
Tab
Sleep 1s
Type "3"
Sleep 2s
Type ":config"
Enter
Sleep 2s
Type ":home"
Enter
Sleep 1ms
Type ":q"
Enter
Sleep 500ms
