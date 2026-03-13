package hyperui

import (
	"strings"

	"github.com/assaidy/hyper/v2"
)

type ButtonVariant uint

const (
	VariantDefault ButtonVariant = iota
	VariantPrimary
	VariantSecondary
)

var variantClasses = map[ButtonVariant]string{
	VariantDefault:   "bg-gray-100 text-gray-900 hover:bg-gray-200 focus:ring-gray-400",
	VariantPrimary:   "bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500",
	VariantSecondary: "bg-gray-800 text-white hover:bg-gray-700 focus:ring-gray-500",
}

type ButtonSize uint

const (
	SizeDefault ButtonSize = iota
	SizeSmall
	SizeLarge
)

var sizeClasses = map[ButtonSize]string{
	SizeDefault: "px-4 py-2 text-sm",
	SizeSmall:   "px-3 py-1.5 text-xs",
	SizeLarge:   "px-6 py-3 text-base",
}

type ButtonParams struct {
	Variant     ButtonVariant
	Size        ButtonSize
	IsFullWidth bool
	Attributes  []h.Attribute
}

// Example:
//
//	// simple
//	Button()("Click Me!")
//
//	// with params
//	Button(ButtonParams{
//		Variant:     VariantPrimary,
//		IsFullWidth: true,
//		Attributes:  []Attribute{
//			AttrID("btn"),
//			AttrClass("some-class"),
//		},
//	})(
//		"Click Me!"
//	)
func Button(params ...ButtonParams) h.ElementBuilder {
	var p ButtonParams
	if len(params) != 0 {
		p = params[0]
	}

	return func(children ...any) h.Element {
		element := h.BUTTON(p.Attributes...)(children...)
		mergeStyles(&element, getStylesClass(p))
		return element
	}
}

func getStylesClass(params ButtonParams) string {
	var stylesBuilder strings.Builder

	stylesBuilder.WriteString("inline-flex items-center justify-center rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2")
	stylesBuilder.WriteByte(' ')

	if class, ok := variantClasses[params.Variant]; ok {
		stylesBuilder.WriteString(class)
		stylesBuilder.WriteByte(' ')
	} else {
		panic("invalid button variant")
	}

	if class, ok := sizeClasses[params.Size]; ok {
		stylesBuilder.WriteString(class)
		stylesBuilder.WriteByte(' ')
	} else {
		panic("invalid button size")
	}

	if params.IsFullWidth {
		stylesBuilder.WriteString("w-full")
	}

	return stylesBuilder.String()
}
