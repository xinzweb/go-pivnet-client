package service

import (
	"fmt"
	"github.com/baotingfang/go-pivnet-client/config"
	. "github.com/baotingfang/go-pivnet-client/utils"
	"github.com/baotingfang/go-pivnet-client/vlog"
	"regexp"
	"strings"
)

type Validator interface {
	Validate() bool
	GetErrorMessages() []string
}

type AbstractValidator struct {
	errorMessages []string
}

func (v AbstractValidator) GetErrorMessages() []string {
	return v.errorMessages
}

type MetaDataValidator struct {
	AbstractValidator
	metadata config.Metadata
}

func NewMetaDataValidator(metadata config.Metadata) *MetaDataValidator {
	return &MetaDataValidator{metadata: metadata}
}

func (mv *MetaDataValidator) Validate() bool {
	var messages []string

	r := mv.metadata.Release

	vlog.Debug("validating date settings in metadata file...")
	if !Empty(r.EndOfAvailabilityDate) && !Empty(r.EndOfAvailabilityDateOffset) {
		messages = append(messages,
			"can not specify both end_of_availability_date and end_of_availability_date_offset")
	}
	if !Empty(r.EndOfSupportDate) && !IsDate(r.EndOfSupportDate) {
		messages = append(messages,
			`end_of_support_date must be a valid date of the format "YYYY-MM-DD"`)
	}

	if !Empty(r.EndOfGuidanceDate) && !IsDate(r.EndOfGuidanceDate) {
		messages = append(messages,
			`end_of_guidance_date must be a valid date of the format "YYYY-MM-DD"`)
	}

	if !Empty(r.EndOfAvailabilityDate) && !IsDate(r.EndOfAvailabilityDate) {
		messages = append(messages,
			`end_of_availability_date must be a valid date of the format "YYYY-MM-DD"`)
	}

	vlog.Debug("validating offset settings in metadata file...")
	if !Empty(r.EndOfAvailabilityDateOffset) &&
		!NewOffsetValidator(r.EndOfAvailabilityDateOffset).Validate() {
		messages = append(messages,
			`end_of_availability_date_offset must be a valid offset of the form "(+\d+[mdyMDY])+"`)
	}

	vlog.Debug("validating file groups in metadata file...")
	fileGroups := mv.metadata.FileGroups
	for _, fileGroup := range fileGroups {
		vlog.Debug("\tvalidating file group: %s", fileGroup.Name)
		for _, productFile := range fileGroup.ProductFiles {
			vlog.Debug("\t\tvalidating product file:%s", productFile.Name)
			pfValidator := NewProductFileValidator(productFile)
			if !pfValidator.Validate() {
				messages = append(messages, pfValidator.GetErrorMessages()...)
			}
		}
	}

	vlog.Debug("validating product files in metadata file...")
	productFiles := mv.metadata.ProductFiles
	for _, productFile := range productFiles {
		vlog.Debug("\tvalidating product file:%s", productFile.Name)
		pfValidator := NewProductFileValidator(productFile)
		if !pfValidator.Validate() {
			messages = append(messages, pfValidator.GetErrorMessages()...)
		}
	}

	mv.errorMessages = messages
	return len(mv.errorMessages) == 0
}

type OffsetValidator struct {
	AbstractValidator
	value string
}

func NewOffsetValidator(value string) *OffsetValidator {
	return &OffsetValidator{value: value}
}

func (ov *OffsetValidator) Validate() bool {
	var messages []string

	offsetRegexp, _ := regexp.Compile(`(\+\d+[mdyMDY])+`)
	if !offsetRegexp.Match([]byte(ov.value)) {
		messages = append(messages,
			fmt.Sprintf(`"%s" is not a valid offset value`, ov.value))
	}

	ov.errorMessages = messages
	return len(ov.errorMessages) == 0
}

type RequiredValidator struct {
	AbstractValidator
	values []interface{}
}

func NewRequiredValidator(values ...interface{}) *RequiredValidator {
	return &RequiredValidator{values: values}
}

func (rv *RequiredValidator) Validate() bool {
	var messages []string

	for index, value := range rv.values {
		switch v := value.(type) {
		case string:
			val := value.(string)
			if Empty(val) {
				var valuesStr []string
				for _, v := range rv.values {
					valuesStr = append(valuesStr, fmt.Sprintf("%v", v))
				}
				messages = append(messages,
					fmt.Sprintf("value is empty, index=%d (| %v |)",
						index, strings.Join(valuesStr, " | ")))
			}
		default:
			messages = append(messages,
				fmt.Sprintf("RequiredValidator only support string type, and not support this type: %T", v))
		}
	}

	rv.errorMessages = messages
	return len(rv.errorMessages) == 0
}

type ProductFileValidator struct {
	AbstractValidator
	pf config.ProductFile
}

func NewProductFileValidator(productFile config.ProductFile) ProductFileValidator {
	return ProductFileValidator{pf: productFile}
}

func (pv *ProductFileValidator) Validate() bool {
	var messages []string

	productFile := pv.pf
	requiredValidator := NewRequiredValidator(productFile.File, productFile.UploadAs, productFile.FileType, productFile.FileVersion)

	if !requiredValidator.Validate() {
		messages = append(messages,
			fmt.Sprintf("One of settings is empty.[file, upload_as, file_type, file_version]"),
		)
		messages = append(messages, requiredValidator.GetErrorMessages()...)
	}

	pv.errorMessages = messages
	return len(pv.errorMessages) == 0
}
