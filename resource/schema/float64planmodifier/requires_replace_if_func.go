package float64planmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequiresReplaceIfFunc is a conditional function used in the RequiresReplaceIf
// plan modifier to determine whether the attribute requires replacement.
type RequiresReplaceIfFunc func(context.Context, planmodifier.Float64Request, *RequiresReplaceIfFuncResponse)

// RequiresReplaceIfFuncResponse is the response type for a RequiresReplaceIfFunc.
type RequiresReplaceIfFuncResponse struct {
	// Diagnostics report errors or warnings related to this logic. An empty
	// or unset slice indicates success, with no warnings or errors generated.
	Diagnostics diag.Diagnostics

	// RequiresReplace should be enabled if the resource should be replaced.
	RequiresReplace bool
}
