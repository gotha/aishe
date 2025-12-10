import {
    AISHE_API_URL,
    REQUEST_TIMEOUT_MS,
    APIClientError,
    aisheAPIRequest,
    type HealthResponse,
    type AnswerResponse,
} from "aishe-client";

/** AIshe API client */
export class AIsheHTTPClient {
    /** Base URL of the AIshe API server */
    private readonly baseUrl: string;
    /** Request timeout in milliseconds */
    private readonly timeout: number;

    /**
     * Initialize AIshe API client
     *
     * @param baseUrl - Base URL of the AIshe API server (default: AISHE_API_URL).
     * @param timeout - Request timeout in milliseconds (default: REQUEST_TIMEOUT_MS).
     */
    constructor(baseUrl?: string, timeout?: number) {
        // TODO: implement this function
        // 1. if baseURL is not provided, use default AISHE_API_URL from 'aishe-client'
        // 2. if timeout is not provided, use default REQUEST_TIMEOUT_MS from 'aishe-client'
        throw new Error("AIsheHTTPClient.constructor: Not implemented");
    }

    /**
     * Close AIshe API client
     */
    close(): void {
        // No-operation: HTTP client doesn't need explicit closing.
        // method includeded for consistency with other clients.
    }

    /**
     * Check AIshe's health
     *
     * Health request is:
     *    GET /health
     *
     * Health response is in the format:
     *
     * {
     *   "status": "ok",
     *   "ollama_accessible": true,
     *   "message": "Additional status message"
     * }
     *
     * @throws APIClientError - If the health check failed.
     *
     * @returns Health response from AIshe server.
     */
    async checkHealth(): Promise<HealthResponse> {
        // TODO: implement this function
        // 1. Build the health endpoint: baseURL + "/health"
        // 2. Make a GET request using aisheAPIRequest()
        //    NOTE: aisheAPIRequest() will handle the HTTP request, timeout, error handling for you.
        // 3. Decode JSON response into HealthResponse
        // 4. Check the response status code (should be "healthy");
        //    if not, throw APIClientError with an appropriate error message
        // 5. Return the health response
        //
        // NOTE: aisheAPIRequest() method signature:
        //    async aisheAPIRequest(method: "GET" | "POST", endpoint: string, timeout?, body?): Promise<unknown>
        //
        throw new Error("AIsheHTTPClient.checkHealth: Not implemented");
    }

    /**
     * Ask AIshe a question
     *
     * Question request is:
     *    POST /api/v1/ask
     *    Body: { "question": "What is JavaScript?" }
     *
     * Answer response is in the format:
     *
     * {
     *   "answer": "JavaScript is a dynamic programming language...",
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
     *
     * @returns Answer from AIshe server.
     */
    async askQuestion(question: string): Promise<AnswerResponse> {
        // TODO: implement this function
        // 1. Check if the question is empty; if so, throw Error with an appropriate error message
        // 2. Build the ask endpoint: baseURL + "/api/v1/ask"
        // 3. Make a POST request using aisheAPIRequest()
        //    NOTE: aisheAPIRequest() will handle the HTTP request, timeout, error handling for you.
        // 4. Decode JSON response into AnswerResponse
        // 5. Return the answer response
        //
        // NOTE: aisheAPIRequest() method signature:
        //    async aisheAPIRequest(method: "GET" | "POST", endpoint: string, timeout?, body?): Promise<unknown>
        //
        throw new Error("AIsheHTTPClient.askQuestion: Not implemented");
    }
}
