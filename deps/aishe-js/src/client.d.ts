import type { HealthResponse, AnswerResponse } from "./models.js";
/** AIshe API client */
export declare class RAGAPIClient {
    /** Base URL of the AIshe API server */
    private readonly baseUrl;
    /** Request timeout in milliseconds */
    private readonly timeout;
    /**
     * Initialize AIshe API client
     *
     * @param baseUrl - Base URL of the AIshe API server (default: AISHE_API_URL).
     * @param timeout - Request timeout in milliseconds (default: REQUEST_TIMEOUT_MS).
     */
    constructor(
    /** Base URL of the AIshe API server */
    baseUrl?: string, 
    /** Request timeout in milliseconds */
    timeout?: number);
    /**
     * Check AIshe's health
     *
     * Health response is in the format:
     *
     * {
     *   "status": "ok",
     *   "ollama_accessible": true,
     *   "message": "Additional status message"
     * }
     *
     * @throws ServerNotReachableError - If the request timed out.
     * @throws ServerError - If the request failed.
     * @throws APIClientError - If the request is malformed.
     *
     * @returns Health response from AIshe server.
     */
    checkHealth(): Promise<HealthResponse>;
    /**
     * Ask AIshe a question
     *
     * Answer response is in the format:
     *
     * {
     *   "answer": "The answer to the question",
     *   "sources": [
     *     {
     *       "number": 1,
     *       "title": "Source title",
     *       "url": "Source URL"
     *     },
     *     ... // more sources
     *   ],
     *   "processing_time": 0.123 // in seconds
     * }
     *
     * @param question - Question to ask.
     *
     * @throws Error - If the question is empty.
     * @throws ServerNotReachableError - If the request timed out.
     * @throws ServerError - If the request failed.
     * @throws APIClientError - If the request or response is malformed.
     *
     * @returns Answer from AIshe server.
     */
    askQuestion(question: string): Promise<AnswerResponse>;
    /**
     * Make a request to AIshe API server
     *
     * IMPORTANT: Please do not supply a body for GET requests.
     *
     * @param method - HTTP method to use (GET, POST, etc.).
     * @param endpoint - Endpoint to ship a request to.
     * @param body - Body of the request (optional). NOTE: skip for GET requests.
     *
     * @throws ServerNotReachableError - If the request timed out.
     * @throws ServerError - If the request failed.
     * @throws APIClientError - If the request is malformed.
     *
     * @returns Response from AIshe server.
     */
    private aisheRequest;
    /**
     * Get the endpoint for checking the health of AIshe server
     *
     * NOTE: Expect endpoint to be `/health` or similar.
     *
     * @returns Endpoint for checking the health of AIshe server.
     */
    private healthEndpoint;
    /**
     * Get the endpoint for asking a question to AIshe server
     *
     * NOTE: Expect endpoint to be `/api/v1/ask` or similar.
     *
     * @returns Endpoint for asking a question to AIshe server.
     */
    private askEndpoint;
    /**
     * Validate health response from AIshe server
     *
     * @param data - Health response from AIshe server.
     *
     * @returns True if the health response is valid, false otherwise.
     */
    private isHealthResponse;
    /**
     * Validate answer response from AIshe server
     *
     * @param data - Answer response from AIshe server.
     *
     * @returns True if the answer response is valid, false otherwise.
     */
    private isAnswerResponse;
}
//# sourceMappingURL=client.d.ts.map