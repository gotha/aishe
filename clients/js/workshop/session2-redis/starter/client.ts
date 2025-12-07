import { createClient, type RedisClientType } from "redis";

import {
    AISHE_API_URL,
    REQUEST_TIMEOUT_MS,
    REDIS_HOST,
    REDIS_PORT,
    REDIS_DATABASE,
    REDIS_USERNAME,
    REDIS_PASSWORD,
    APIClientError,
    aisheAPIRequest,
    generateCacheKey,
    type HealthResponse,
    type AnswerResponse,
    ServerNotReachableError,
} from "aishe-client";

/** AIshe API client */
export class AIsheHTTPClient {
    /** Base URL of the AIshe API server */
    private readonly baseUrl: string;
    /** Request timeout in milliseconds */
    private readonly timeout: number;
    /** Redis client */
    private readonly redisClient: RedisClientType;

    /**
     * Initialize AIshe API client
     *
     * @param redisClient - Redis client.
     * @param baseUrl - Base URL of the AIshe API server (default: AISHE_API_URL).
     * @param timeout - Request timeout in milliseconds (default: REQUEST_TIMEOUT_MS).
     */
    private constructor(redisClient: RedisClientType, baseUrl?: string, timeout?: number) {
        // TODO: implement this function
        // 1. if baseURL is not provided, use default AISHE_API_URL from 'aishe-client'
        // 2. if timeout is not provided, use default REQUEST_TIMEOUT_MS from 'aishe-client'
        // [NEW] 3. assign redisClient to property this.redisClient

        throw new Error("AIsheHTTPClient.constructor: Not implemented");
    }

    static async create(
        baseUrl?: string,
        timeout?: number,
        redisHost?: string,
        redisPort?: number,
        redisUsername?: string,
        redisPassword?: string,
        redisDatabase?: number,
    ): Promise<AIsheHTTPClient> {
        // TODO: implement this function
        // [NEW] 1. create a Redis client using createClient()
        //       NOTE: use default values from 'aishe-client' for Redis configuration parameters
        // [NEW] 2. connect to Redis
        // [NEW] 3. add a basic healthcheck by pinging Redis
        //       NOTE: if Redis is not reachable, throw ServerNotReachableError
        throw new Error("AIsheHTTPClient.create: Not implemented");
    }

    /**
     * Close AIshe API client
     */
    async close(): Promise<void> {
        // TODO: implement this function
        // [NEW] 1. quit the Redis client
        throw new Error("AIsheHTTPClient.close: Not implemented");
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
        // NOTE: Use your implementation from session 1 OR
        //       reference implementation in `session1-basic/solution/client.ts`
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
        // [NEW] 2. Generate a cache key using the question
        // [NEW] 3. Check if the question is cached in Redis
        //       - if found (cache HIT), return the cached answer
        //       - otherwise, proceed to step 4
        // 4. Build the ask endpoint: baseURL + "/api/v1/ask"
        // 5. Make a POST request using aisheAPIRequest()
        //    NOTE: aisheAPIRequest() will handle the HTTP request, timeout, error handling for you.
        // [NEW] 6. Cache raw answer in Redis (as a string)
        // 7. Decode JSON response into AnswerResponse
        // 8. Return the answer response
        //
        // NOTE: You should save the raw answer in Redis as a string and use
        //       the JSON module to parse it back & forward.
        //       If you have some extra time, you can try implementing caching
        //       with native Redis JSON data types.
        //
        // NOTE: Use your implementation from session 1 OR
        //       reference implementation in `session1-basic/solution/client.ts`
        throw new Error("AIsheHTTPClient.askQuestion: Not implemented");
    }

    /**
     * Check if the question is cached in Redis
     *
     * NOTE: use this method to determine the source of your answer
     * (Redis cache HIT or AIshe API)
     * 
     * @param question - Question to check.
     *
     * @returns True if the question is cached (cache HIT), false otherwise (cache MISS).
     */
    async isCached(question: string): Promise<boolean> {
        // TODO: implement this function
        // 1. Generate a cache key using the question
        // 2. Check if the question is cached in Redis
        //    - if found (cache HIT), return true
        //    - otherwise, return false
        throw new Error("AIsheHTTPClient.isCached: Not implemented");
    }
}
