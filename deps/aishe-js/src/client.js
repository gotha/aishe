import { AISHE_API_URL, REQUEST_TIMEOUT_MS } from "./config.js";
import { ServerError, APIClientError, ServerNotReachableError } from "./errors.js";
/** AIshe API client */
export class RAGAPIClient {
    /**
     * Initialize AIshe API client
     *
     * @param baseUrl - Base URL of the AIshe API server (default: AISHE_API_URL).
     * @param timeout - Request timeout in milliseconds (default: REQUEST_TIMEOUT_MS).
     */
    constructor(
    /** Base URL of the AIshe API server */
    baseUrl = AISHE_API_URL, 
    /** Request timeout in milliseconds */
    timeout = REQUEST_TIMEOUT_MS) {
        this.baseUrl = baseUrl;
        this.timeout = timeout;
    }
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
    async checkHealth() {
        const endpoint = this.healthEndpoint();
        const method = "GET";
        const data = await this.aisheRequest(method, endpoint);
        const isValid = this.isHealthResponse(data);
        if (!isValid) {
            throw new APIClientError("Malformed health response from server");
        }
        return data;
    }
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
    async askQuestion(question) {
        if (!question || !question.trim()) {
            throw new Error("Question cannot be empty");
        }
        const endpoint = this.askEndpoint();
        const method = "POST";
        const body = { question: question.trim() };
        const data = await this.aisheRequest(method, endpoint, body);
        const isValid = this.isAnswerResponse(data);
        if (!isValid) {
            throw new APIClientError("Malformed answer response from server");
        }
        return data;
    }
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
    async aisheRequest(method, endpoint, body) {
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), this.timeout);
        try {
            const response = await fetch(endpoint, {
                method: method,
                headers: {
                    "Content-Type": "application/json",
                },
                body: body ? JSON.stringify(body) : null,
                signal: controller.signal,
            });
            if (!response.ok) {
                throw new ServerError(`${method} request failed! Status: ${response.status}`);
            }
            return await response.json();
        }
        catch (error) {
            if (error instanceof Error && error.name === "AbortError") {
                throw new ServerNotReachableError(`${method} request timed out after ${this.timeout}ms`);
            }
            else if (error instanceof Error && error.name === "ServerError") {
                throw error;
            }
            throw new APIClientError(`${method} request failed! Unexpected error: ${error}`);
        }
        finally {
            clearTimeout(timeoutId);
        }
    }
    /**
     * Get the endpoint for checking the health of AIshe server
     *
     * NOTE: Expect endpoint to be `/health` or similar.
     *
     * @returns Endpoint for checking the health of AIshe server.
     */
    healthEndpoint() {
        return `${this.baseUrl}/health`;
    }
    /**
     * Get the endpoint for asking a question to AIshe server
     *
     * NOTE: Expect endpoint to be `/api/v1/ask` or similar.
     *
     * @returns Endpoint for asking a question to AIshe server.
     */
    askEndpoint() {
        return `${this.baseUrl}/api/v1/ask`;
    }
    /**
     * Validate health response from AIshe server
     *
     * @param data - Health response from AIshe server.
     *
     * @returns True if the health response is valid, false otherwise.
     */
    isHealthResponse(data) {
        return (typeof data === "object" &&
            data !== null &&
            typeof data.status === "string" &&
            typeof data.ollama_accessible === "boolean");
    }
    /**
     * Validate answer response from AIshe server
     *
     * @param data - Answer response from AIshe server.
     *
     * @returns True if the answer response is valid, false otherwise.
     */
    isAnswerResponse(data) {
        return (typeof data === "object" &&
            data !== null &&
            typeof data.answer === "string" &&
            Array.isArray(data.sources) &&
            typeof data.processing_time === "number");
    }
}
//# sourceMappingURL=client.js.map