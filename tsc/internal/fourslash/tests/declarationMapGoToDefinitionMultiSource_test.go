package fourslash_test

import (
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/fourslash"
	"github.com/microsoft/TypeScript/tsc/internal/testutil"
)

// Regression test for #63874: go-to-definition must use the source file named
// by the declaration-map mapping segment, not `sources[0]`. The map below is
// similar to maps shipped by npm packages whose `.d.ts` maps contain several
// sources; the statement-start segment maps to `other.ts` (`sources[0]`) while
// the identifier segment maps to `index.ts` (`sources[1]`).
func TestDeclarationMapGoToDefinitionMultiSource(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @lib: es5
// @Filename: index.ts
export const foo = 0;
// @Filename: other.ts
export const unrelated = 1;
// @Filename: indexdef.d.ts.map
{"version":3,"file":"indexdef.d.ts","sources":["other.ts","index.ts"],"names":[],"mappings":"AAAA,qBCAa,GAAG"}
// @Filename: indexdef.d.ts
export declare const foo: number;
//# sourceMappingURL=indexdef.d.ts.map
// @Filename: mymodule.ts
import { foo } from "./indexdef";
export const value = [|/*1*/foo|];`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.MarkTestAsStradaServer()
	f.VerifyBaselineGoToDefinition(t, true, "1")
}
