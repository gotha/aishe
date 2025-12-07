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

        this.redisClient = redisClient;
        this.baseUrl = baseUrl ?? AISHE_API_URL;
        this.timeout = timeout ?? REQUEST_TIMEOUT_MS;
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

        const redisClient: RedisClientType = createClient({
            socket: {
                host: redisHost ?? REDIS_HOST,
                port: redisPort ?? REDIS_PORT,
            },
            username: redisUsername ?? REDIS_USERNAME,
            password: redisPassword ?? REDIS_PASSWORD,
            database: redisDatabase ?? REDIS_DATABASE,
        });
        await redisClient.connect();

        // Healthcheck Redis
        const health = await redisClient.ping();
        if (!health) {
            throw new ServerNotReachableError("Redis is not reachable. Hint: check your Redis configuration.");
        }

        return new AIsheHTTPClient(redisClient, baseUrl, timeout);
    }

    /**
     * Close AIshe API client
     */
    async close(): Promise<void> {
        // TODO: implement this function
        // [NEW] 1. quit the Redis client
        await this.redisClient.quit();
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
        if (!question || !question.trim()) {
            throw new Error("Question cannot be empty");
        }

        // First, check if the question is cached in Redis
        const cacheKey = generateCacheKey(question);
        const cachedAnswer = await this.redisClient.get(cacheKey);
        if (cachedAnswer) {
            console.log(`INFO: Cache HIT for question: ${question}`);
            return JSON.parse(cachedAnswer) as AnswerResponse;
        }

        // Fetch from AIshe API
        console.log(`INFO: Cache MISS for question: ${question}`);
        const endpoint = `${this.baseUrl}/api/v1/ask`;
        const body = { question: question.trim() };
        const response = await aisheAPIRequest("POST", endpoint, this.timeout, body);
        const answerResponse = response as AnswerResponse;

        // Cache the answer
        await this.redisClient.set(cacheKey, JSON.stringify(answerResponse));
        console.log(`INFO: Cached answer for question: ${question}`);

        return answerResponse;
    }

    /**
     * Check if the question is cached in Redis
     * 
     * NOTE: use this method to determine the source of your answer
     *       (Redis cache HIT or AIshe API)
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
        const cacheKey = generateCacheKey(question);
        return (await this.redisClient.exists(cacheKey)) > 0;
    }
}
