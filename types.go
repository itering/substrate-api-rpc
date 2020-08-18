package substrate

import (
	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types"
)

func RegCustomTypes(sourceCode []byte) {
	types.RegCustomTypes(source.LoadTypeRegistry(sourceCode))
}
