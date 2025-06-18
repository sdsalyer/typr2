/*
Keyboard Layout Editor (KLE) implementation
https://github.com/ijprest/kle-serial
*/
package main

import (
	// "encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/yosuke-furukawa/json5/encoding/json5"
	"maps"
)

// Represents a keyboard layout in KLE format
// See: https://github.com/ijprest/kle-serial?tab=readme-ov-file#keyboard-objects
type Keyboard struct {
	Meta KeyboardMetadata `json:"meta"`
	Keys []Key            `json:"keys"`
}

// Represents a keyboard's metadata in KLE format (Name, Author, etc.)
// See: https://github.com/ijprest/kle-serial?tab=readme-ov-file#keyboard-metadata
type KeyboardMetadata struct {
	Author      string `json:"author"`
	Backcolor   string `json:"backcolor"`
	Name        string `json:"name"`
	Notes       string `json:"notes"`
	Radii       string `json:"radii"`
	SwitchBrand string `json:"switchBrand"`
	SwitchMount string `json:"switchMount"`
	SwitchType  string `json:"switchType"`
}

// See: https://github.com/ijprest/kle-serial?tab=readme-ov-file#keys
type Key struct {
	Color  string   `json:"color"`
	Labels []string `json:"labels"` // An array of up to 12 labels
	/*
		These are split by "\n" from the JSON
		Positioned as such:
		[
		  0,  1,  2,
		  3,  4,  5,
		  6,  7,  8,
		  9, 10, 11 // These are "front legends"
		]
	*/

	// TODO: we are ignoring the x,y coordinates and assuming row and col position
	// X      float64 `json:"x"`
	// Y      float64 `json:"y"`

	Width  float64 `json:"width"`
	Height float64 `json:"height"`

	X2      float64 `json:"x2"`
	Y2      float64 `json:"y2"`
	Width2  float64 `json:"width2"`
	Height2 float64 `json:"height2"`

	RotationX     float64 `json:"rotation_x"`
	RotationY     float64 `json:"rotation_y"`
	RotationAngle float64 `json:"rotation_angle"`

	Decal   bool `json:"decal"`
	Ghost   bool `json:"ghost"`
	Stepped bool `json:"stepped"`
	Nub     bool `json:"nub"` // Bump for "homing" (i.e. F and J on QWERTY home row)

	Profile string `json:"profile"`

	SM string `json:"sm"` // switch mount
	SB string `json:"sb"` // switch brand
	ST string `json:"st"` // switch type

	// Additional fields for rendering
	Alignment int    `json:"alignment"`
	FontSize  int    `json:"fontSize"`
	TextColor string `json:"textColor"`
	X         int    `json:"x_pos"` // column/position in the row
	Y         int    `json:"y_pos"` // row number
}

// parseKLELayout parses the KLE JSON format into our Keyboard struct
func parseKLELayout(data []byte) (Keyboard, error) {
	var rawData []any
	if err := json5.Unmarshal(data, &rawData); err != nil {
		return Keyboard{}, fmt.Errorf("failed to parse JSON5: %w", err)
	}

	keyboard := Keyboard{
		Keys: []Key{},
	}

	// Parse metadata from first object
	if len(rawData) > 0 {
		if metaObj, ok := rawData[0].(map[string]any); ok {
			keyboard.Meta = parseMetadata(metaObj)
		}
	}

	// Parse key rows
	// TODO: while these can be decimals in the KLE layout, i'm not sure it makes
	//       sense from a TUI perspective to have fractional rows or rotated keys
	//       For now, changing to int for positional value and ignore the KLE x/y
	currentY := 0
	currentX := 0

	for _, row := range rawData {

		// Metadata is already parsed and inherently skipped by iterating
		// arrays from the JSON here
		if rowArray, ok := row.([]any); ok {
			currentX = 0
			keys := parseKeyRow(rowArray, currentX, currentY)
			keyboard.Keys = append(keyboard.Keys, keys...)
			currentY += 1
		}
	}

	return keyboard, nil
}

func parseMetadata(obj map[string]any) KeyboardMetadata {
	meta := KeyboardMetadata{}

	if name, ok := obj["name"].(string); ok {
		meta.Name = name
	}
	if author, ok := obj["author"].(string); ok {
		meta.Author = author
	}
	if notes, ok := obj["notes"].(string); ok {
		meta.Notes = notes
	}
	if radii, ok := obj["radii"].(string); ok {
		meta.Radii = radii
	}
	if switchMount, ok := obj["switchMount"].(string); ok {
		meta.SwitchMount = switchMount
	}

	return meta
}

func parseKeyRow(row []any, startX int, y int) []Key {
	var keys []Key
	currentX := startX

	// Current key properties (carried forward)
	var currentProps map[string]any

	for _, item := range row {
		switch v := item.(type) {
		case map[string]any:
			// Key properties
			currentProps = mergeProps(currentProps, v)

			// TODO: we ignore the x coordinate and assume an x position
			// Handle X offset
			// if x, ok := v["x"].(float64); ok {
			// 	currentX += x
			// }

		case string:
			// Key label - create key
			key := Key{
				X:      currentX,
				Y:      y,
				Width:  1.0,
				Height: 1.0,
			}

			// Apply current properties
			if currentProps != nil {
				applyKeyProps(&key, currentProps)
				// Reset the properties so subsequent keys don't get the wrong props
				currentProps = nil
			}

			key.Labels = parseLabels(v, key.Alignment)

			keys = append(keys, key)
			// TODO: we're incrementing position, not adding to coordinate
			// currentX += key.Width
			currentX += 1
		}
	}

	return keys
}

func mergeProps(existing, new map[string]any) map[string]any {
	if existing == nil {
		existing = make(map[string]any)
	}

	maps.Copy(existing, new)

	return existing
}

func applyKeyProps(key *Key, props map[string]any) {
	if w, ok := props["w"].(float64); ok {
		key.Width = w
	}

	if h, ok := props["h"].(float64); ok {
		key.Height = h
	}

	if c, ok := props["c"].(string); ok {
		key.Color = c
	}

	if t, ok := props["t"].(string); ok {
		key.TextColor = t
	}

	if f, ok := props["f"].(float64); ok {
		key.FontSize = int(f)
	}

	if a, ok := props["a"].(float64); ok {
		key.Alignment = int(a)
	}

	if p, ok := props["p"].(string); ok {
		key.Profile = p
	}

	if d, ok := props["d"].(bool); ok {
		key.Decal = d
	}

	if g, ok := props["g"].(bool); ok {
		key.Ghost = g
	}

	if n, ok := props["n"].(bool); ok {
		key.Nub = n
	}

	// etc...
}

// Read and parse KLE layout JSON file
func loadKeyboard(filename string) (Keyboard, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Keyboard{}, fmt.Errorf("failed to read file: %w", err)
	}

	return parseKLELayout(data)
}

// Reorder labels based on alignment flags
// See: https://github.com/ijprest/kle-serial/blob/4080386fcdcb66a391e1b4857532512f9ca4121e/index.ts#L86-L92
func reorderLabels(labels []string, alignment int) []string {
	// Map from serialized label position to normalized position,
	// depending on the alignment flags.
	labelMap := [][]int{
		// 0   1   2   3   4   5   6   7   8   9  10  11   // alignment flags
		{0, 6, 2, 8, 9, 11, 3, 5, 1, 4, 7, 10},          // 0 = no centering
		{1, 7, -1, -1, 9, 11, 4, -1, -1, -1, -1, 10},    // 1 = center x
		{3, -1, 5, -1, 9, 11, -1, -1, 4, -1, -1, 10},    // 2 = center y
		{4, -1, -1, -1, 9, 11, -1, -1, -1, -1, -1, 10},  // 3 = center x & y
		{0, 6, 2, 8, 10, -1, 3, 5, 1, 4, 7, -1},         // 4 = center front (default)
		{1, 7, -1, -1, 10, -1, 4, -1, -1, -1, -1, -1},   // 5 = center front & x
		{3, -1, 5, -1, 10, -1, -1, -1, 4, -1, -1, -1},   // 6 = center front & y
		{4, -1, -1, -1, 10, -1, -1, -1, -1, -1, -1, -1}, // 7 = center front & x & y
	}

	var retVal []string = make([]string, len(labels))
	for i := range labels {
		newIndex := labelMap[alignment][i]
		if newIndex == -1 {
			continue // Don't reorder this index
		}
		retVal[newIndex] = labels[i]
	}

	return retVal
}

func parseLabels(labelStr string, alignment int) []string {
	// Split labels by newline and trim whitespace

	labels := strings.Split(labelStr, "\n")
	for i := range 12 {
		if i >= len(labels) {
			labels = append(labels, "") // Fill missing labels with empty strings
		} else {
			// TODO: sanitize user input, may contain arbitrary HTML content
			newLabel := strings.TrimSpace(labels[i])
			if len(newLabel) == 0 {
				// Must be the SPACE key
				newLabel = "␣"
			}
			labels[i] = newLabel
		}
	}
	if len(labels) != 12 {
		panic(fmt.Sprintf("Expected 12 labels, got %d: %v", len(labels), labels))
	}

	return reorderLabels(labels, alignment)
}

// This should format the value of any marshalable type into a pretty-printed JSON5 string
func PrettyPrint[T any](v T) (string, error) {
	prettyJSON, err := json5.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}
