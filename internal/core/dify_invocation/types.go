package dify_invocation

import (
	"github.com/go-playground/validator/v10"
	"github.com/langgenius/dify-plugin-daemon/internal/types/entities/app_entities"
	"github.com/langgenius/dify-plugin-daemon/internal/types/entities/plugin_entities"
	"github.com/langgenius/dify-plugin-daemon/internal/types/entities/requests"
	"github.com/langgenius/dify-plugin-daemon/internal/types/validators"
)

type BaseInvokeDifyRequest struct {
	TenantId string     `json:"tenant_id"`
	UserId   string     `json:"user_id"`
	Type     InvokeType `json:"type"`
}

type InvokeType string

const (
	INVOKE_TYPE_LLM                      InvokeType = "llm"
	INVOKE_TYPE_TEXT_EMBEDDING           InvokeType = "text_embedding"
	INVOKE_TYPE_RERANK                   InvokeType = "rerank"
	INVOKE_TYPE_TTS                      InvokeType = "tts"
	INVOKE_TYPE_SPEECH2TEXT              InvokeType = "speech2text"
	INVOKE_TYPE_MODERATION               InvokeType = "moderation"
	INVOKE_TYPE_TOOL                     InvokeType = "tool"
	INVOKE_TYPE_NODE_PARAMETER_EXTRACTOR InvokeType = "node_parameter_extractor"
	INVOKE_TYPE_NODE_QUESTION_CLASSIFIER InvokeType = "node_question_classifier"
	INVOKE_TYPE_APP                      InvokeType = "app"
	INVOKE_TYPE_STORAGE                  InvokeType = "storage"
	INVOKE_TYPE_ENCRYPT                  InvokeType = "encrypt"
)

type InvokeLLMRequest struct {
	BaseInvokeDifyRequest
	requests.BaseRequestInvokeModel
	requests.InvokeLLMSchema
}

type InvokeTextEmbeddingRequest struct {
	BaseInvokeDifyRequest
	requests.BaseRequestInvokeModel
	requests.InvokeTextEmbeddingSchema
}

type InvokeRerankRequest struct {
	BaseInvokeDifyRequest
	requests.BaseRequestInvokeModel
	requests.InvokeRerankSchema
}

type InvokeTTSRequest struct {
	BaseInvokeDifyRequest
	requests.BaseRequestInvokeModel
	requests.InvokeTTSSchema
}

type InvokeSpeech2TextRequest struct {
	BaseInvokeDifyRequest
	requests.BaseRequestInvokeModel
	requests.InvokeSpeech2TextSchema
}

type InvokeModerationRequest struct {
	BaseInvokeDifyRequest
	requests.BaseRequestInvokeModel
	requests.InvokeModerationSchema
}

type InvokeAppSchema struct {
	AppId          string                  `json:"app_id" validate:"required"`
	Inputs         map[string]any          `json:"inputs" validate:"omitempty"`
	Query          string                  `json:"query" validate:"omitempty"`
	ResponseMode   string                  `json:"response_mode"`
	ConversationId string                  `json:"conversation_id"`
	User           string                  `json:"user" validate:"omitempty"`
	Files          []*app_entities.FileVar `json:"files" validate:"omitempty,dive"`
}

type StorageOpt string

const (
	STORAGE_OPT_GET StorageOpt = "get"
	STORAGE_OPT_SET StorageOpt = "set"
	STORAGE_OPT_DEL StorageOpt = "del"
)

func isStorageOpt(fl validator.FieldLevel) bool {
	opt := StorageOpt(fl.Field().String())
	return opt == STORAGE_OPT_GET || opt == STORAGE_OPT_SET || opt == STORAGE_OPT_DEL
}

func init() {
	validators.GlobalEntitiesValidator.RegisterValidation("storage_opt", isStorageOpt)
}

type InvokeStorageRequest struct {
	Opt   StorageOpt `json:"opt" validate:"required,storage_opt"`
	Key   string     `json:"key" validate:"required"`
	Value string     `json:"value"` // encoded in hex, optional
}

type InvokeAppRequest struct {
	BaseInvokeDifyRequest

	InvokeAppSchema
}

type ModelConfig struct {
	Provider         string         `json:"provider" validate:"required"`
	Name             string         `json:"name" validate:"required"`
	Mode             string         `json:"mode" validate:"required"`
	CompletionParams map[string]any `json:"completion_params" validate:"omitempty"`
}

type InvokeParameterExtractorRequest struct {
	BaseInvokeDifyRequest

	Parameters []struct {
		Name        string   `json:"name" validate:"required"`
		Type        string   `json:"type" validate:"required,oneof=string number bool select array[string] array[number] array[object]"`
		Options     []string `json:"options" validate:"omitempty"`
		Description string   `json:"description" validate:"omitempty"`
		Required    bool     `json:"required" validate:"omitempty"`
	} `json:"parameters" validate:"required,dive"`

	Model       ModelConfig `json:"model" validate:"required"`
	Instruction string      `json:"instruction" validate:"omitempty"`
	Query       string      `json:"query" validate:"required"`
}

type InvokeQuestionClassifierRequest struct {
	BaseInvokeDifyRequest

	Classes []struct {
		ID   string `json:"id" validate:"required"`
		Name string `json:"name" validate:"required"`
	} `json:"classes" validate:"required,dive"`

	Model       ModelConfig `json:"model" validate:"required"`
	Instruction string      `json:"instruction" validate:"omitempty"`
	Query       string      `json:"query" validate:"required"`
}

type EncryptOpt string

const (
	ENCRYPT_OPT_ENCRYPT EncryptOpt = "encrypt"
	ENCRYPT_OPT_DECRYPT EncryptOpt = "decrypt"
	ENCRYPT_OPT_CLEAR   EncryptOpt = "clear"
)

func isEncryptOpt(fl validator.FieldLevel) bool {
	opt := EncryptOpt(fl.Field().String())
	return opt == ENCRYPT_OPT_ENCRYPT || opt == ENCRYPT_OPT_DECRYPT || opt == ENCRYPT_OPT_CLEAR
}

type EncryptNamespace string

const (
	ENCRYPT_NAMESPACE_ENDPOINT EncryptNamespace = "endpoint"
)

func isEncryptNamespace(fl validator.FieldLevel) bool {
	opt := EncryptNamespace(fl.Field().String())
	return opt == ENCRYPT_NAMESPACE_ENDPOINT
}

func init() {
	validators.GlobalEntitiesValidator.RegisterValidation("encrypt_opt", isEncryptOpt)
	validators.GlobalEntitiesValidator.RegisterValidation("encrypt_namespace", isEncryptNamespace)
}

type InvokeEncryptSchema struct {
	Opt       EncryptOpt                                `json:"opt" validate:"required,encrypt_opt"`
	Namespace EncryptNamespace                          `json:"namespace" validate:"required,encrypt_namespace"`
	Identity  string                                    `json:"identity" validate:"required"`
	Data      map[string]any                            `json:"data" validate:"omitempty"`
	Config    map[string]plugin_entities.ProviderConfig `json:"config" validate:"omitempty,dive"`
}

type InvokeEncryptRequest struct {
	BaseInvokeDifyRequest

	InvokeEncryptSchema
}

func (r *InvokeEncryptRequest) EncryptRequired(settings map[string]any) bool {
	if r.Config == nil {
		return false
	}

	// filter out which key needs encrypt
	for k := range settings {
		if config, ok := r.Config[k]; ok {
			if config.Type == plugin_entities.CONFIG_TYPE_SECRET_INPUT {
				return true
			}
		}
	}

	return false
}

type InvokeToolRequest struct {
	BaseInvokeDifyRequest
	ToolType requests.ToolType `json:"tool_type" validate:"required,tool_type"`
	requests.InvokeToolSchema
}

type InvokeNodeResponse struct {
	ProcessData map[string]any `json:"process_data" validate:"required"`
	Outputs     map[string]any `json:"outputs" validate:"required"`
	Inputs      map[string]any `json:"inputs" validate:"required"`
}

type InvokeEncryptionResponse struct {
	Error string         `json:"error"`
	Data  map[string]any `json:"data"`
}
