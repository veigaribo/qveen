package params

// This file has been generated by Qveen from
// `qveen/params_errors.toml`.
// Please do not modify it directly.

import (
	"fmt"
	"github.com/veigaribo/qveen/utils"
	"github.com/veigaribo/qveen/prompts"
	"strings"
)

type ParamError struct {
	Path   []any
	Reason string
}

func MakeParamError(path []any, reason string) ParamError {
	return ParamError{
		Path:   path,
		Reason: reason,
	}
}

func (e ParamError) Error() string {
	var builder strings.Builder

	builder.WriteRune('`')
	utils.WritePathString(e.Path, &builder)
	builder.WriteRune('`')
	builder.WriteRune(' ')
	builder.WriteString(e.Reason)

	return builder.String()
}

// Specific errors for use with `errors.As`.

type MetaWrongTypeError struct {
	Err ParamError
}

func MakeMetaWrongTypeError(path []any) MetaWrongTypeError {
	return MetaWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain a table.",
		),
	}
}

func (e MetaWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaRootTemplateWrongTypeError struct {
	Err ParamError
}

func MakeMetaRootTemplateWrongTypeError(path []any) MetaRootTemplateWrongTypeError {
	return MetaRootTemplateWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain a string.",
		),
	}
}

func (e MetaRootTemplateWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaRootTemplateWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaRootOutputWrongTypeError struct {
	Err ParamError
}

func MakeMetaRootOutputWrongTypeError(path []any) MetaRootOutputWrongTypeError {
	return MetaRootOutputWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain a string.",
		),
	}
}

func (e MetaRootOutputWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaRootOutputWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaPairsWrongTypeError struct {
	Err ParamError
}

func MakeMetaPairsWrongTypeError(path []any) MetaPairsWrongTypeError {
	return MetaPairsWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain an array.",
		),
	}
}

func (e MetaPairsWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaPairsWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaPairWrongTypeError struct {
	Err ParamError
}

func MakeMetaPairWrongTypeError(path []any) MetaPairWrongTypeError {
	return MetaPairWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain a table.",
		),
	}
}

func (e MetaPairWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaPairWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaRootTemplateMissingInMultipleError struct {
	Err ParamError
}

func MakeMetaRootTemplateMissingInMultipleError(path []any) MetaRootTemplateMissingInMultipleError {
	return MetaRootTemplateMissingInMultipleError{
		Err: MakeParamError(
			path,
			"required field is required for multiple files but is missing.",
		),
	}
}

func (e MetaRootTemplateMissingInMultipleError) Error() string {
	return e.Err.Error()
}

func (e MetaRootTemplateMissingInMultipleError) Unwrap() error {
	return e.Err
}

type MetaRootOutputMissingInMultipleError struct {
	Err ParamError
}

func MakeMetaRootOutputMissingInMultipleError(path []any) MetaRootOutputMissingInMultipleError {
	return MetaRootOutputMissingInMultipleError{
		Err: MakeParamError(
			path,
			"required field is required for multiple files but is missing.",
		),
	}
}

func (e MetaRootOutputMissingInMultipleError) Error() string {
	return e.Err.Error()
}

func (e MetaRootOutputMissingInMultipleError) Unwrap() error {
	return e.Err
}

type MetaPairTemplateMissingError struct {
	Err ParamError
}

func MakeMetaPairTemplateMissingError(path []any) MetaPairTemplateMissingError {
	return MetaPairTemplateMissingError{
		Err: MakeParamError(
			path,
			"missing required field.",
		),
	}
}

func (e MetaPairTemplateMissingError) Error() string {
	return e.Err.Error()
}

func (e MetaPairTemplateMissingError) Unwrap() error {
	return e.Err
}

type MetaPairTemplateWrongTypeError struct {
	Err ParamError
}

func MakeMetaPairTemplateWrongTypeError(path []any) MetaPairTemplateWrongTypeError {
	return MetaPairTemplateWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain a string.",
		),
	}
}

func (e MetaPairTemplateWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaPairTemplateWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaPairOutputMissingError struct {
	Err ParamError
}

func MakeMetaPairOutputMissingError(path []any) MetaPairOutputMissingError {
	return MetaPairOutputMissingError{
		Err: MakeParamError(
			path,
			"missing required field.",
		),
	}
}

func (e MetaPairOutputMissingError) Error() string {
	return e.Err.Error()
}

func (e MetaPairOutputMissingError) Unwrap() error {
	return e.Err
}

type MetaPairOutputWrongTypeError struct {
	Err ParamError
}

func MakeMetaPairOutputWrongTypeError(path []any) MetaPairOutputWrongTypeError {
	return MetaPairOutputWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain a string.",
		),
	}
}

func (e MetaPairOutputWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaPairOutputWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaPromptsWrongTypeError struct {
	Err ParamError
}

func MakeMetaPromptsWrongTypeError(path []any) MetaPromptsWrongTypeError {
	return MetaPromptsWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain an array.",
		),
	}
}

func (e MetaPromptsWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaPromptsWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaPromptWrongTypeError struct {
	Err ParamError
}

func MakeMetaPromptWrongTypeError(path []any) MetaPromptWrongTypeError {
	return MetaPromptWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain a table.",
		),
	}
}

func (e MetaPromptWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaPromptWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaPromptNameMissingError struct {
	Err ParamError
}

func MakeMetaPromptNameMissingError(path []any) MetaPromptNameMissingError {
	return MetaPromptNameMissingError{
		Err: MakeParamError(
			path,
			"missing required field.",
		),
	}
}

func (e MetaPromptNameMissingError) Error() string {
	return e.Err.Error()
}

func (e MetaPromptNameMissingError) Unwrap() error {
	return e.Err
}

type MetaPromptNameWrongTypeError struct {
	Err ParamError
}

func MakeMetaPromptNameWrongTypeError(path []any) MetaPromptNameWrongTypeError {
	return MetaPromptNameWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain a string.",
		),
	}
}

func (e MetaPromptNameWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaPromptNameWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaPromptKindWrongTypeError struct {
	Err ParamError
}

func MakeMetaPromptKindWrongTypeError(path []any) MetaPromptKindWrongTypeError {
	return MetaPromptKindWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain a string.",
		),
	}
}

func (e MetaPromptKindWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaPromptKindWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaPromptKindInvalidError struct {
	Err ParamError
}

func MakeMetaPromptKindInvalidError(path []any) MetaPromptKindInvalidError {
	return MetaPromptKindInvalidError{
		Err: MakeParamError(
			path,
			fmt.Sprintf("field does not contain one of the allowed values: %v.", prompts.SupportedPromptKinds),
		),
	}
}

func (e MetaPromptKindInvalidError) Error() string {
	return e.Err.Error()
}

func (e MetaPromptKindInvalidError) Unwrap() error {
	return e.Err
}

type MetaPromptTitleWrongTypeError struct {
	Err ParamError
}

func MakeMetaPromptTitleWrongTypeError(path []any) MetaPromptTitleWrongTypeError {
	return MetaPromptTitleWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain a string.",
		),
	}
}

func (e MetaPromptTitleWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaPromptTitleWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaPromptOptionsMissingError struct {
	Err ParamError
}

func MakeMetaPromptOptionsMissingError(path []any) MetaPromptOptionsMissingError {
	return MetaPromptOptionsMissingError{
		Err: MakeParamError(
			path,
			"required field for `select` is missing.",
		),
	}
}

func (e MetaPromptOptionsMissingError) Error() string {
	return e.Err.Error()
}

func (e MetaPromptOptionsMissingError) Unwrap() error {
	return e.Err
}

type MetaPromptOptionsWrongTypeError struct {
	Err ParamError
}

func MakeMetaPromptOptionsWrongTypeError(path []any) MetaPromptOptionsWrongTypeError {
	return MetaPromptOptionsWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain an array.",
		),
	}
}

func (e MetaPromptOptionsWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaPromptOptionsWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaPromptOptionTitleMissingError struct {
	Err ParamError
}

func MakeMetaPromptOptionTitleMissingError(path []any) MetaPromptOptionTitleMissingError {
	return MetaPromptOptionTitleMissingError{
		Err: MakeParamError(
			path,
			"missing required field.",
		),
	}
}

func (e MetaPromptOptionTitleMissingError) Error() string {
	return e.Err.Error()
}

func (e MetaPromptOptionTitleMissingError) Unwrap() error {
	return e.Err
}

type MetaPromptOptionTitleWrongTypeError struct {
	Err ParamError
}

func MakeMetaPromptOptionTitleWrongTypeError(path []any) MetaPromptOptionTitleWrongTypeError {
	return MetaPromptOptionTitleWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain a string.",
		),
	}
}

func (e MetaPromptOptionTitleWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaPromptOptionTitleWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaPromptOptionWrongTypeError struct {
	Err ParamError
}

func MakeMetaPromptOptionWrongTypeError(path []any) MetaPromptOptionWrongTypeError {
	return MetaPromptOptionWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but contains neither a string nor a table.",
		),
	}
}

func (e MetaPromptOptionWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaPromptOptionWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaLeftDelimWrongTypeError struct {
	Err ParamError
}

func MakeMetaLeftDelimWrongTypeError(path []any) MetaLeftDelimWrongTypeError {
	return MetaLeftDelimWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain a string.",
		),
	}
}

func (e MetaLeftDelimWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaLeftDelimWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaRightDelimWrongTypeError struct {
	Err ParamError
}

func MakeMetaRightDelimWrongTypeError(path []any) MetaRightDelimWrongTypeError {
	return MetaRightDelimWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain a string.",
		),
	}
}

func (e MetaRightDelimWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaRightDelimWrongTypeError) Unwrap() error {
	return e.Err
}

type MetaCaseWrongTypeError struct {
	Err ParamError
}

func MakeMetaCaseWrongTypeError(path []any) MetaCaseWrongTypeError {
	return MetaCaseWrongTypeError{
		Err: MakeParamError(
			path,
			"field present but does not contain a string.",
		),
	}
}

func (e MetaCaseWrongTypeError) Error() string {
	return e.Err.Error()
}

func (e MetaCaseWrongTypeError) Unwrap() error {
	return e.Err
}
