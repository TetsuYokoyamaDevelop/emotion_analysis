package model

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			FunctionCall struct {
				Arguments string `json:"arguments"`
			} `json:"function_call"`
		} `json:"message"`
	} `json:"choices"`
}
