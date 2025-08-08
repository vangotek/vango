# VanGo Theme Application Fix

## Issues Identified

1. **Template Loading Priority**: Default layouts are loaded first, overriding theme templates
2. **Missing Base Template Support**: No support for Hugo-style template inheritance
3. **Template Lookup Logic**: Doesn't properly resolve theme templates
4. **Asset Path Resolution**: Theme assets may not be correctly linked

## Required Fixes

### 1. Fix Template Loading Order in `internal/template/engine.go`

```go
// LoadTemplates should prioritize theme templates over default layouts
func (e *Engine) LoadTemplates(themeLayoutDir string) error {
	// Load theme templates first (higher priority)
	if themeLayoutDir != "" && themeLayoutDir != e.config.LayoutDir {
		if err := e.parseAndAddTemplates(themeLayoutDir); err != nil {
			return fmt.Errorf("failed to parse theme templates: %w", err)
		}
	}

	// Then load default templates (lower priority - won't override existing)
	if err := e.parseAndAddTemplatesWithOverride(e.config.LayoutDir, false); err != nil {
		return fmt.Errorf("failed to parse default templates: %w", err)
	}

	return nil
}
```

### 2. Add Template Hierarchy Support

The engine needs to support base templates and template blocks. The modern-app theme uses:
- `baseof.html` - Base template with blocks
- `single.html` - Extends baseof with `{{ define "content" }}`
- `list.html` - Extends baseof with `{{ define "content" }}`

### 3. Fix Template Name Resolution

Current template lookup logic needs to handle:
- Base template inheritance
- Theme-specific template paths
- Fallback to default templates

### 4. Ensure Theme Assets Are Properly Linked

The `themeAsset` function should correctly resolve paths to copied theme assets.

## Implementation Priority

1. **HIGH**: Fix template loading order
2. **HIGH**: Add base template support 
3. **MEDIUM**: Improve template lookup logic
4. **LOW**: Enhanced asset path resolution

## Expected Result

After fixes, the modern-app theme should:
- Use its own `baseof.html` base template
- Apply the modern CSS styling
- Show proper navigation and layout structure
- Include theme-specific features (dark mode toggle, etc.)

## Test Command

After implementing fixes:
```bash
cd test/vango-test
../../vango build
# Check that public/welcome/index.html uses modern-app theme templates
```
