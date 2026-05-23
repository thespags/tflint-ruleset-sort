package custom

import (
	"fmt"
	"strings"

	"github.com/thespags/tflint-ruleset-sort/config"
)

func parseConfigResource(resource *config.Resource) (*Resource, error) {
	out := &Resource{
		SkipSortLiterals: append([]string(nil), resource.SkipSortLiterals...),
	}

	if len(resource.Keys) == 0 {
		return out, nil
	}

	keys := strings.Split(resource.Keys[0], ".")
	keyAttributes := make([]string, 0, len(resource.Keys))
	keyBlocks := keys[:len(keys)-1]

	for _, key := range resource.Keys {
		keys = strings.Split(key, ".")
		for i := range len(keys) - 1 {
			if keys[i] != keyBlocks[i] {
				return nil, fmt.Errorf(
					"invalid configuration for `%s`: all keys must have the same prefix `%s`: unexpected `%s`",
					resource.Kind,
					strings.Join(keyBlocks, "."),
					strings.Join(keys[:i+1], "."),
				)
			}
		}

		keyAttributes = append(keyAttributes, keys[len(keys)-1])
	}

	out.KeyBlocks = keyBlocks
	out.KeyAttributes = keyAttributes

	return out, nil
}
