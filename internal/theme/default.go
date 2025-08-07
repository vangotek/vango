// internal/theme/default.go
package theme

import _ "embed"


var defaultSingleTemplate string


var defaultListTemplate string


var defaultCSS string

func (tm *ThemeManager) GetDefaultTheme() *Theme {
    return &Theme{
        Name: "default",
        Version: "1.0.0",
        Description: "Built-in default theme for Vango",
        Author: "Vango Team",
        Templates: map[string]string{
            "layouts/_default/single.html": defaultSingleTemplate,
            "layouts/_default/list.html": defaultListTemplate,
        },
        CSS: defaultCSS,
    }
}