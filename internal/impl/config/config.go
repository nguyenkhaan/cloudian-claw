package config

import "os"

func OPENROUTER_API_KEY() string {
	return os.Getenv("OPENROUTER_API_KEY")
}

func GATEWAY_API_TOKEN() string {
	return os.Getenv("GATEWAY_API_TOKEN")
}

func SERVER_HOST() string {
	return os.Getenv("SERVER_HOST")
}

func SERVER_PORT() string {
	return os.Getenv("SERVER_PORT")
}

func POSTGRES_URL() string {
	return os.Getenv("POSTGRES_URL")
}

func DEFAULT_PROVIDER() string {
	return os.Getenv("DEFAULT_PROVIDER")
}

func DEFAULT_MODEL() string {
	return os.Getenv("DEFAULT_MODEL")
}

func AI_TEMPERATURE() string {
	return os.Getenv("AI_TEMPERATURE")
}

func AI_MAX_TOKENS() string {
	return os.Getenv("AI_MAX_TOKENS")
}

func QUOTA_MAX_REQUESTS_PER_DAY() string {
	return os.Getenv("QUOTA_MAX_REQUESTS_PER_DAY")
}

func GLOBAL_TOOLS_ENABLED() string {
	return os.Getenv("GLOBAL_TOOLS_ENABLED")
}

func MAX_TOOL_CALLS() string {
	return os.Getenv("MAX_TOOL_CALLS")
}

func MAX_PROVIDER_RETRIES() string {
	return os.Getenv("MAX_PROVIDER_RETRIES")
}

func MAX_EXECUTION_DURATION() string {
	return os.Getenv("MAX_EXECUTION_DURATION")
}

func MAX_CONTEXT_TOKENS() string {
	return os.Getenv("MAX_CONTEXT_TOKENS")
}

func MAX_CONTINUATION_DEPTH() string {
	return os.Getenv("MAX_CONTINUATION_DEPTH")
}