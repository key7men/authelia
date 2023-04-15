package configuration

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"

	"github.com/authelia/authelia/v4/internal/configuration/schema"
	"github.com/authelia/authelia/v4/internal/utils"
)

func koanfGetKeys(ko *koanf.Koanf) (keys []string) {
	keys = ko.Keys()

	for key, value := range ko.All() {
		slc, ok := value.([]any)
		if !ok {
			continue
		}

		for _, item := range slc {
			m, mok := item.(map[string]any)
			if !mok {
				continue
			}

			for k := range m {
				full := fmt.Sprintf("%s[].%s", key, k)
				if !utils.IsStringInSlice(full, keys) {
					keys = append(keys, full)
				}
			}
		}
	}

	return keys
}

func koanfRemapKeys(val *schema.StructValidator, ko *koanf.Koanf, ds map[string]Deprecation) (final *koanf.Koanf, err error) {
	keys := ko.All()

	keys = koanfRemapKeysStandard(keys, val, ds)
	keys = koanfRemapKeysMapped(keys, val, ds)

	final = koanf.New(".")

	if err = final.Load(confmap.Provider(keys, "."), nil); err != nil {
		return nil, err
	}

	return final, nil
}

func koanfRemapKeysStandard(keys map[string]any, val *schema.StructValidator, ds map[string]Deprecation) (keysFinal map[string]any) {
	var (
		ok    bool
		d     Deprecation
		key   string
		value any
	)

	keysFinal = make(map[string]any)

	for key, value = range keys {
		if d, ok = ds[key]; ok {
			if !d.AutoMap {
				val.Push(fmt.Errorf("invalid configuration key '%s' was replaced by '%s'", d.Key, d.NewKey))

				keysFinal[key] = value

				continue
			} else {
				val.PushWarning(fmt.Errorf("configuration key '%s' is deprecated in %s and has been replaced by '%s': "+
					"this has been automatically mapped for you but you will need to adjust your configuration to remove this message", d.Key, d.Version.String(), d.NewKey))
			}

			if !mapHasKey(d.NewKey, keys) && !mapHasKey(d.NewKey, keysFinal) {
				if d.MapFunc != nil {
					keysFinal[d.NewKey] = d.MapFunc(value)
				} else {
					keysFinal[d.NewKey] = value
				}
			}

			continue
		}

		keysFinal[key] = value
	}

	return keysFinal
}

func koanfRemapKeysMapped(keys map[string]any, val *schema.StructValidator, ds map[string]Deprecation) (keysFinal map[string]any) {
	var (
		key           string
		value         any
		slc, slcFinal []any
		ok            bool
		m             map[string]any
		d             Deprecation
	)

	keysFinal = make(map[string]any)

	for key, value = range keys {
		if slc, ok = value.([]any); !ok {
			keysFinal[key] = value

			continue
		}

		slcFinal = make([]any, len(slc))

		for i, item := range slc {
			if m, ok = item.(map[string]any); !ok {
				slcFinal[i] = item

				continue
			}

			itemFinal := make(map[string]any)

			for subkey, element := range m {
				prefix := fmt.Sprintf("%s[].", key)

				fullKey := prefix + subkey

				if d, ok = ds[fullKey]; ok {
					if !d.AutoMap {
						val.Push(fmt.Errorf("invalid configuration key '%s' was replaced by '%s'", d.Key, d.NewKey))

						itemFinal[subkey] = element

						continue
					} else {
						val.PushWarning(fmt.Errorf("configuration key '%s' is deprecated in %s and has been replaced by '%s': "+
							"this has been automatically mapped for you but you will need to adjust your configuration to remove this message", d.Key, d.Version.String(), d.NewKey))
					}

					newkey := strings.Replace(d.NewKey, prefix, "", 1)

					if !mapHasKey(newkey, m) && !mapHasKey(newkey, itemFinal) {
						if d.MapFunc != nil {
							itemFinal[newkey] = d.MapFunc(element)
						} else {
							itemFinal[newkey] = element
						}
					}
				} else {
					itemFinal[subkey] = element
				}
			}

			slcFinal[i] = itemFinal
		}

		keysFinal[key] = slcFinal
	}

	return keysFinal
}
