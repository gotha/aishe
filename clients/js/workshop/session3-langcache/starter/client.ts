import { LangCache } from "@redis-ai/langcache";

import {
    AISHE_API_URL,
    REQUEST_TIMEOUT_MS,
    APIClientError,
    aisheAPIRequest,
    type HealthResponse,
    type AnswerResponse,
    LANGCACHE_STRICT_SIMILARITY_THRESHOLD,
    LANGCACHE_CLOSE_SIMILARITY_THRESHOLD,
    LANGCACHE_LOOSE_SIMILARITY_THRESHOLD,
} from "aishe-client";

// TODO: replace with your own API KEY, CACHE ID and URL
const LANGCACHE_API_KEY: string = "YOUR_API_KEY";
const LANGCACHE_CACHE_ID: string = "YOUR_CACHE_ID";
const LANGCACHE_SERVER_URL: string = "YOUR_SERVER_URL";

/** AIshe API client */
export class AIsheHTTPClient {
    /** Base URL of the AIshe API server */
    private readonly baseUrl: string;
    /** Request timeout in milliseconds */
    private readonly timeout: number;
    /** LangCache client */
    private readonly langCache: LangCache;
    /** Similarity threshold */
    private readonly similarityThreshold: number;

    /**
     * Initialize AIshe API client
     *
     * @param baseUrl - Base URL of the AIshe API server (default: AISHE_API_URL).
     * @param timeout - Request timeout in milliseconds (default: REQUEST_TIMEOUT_MS).
     */
    private constructor(baseUrl?: string, timeout?: number, similarityThreshold?: number) {
        // TODO: implement this function
        // 1. if baseURL is not provided, use default AISHE_API_URL from 'aishe-client'
        // 2. if timeout is not provided, use default REQUEST_TIMEOUT_MS from 'aishe-client'
        // ~~3. assign redisClient to property this.redisClient~~
        // [NEW] 3. if similarityThreshold is not provided, use one of the defaults:
        //           - LANGCACHE_STRICT_SIMILARITY_THRESHOLD from 'aishe-client'
        //           - LANGCACHE_CLOSE_SIMILARITY_THRESHOLD from 'aishe-client'
        //           - LANGCACHE_LOOSE_SIMILARITY_THRESHOLD from 'aishe-client'
        // [NEW] 4. create a new LangCache client

        throw new Error("AIsheHTTPClient.constructor: Not implemented");
    }

    static async create(baseUrl?: string, timeout?: number, similarityThreshold?: number): Promise<AIsheHTTPClient> {
        // TODO: implement this function
        // ~~ 1. create a Redis client using createClient()
        // ~~      NOTE: use default values from 'aishe-client' for Redis configuration parameters
        // ~~ 2. connect to Redis
        // ~~ 3. add a basic healthcheck by pinging Redis
        // ~~      NOTE: if Redis is not reachable, throw ServerNotReachableError
        // [NEW] 1. create a new AIshe HTTP client & implement constructor()

        throw new Error("AIsheHTTPClient.create: Not implemented");
    }

    /**
     * Close AIshe API client
     */
    async close(): Promise<void> {
        // TODO: implement this function
        // ~~1. quit the Redis client~~
        // [NEW] 1. no operation needed; clean up this method
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
        // ~~2. Generate a cache key using the question~~
        // ~~3. Check if the question is cached in Redis~~
        // [NEW] 2. Check if the question is cached in LangCache
        //       - if found (cache HIT), return the cached answer
        //       - otherwise proceed to step 4
        // 4. Build the ask endpoint: baseURL + "/api/v1/ask"
        // 5. Make a POST request using aisheAPIRequest()
        //    NOTE: aisheAPIRequest() will handle the HTTP request, timeout, error handling for you.
        // ~~6. Cache raw answer in Redis (as a string)~~
        // [NEW] 6. Cache raw answer in LangCache
        // 7. Decode JSON response into AnswerResponse
        // 8. Return the answer response
        //
        // NOTE: LangCache search response is in the format:
        //
        // {
        //   "data": [
        //     {
        //       "id": "...",
        //       "prompt": "What is JavaScript?",
        //       "response": "JavaScript is a dynamic programming language...",
        //       "attributes": { ... },
        //       "similarity": 0.95,
        //       "searchStrategy": "semantic"
        //     },
        //   ],
        //
        throw new Error("AIsheHTTPClient.askQuestion: Not implemented");
    }

    /**
     * Check if the question is cached in LangCache
     *
     * NOTE: use this method to determine the source of your answer
     *       (LangCache cache HIT or AIshe API)
     *
     * @param question - Question to check.
     *
     * @returns True if the question is cached (cache HIT), false otherwise (cache MISS).
     */
    async isCached(question: string): Promise<boolean> {
        // TODO: implement this function
        // ~~1. Generate a cache key using the question~~
        // ~~2. Check if the question is cached in Redis~~
        // [NEW] 1. Check if the question is cached in LangCache
        //    - if found (cache HIT), return true
        //    - otherwise, return false
        throw new Error("AIsheHTTPClient.isCached: Not implemented");
    }
}
