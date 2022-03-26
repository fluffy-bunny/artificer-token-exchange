// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package claimsprovider

import (
	"reflect"

	di "github.com/fluffy-bunny/sarulabsdi"
)

// ReflectTypeIClaimsProvider used when your service claims to implement IClaimsProvider
var ReflectTypeIClaimsProvider = di.GetInterfaceReflectType((*IClaimsProvider)(nil))

// AddSingletonIClaimsProvider adds a type that implements IClaimsProvider
func AddSingletonIClaimsProvider(builder *di.Builder, implType reflect.Type, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	di.AddSingleton(builder, implType, implementedTypes...)
}

// AddSingletonIClaimsProviderWithMetadata adds a type that implements IClaimsProvider
func AddSingletonIClaimsProviderWithMetadata(builder *di.Builder, implType reflect.Type, metaData map[string]interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	di.AddSingletonWithMetadata(builder, implType, metaData, implementedTypes...)
}

// AddSingletonIClaimsProviderByObj adds a prebuilt obj
func AddSingletonIClaimsProviderByObj(builder *di.Builder, obj interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	di.AddSingletonWithImplementedTypesByObj(builder, obj, implementedTypes...)
}

// AddSingletonIClaimsProviderByObjWithMetadata adds a prebuilt obj
func AddSingletonIClaimsProviderByObjWithMetadata(builder *di.Builder, obj interface{}, metaData map[string]interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	di.AddSingletonWithImplementedTypesByObjWithMetadata(builder, obj, metaData, implementedTypes...)
}

// AddSingletonIClaimsProviderByFunc adds a type by a custom func
func AddSingletonIClaimsProviderByFunc(builder *di.Builder, implType reflect.Type, build func(ctn di.Container) (interface{}, error), implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	di.AddSingletonWithImplementedTypesByFunc(builder, implType, build, implementedTypes...)
}

// AddSingletonIClaimsProviderByFuncWithMetadata adds a type by a custom func
func AddSingletonIClaimsProviderByFuncWithMetadata(builder *di.Builder, implType reflect.Type, build func(ctn di.Container) (interface{}, error), metaData map[string]interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	di.AddSingletonWithImplementedTypesByFuncWithMetadata(builder, implType, build, metaData, implementedTypes...)
}

// AddTransientIClaimsProvider adds a type that implements IClaimsProvider
func AddTransientIClaimsProvider(builder *di.Builder, implType reflect.Type, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	di.AddTransientWithImplementedTypes(builder, implType, implementedTypes...)
}

// AddTransientIClaimsProviderWithMetadata adds a type that implements IClaimsProvider
func AddTransientIClaimsProviderWithMetadata(builder *di.Builder, implType reflect.Type, metaData map[string]interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	di.AddTransientWithImplementedTypesWithMetadata(builder, implType, metaData, implementedTypes...)
}

// AddTransientIClaimsProviderByFunc adds a type by a custom func
func AddTransientIClaimsProviderByFunc(builder *di.Builder, implType reflect.Type, build func(ctn di.Container) (interface{}, error), implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	di.AddTransientWithImplementedTypesByFunc(builder, implType, build, implementedTypes...)
}

// AddTransientIClaimsProviderByFuncWithMetadata adds a type by a custom func
func AddTransientIClaimsProviderByFuncWithMetadata(builder *di.Builder, implType reflect.Type, build func(ctn di.Container) (interface{}, error), metaData map[string]interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	di.AddTransientWithImplementedTypesByFuncWithMetadata(builder, implType, build, metaData, implementedTypes...)
}

// AddScopedIClaimsProvider adds a type that implements IClaimsProvider
func AddScopedIClaimsProvider(builder *di.Builder, implType reflect.Type, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	di.AddScopedWithImplementedTypes(builder, implType, implementedTypes...)
}

// AddScopedIClaimsProviderWithMetadata adds a type that implements IClaimsProvider
func AddScopedIClaimsProviderWithMetadata(builder *di.Builder, implType reflect.Type, metaData map[string]interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	di.AddScopedWithImplementedTypesWithMetadata(builder, implType, metaData, implementedTypes...)
}

// AddScopedIClaimsProviderByFunc adds a type by a custom func
func AddScopedIClaimsProviderByFunc(builder *di.Builder, implType reflect.Type, build func(ctn di.Container) (interface{}, error), implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	di.AddScopedWithImplementedTypesByFunc(builder, implType, build, implementedTypes...)
}

// AddScopedIClaimsProviderByFuncWithMetadata adds a type by a custom func
func AddScopedIClaimsProviderByFuncWithMetadata(builder *di.Builder, implType reflect.Type, build func(ctn di.Container) (interface{}, error), metaData map[string]interface{}, implementedTypes ...reflect.Type) {
	implementedTypes = append(implementedTypes, ReflectTypeIClaimsProvider)
	di.AddScopedWithImplementedTypesByFuncWithMetadata(builder, implType, build, metaData, implementedTypes...)
}

// RemoveAllIClaimsProvider removes all IClaimsProvider from the DI
func RemoveAllIClaimsProvider(builder *di.Builder) {
	builder.RemoveAllByType(ReflectTypeIClaimsProvider)
}

// GetIClaimsProviderFromContainer alternative to SafeGetIClaimsProviderFromContainer but panics of object is not present
func GetIClaimsProviderFromContainer(ctn di.Container) IClaimsProvider {
	return ctn.GetByType(ReflectTypeIClaimsProvider).(IClaimsProvider)
}

// GetManyIClaimsProviderFromContainer alternative to SafeGetManyIClaimsProviderFromContainer but panics of object is not present
func GetManyIClaimsProviderFromContainer(ctn di.Container) []IClaimsProvider {
	objs := ctn.GetManyByType(ReflectTypeIClaimsProvider)
	var results []IClaimsProvider
	for _, obj := range objs {
		results = append(results, obj.(IClaimsProvider))
	}
	return results
}

// SafeGetIClaimsProviderFromContainer trys to get the object by type, will not panic, returns nil and error
func SafeGetIClaimsProviderFromContainer(ctn di.Container) (IClaimsProvider, error) {
	obj, err := ctn.SafeGetByType(ReflectTypeIClaimsProvider)
	if err != nil {
		return nil, err
	}
	return obj.(IClaimsProvider), nil
}

// GetIClaimsProviderDefinition returns that last definition registered that this container can provide
func GetIClaimsProviderDefinition(ctn di.Container) *di.Def {
	def := ctn.GetDefinitionByType(ReflectTypeIClaimsProvider)
	return def
}

// GetIClaimsProviderDefinitions returns all definitions that this container can provide
func GetIClaimsProviderDefinitions(ctn di.Container) []*di.Def {
	defs := ctn.GetDefinitionsByType(ReflectTypeIClaimsProvider)
	return defs
}

// SafeGetManyIClaimsProviderFromContainer trys to get the object by type, will not panic, returns nil and error
func SafeGetManyIClaimsProviderFromContainer(ctn di.Container) ([]IClaimsProvider, error) {
	objs, err := ctn.SafeGetManyByType(ReflectTypeIClaimsProvider)
	if err != nil {
		return nil, err
	}
	var results []IClaimsProvider
	for _, obj := range objs {
		results = append(results, obj.(IClaimsProvider))
	}
	return results, nil
}
