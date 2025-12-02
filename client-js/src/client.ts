import { AISHE_API_URL, REQUEST_TIMEOUT_MS } from "./config.js";
import { ServerError, APIClientError, ServerNotReachableError } from "./errors.js";
import type { HealthResponse, AnswerResponse } from "./models.js";

export class RAGAPIClient {
    private baseUrl: string;
    private timeout: number;  // in milliseconds

    constructor(baseUrl: string = AISHE_API_URL, timeout: number = REQUEST_TIMEOUT_MS) {
        this.baseUrl = baseUrl;
        this.timeout = timeout;  // in milliseconds
    }

    async checkHealth(): Promise<HealthResponse> {
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), this.timeout);

        try {
            const response = await fetch(`${this.baseUrl}/health`, {
                method: "GET",
                signal: controller.signal,
            });
            if (!response.ok) {
                throw new ServerError(`Health check failed! status: ${response.status}`);
            }

            const data = await response.json() 
            if (typeof data !== "object" || data === null || typeof data.status !== "string") {
                throw new APIClientError("Malformed health response from server");
            }
            return data as HealthResponse;
        } catch (error) {
            if (error instanceof Error && error.name === "AbortError") {
                throw new ServerNotReachableError(`Request timed out after ${this.timeout}ms`);
            } 
            throw new APIClientError(`Unexpected error: ${error}`);
        } finally {
            clearTimeout(timeoutId);
        }
    }

    async askQuestion(question: string): Promise<AnswerResponse> {
        if (!question || !question.trim()) {
            throw new Error("Question cannot be empty");
        }

        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), this.timeout);

        try {
            const response = await fetch(`${this.baseUrl}/api/v1/ask`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ question: question.trim() }),
                signal: controller.signal,
            });
            if (!response.ok) {
                throw new ServerError(`Question answering failed! status: ${response.status}`);
            }
            
            const data = await response.json();
            if (typeof data !== "object" || data === null || typeof data.answer !== "string" || !Array.isArray(data.sources) || typeof data.processing_time !== "number") {
                throw new APIClientError("Malformed answer response from server");
            }
            return data as AnswerResponse;
        } catch (error) {
            if (error instanceof Error && error.name === "AbortError") {
                throw new ServerNotReachableError(`Request timed out after ${this.timeout}ms`);
            } 
            throw new APIClientError(`Unexpected error: ${error}`);
        } finally {
            clearTimeout(timeoutId);
        }
    }
}
