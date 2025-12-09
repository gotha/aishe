import { LangCache } from "@redis-ai/langcache";

import {
    AISHE_API_URL,
    REQUEST_TIMEOUT_MS,
    APIClientError,
    aisheAPIRequest,
    type HealthResponse,
    type AnswerResponse,
    LANGCACHE_CLOSE_SIMILARITY_THRESHOLD,
    LANGCACHE_API_KEY,
    LANGCACHE_CACHE_ID,
    LANGCACHE_SERVER_URL,
} from "aishe-client";

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

        if (!LANGCACHE_API_KEY || !LANGCACHE_CACHE_ID || !LANGCACHE_SERVER_URL) {
            throw new Error("LangCache API key, cache ID, and server URL are not set");
        }

        this.baseUrl = baseUrl ?? AISHE_API_URL;
        this.timeout = timeout ?? REQUEST_TIMEOUT_MS;
        this.similarityThreshold = similarityThreshold ?? LANGCACHE_CLOSE_SIMILARITY_THRESHOLD;
        this.langCache = new LangCache({
            apiKey: LANGCACHE_API_KEY,
            cacheId: LANGCACHE_CACHE_ID,
            serverURL: LANGCACHE_SERVER_URL,
        });
    }

    static async create(baseUrl?: string, timeout?: number, similarityThreshold?: number): Promise<AIsheHTTPClient> {
        // TODO: implement this function
        // ~~ 1. create a Redis client using createClient()
        // ~~      NOTE: use default values from 'aishe-client' for Redis configuration parameters
        // ~~ 2. connect to Redis
        // ~~ 3. add a basic healthcheck by pinging Redis
        // ~~      NOTE: if Redis is not reachable, throw ServerNotReachableError
        // [NEW] 1. create a new AIshe HTTP client & implement constructor()

        return new AIsheHTTPClient(baseUrl, timeout, similarityThreshold);
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
        const endpoint = `${this.baseUrl}/health`;
        const response = await aisheAPIRequest("GET", endpoint, this.timeout);
        const healthResponse = response as HealthResponse;
        if (healthResponse.status !== "healthy") {
            throw new APIClientError("Health check failed");
        }
        return healthResponse;
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
        if (!question || !question.trim()) {
            throw new Error("Question cannot be empty");
        }

        // First, check if the question is cached in LangCache
        const searchResponse = await this.langCache.search({
            prompt: question,
            similarityThreshold: this.similarityThreshold,
        });
        if (searchResponse.data.length > 0 && searchResponse.data[0]?.response) {
            console.log(`INFO: Cache HIT for question: ${question}`);
            console.log(`INFO: Similarity: ${searchResponse.data[0].similarity}`);
            return JSON.parse(searchResponse.data[0].response) as AnswerResponse;
        }

        // Fetch from AIshe API
        console.log(`INFO: Cache MISS for question: ${question}`);
        const endpoint = `${this.baseUrl}/api/v1/ask`;
        const body = { question: question.trim() };
        const response = await aisheAPIRequest("POST", endpoint, this.timeout, body);
        const answerResponse = response as AnswerResponse;

        // Cache the answer
        await this.langCache.set({
            prompt: question,
            response: JSON.stringify(answerResponse),
        });
        console.log(`INFO: Cached answer for question: ${question}`);

        return answerResponse;
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
        const searchResponse = await this.langCache.search({ prompt: question });
        return searchResponse.data.length > 0;
    }
}
