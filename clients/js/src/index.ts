// Public API surface for the AIshe JavaScript/TypeScript client.

// Main client
export { RAGAPIClient } from "./client.js";

// Types (interfaces)
export type { AnswerResponse, HealthResponse, ErrorResponse, Source, QuestionRequest } from "./models.js";

// Error classes
export { APIClientError, ServerError, ServerNotReachableError } from "./errors.js";

// Configuration display function
export {
    AISHE_API_URL,
    REQUEST_TIMEOUT_MS,
    REDIS_HOST,
    REDIS_PORT,
    REDIS_DATABASE,
    REDIS_USERNAME,
    REDIS_PASSWORD,
    REDIS_URL,
    REDIS_CACHE_KEY_PREFIX,
    LANGCACHE_STRICT_SIMILARITY_THRESHOLD,
    LANGCACHE_CLOSE_SIMILARITY_THRESHOLD,
    LANGCACHE_LOOSE_SIMILARITY_THRESHOLD,
    displayConfig,
} from "./config.js";

// Command-line interface
export { AIsheCLI } from "./cli.js";

// Request function
export { aisheAPIRequest, generateCacheKey } from "./utils.js";
