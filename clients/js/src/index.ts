// Public API surface for the AIshe JavaScript/TypeScript client.

// Main client
export { RAGAPIClient } from "./client.js";

// Types (interfaces)
export type { AnswerResponse, HealthResponse, ErrorResponse, Source, QuestionRequest } from "./models.js";

// Error classes
export { APIClientError, ServerError, ServerNotReachableError } from "./errors.js";

// Configuration display function
export { AISHE_API_URL, REQUEST_TIMEOUT_MS, displayConfig } from "./config.js";

// Command-line interface
export { AIsheCLI } from "./cli.js";

// Request function
export { aisheAPIRequest } from "./requests.js";
