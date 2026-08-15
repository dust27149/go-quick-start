package merge

import "reflect"

// MergeNonEmpty 通用合并：patch 中非空字段覆盖 base。
func MergeNonEmpty[T any](base, patch T) T {
	baseVal := reflect.ValueOf(base)
	patchVal := reflect.ValueOf(patch)
	mergedVal := mergeValue(baseVal, patchVal)

	merged, ok := mergedVal.Interface().(T)
	if !ok {
		return base
	}
	return merged
}

// mergeValue 递归合并两个 reflect.Value，返回合并后的结果。
func mergeValue(base, patch reflect.Value) reflect.Value {
	if !patch.IsValid() {
		return base
	}

	if base.Kind() != patch.Kind() {
		if isEmptyValue(patch) {
			return base
		}
		return patch
	}

	switch patch.Kind() {
	case reflect.Struct:
		result := reflect.New(base.Type()).Elem()
		result.Set(base)
		for i := 0; i < patch.NumField(); i++ {
			field := result.Field(i)
			if !field.CanSet() {
				continue
			}
			field.Set(mergeValue(field, patch.Field(i)))
		}
		return result
	case reflect.Pointer:
		if patch.IsNil() {
			return base
		}
		if base.IsNil() {
			return patch
		}
		if base.Elem().Kind() == reflect.Struct && patch.Elem().Kind() == reflect.Struct {
			mergedElem := mergeValue(base.Elem(), patch.Elem())
			result := reflect.New(mergedElem.Type())
			result.Elem().Set(mergedElem)
			return result
		}
		return patch
	default:
		if isEmptyValue(patch) {
			return base
		}
		return patch
	}
}

// isEmptyValue 检查 reflect.Value 是否为空值。
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Slice, reflect.Map:
		return v.IsNil() || v.Len() == 0
	default:
		return v.IsZero()
	}
}
