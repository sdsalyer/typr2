# typr2

A TUI typing tutor written in Go.

A re-write and re-imagining of the original [typr](https://sr.ht/~berts/typr/) project.

![typr2 screenshot](examples/screenshot.png)

## Demo

![typr2 demo gif](examples/demo.gif)

## Project

I looked at the top Go projects on Github and found that their project
structures are all over the place. However, one of the most-starred Go projects
is the project layout listed below.

This project uses the format standard at
[golang-standards/project-layout](https://github.com/golang-standards/project-layout)

These tenets from Melkey might also be useful:

1. Keep binaries and libraries separate.
  - `internal/` for private libs
  - `pkg/` for public libs
2. Avoid nested packages (keep it flat).
3. Keep package names descriptive.
4. Put non-go files at the root of the project.
5. Keep tests in the tests folder.
