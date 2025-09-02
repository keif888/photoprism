package vision

import (
	"errors"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/pkg/media"
)

// CaptionPromptDefault is the default prompt used to generate captions.
var CaptionPromptDefault = "Create a caption with exactly one sentence in the active voice that describes the main visual content." +
	" Begin with the main subject and clear action. Avoid text formatting, meta-language, and filler words."

// CaptionModelDefault specifies the default model used to generate captions,
// see https://ollama.com/search?c=vision for a list of available models.
var CaptionModelDefault = "gemma3"

// Caption returns generated captions for the specified images.
func Caption(images Files, mediaSrc media.Src) (result *CaptionResult, model *Model, err error) {
	// Return if there is no configuration or no image classification models are configured.
	if Config == nil {
		return result, model, errors.New("vision service is not configured")
	} else if model = Config.Model(ModelTypeCaption); model != nil {
		// Use remote service API if a server endpoint has been configured.
		if uri, method := model.Endpoint(); uri != "" && method != "" {
			var apiRequest *ApiRequest
			var apiResponse *ApiResponse

			if apiRequest, err = NewApiRequest(model.EndpointRequestFormat(), images, model.EndpointFileScheme()); err != nil {
				return result, model, err
			}

			switch model.Service.RequestFormat {
			case ApiFormatOllama:
				apiRequest.Model, _, _ = model.Model()
			default:
				_, apiRequest.Model, apiRequest.Version = model.Model()
			}

			// Set system prompt if configured.
			apiRequest.System = model.GetSystemPrompt()

			// Set caption prompt if configured.
			apiRequest.Prompt = model.GetPrompt()

			// Set caption model request options.
			apiRequest.Options = model.GetOptions()

			// Log JSON request data in trace mode.
			apiRequest.WriteLog()

			// Todo: Refactor response handling to support different API response formats,
			//       including those used by Ollama and OpenAI.
			if apiResponse, err = PerformApiRequest(apiRequest, uri, method, model.EndpointKey()); err != nil {
				return result, model, err
			} else if apiResponse.Result.Caption == nil {
				return result, model, errors.New("invalid caption model response")
			}

			// Set image as the default caption source.
			if apiResponse.Result.Caption.Text != "" && apiResponse.Result.Caption.Source == "" {
				apiResponse.Result.Caption.Source = entity.SrcImage
			}

			result = apiResponse.Result.Caption
		} else {
			return result, model, errors.New("invalid caption model configuration")
		}
	} else {
		return result, model, errors.New("missing caption model")
	}

	return result, model, nil
}
