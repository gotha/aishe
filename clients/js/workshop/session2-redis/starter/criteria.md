# Session 2 - Criteria for Success

If you succeeded in time with session 1, copy your implementation of session 1
(`main.ts` & `client.ts`) inside _session2-redis/starter/_.

Otherwise, check out the reference implementation inside _session1-basic/solution/_
and use it to complete session 2.

## Criteria

`main.ts`:

- [x] Create a basic AIshe HTTP client
- [x] Check AIshe's health
- [x] Print health response
- [x] Ask AIshe a question
- [x] Display AIshe's results
- [x] Time asking AIshe a question
- [ ] Check if cache key is in Redis before asking a question
- [ ] Display the source of your retrieved answer (Redis cache HIT or AIshe API)

`client.ts`

- [x] implement AIsheHTTPClient constructor
- [x] implement checkHealth() method
- [x] implement askQuestion() method
- [ ] modify constructor() to create a Redis client
- [ ] implement create() static method + add a basic healthcheck by pinging Redis
- [ ] implement close()
- [ ] modify askQuestion() to check the cache before fetching from AIshe API
    - [ ] if answer is cahed, retrieve cached value
    - [ ] otherwise, generate a cache key, fetch from AIshe API, cache & return your result

## Redis Client Basics

```ts
import { createClient, type RedisClientType } from "redis";

const redisClient: RedisClientType = createClient({
    url: "redis://localhost:6379",
});
await redisClient.connect();

try {
    await redisClient.set("framework", "Redis 8.4");
    const framework = await redisClient.get("framework");
    console.log("framework =", framework);
} finally {
    await redisClient.quit();
}
```

## References

AIshe client models, types & errors are available from the `aishe-client` library.
You may import them via `import {...} from 'aishe-client';`

Available requests, models, types, errors, configs:

### Default configuration values

```ts
/**
 * AIshe API server host
 */
export const SERVER_HOST: string = process.env.SERVER_HOST || "localhost";

/**
 * AIshe API server port
 */
export const SERVER_PORT: number = parseInt(process.env.SERVER_PORT || "8000");

/**
 * AIshe API URL
 */
export const AISHE_API_URL: string = process.env.AISHE_API_URL || `http://${SERVER_HOST}:${SERVER_PORT}`;

/**
 * Request timeout in milliseconds
 */
export const REQUEST_TIMEOUT_MS: number = parseInt(process.env.REQUEST_TIMEOUT || "120") * 1000;

/**
 * Redis host
 */
export const REDIS_HOST: string = process.env.REDIS_HOST || "localhost";

/**
 * Redis port
 */
export const REDIS_PORT: number = parseInt(process.env.REDIS_PORT || "6379");

/**
 * Redis database
 */
export const REDIS_DATABASE: number = parseInt(process.env.REDIS_DATABASE || "0");

/**
 * Redis username
 */
export const REDIS_USERNAME: string = process.env.REDIS_USERNAME || "default";

/**
 * Redis password
 */
export const REDIS_PASSWORD: string = process.env.REDIS_PASSWORD || "";

/**
 * Redis URL
 */
export const REDIS_URL: string =
    process.env.REDIS_URL ||
    `redis://${REDIS_USERNAME}:${REDIS_PASSWORD}@${REDIS_HOST}:${REDIS_PORT}/${REDIS_DATABASE}`;
```

### Health and Answer response models from AIshe

```ts
/**
 * Request model for asking a question
 */
export interface QuestionRequest {
    /** The question to answer */
    question: string;
}

/**
 * Model for source reference
 */
export interface Source {
    /** Source reference number */
    number: number;

    /** Source article title */
    title: string;

    /** URL of the source article */
    url: string;
}

/**
 * Response model for an answer
 */
export interface AnswerResponse {
    /** Generated answer to the given question */
    answer: string;

    /** List of sources used to generate the answer */
    sources: Source[];

    /** Time taken to process the question in seconds */
    processing_time: number;
}

/**
 * Response model for health check
 */
export interface HealthResponse {
    /** Server status */
    status: string;

    /** Whether Ollama service is accessible */
    ollama_accessible: boolean;

    /** Additional status message */
    message?: string;
}
```

### Errors thrown by AIshe

```ts
/**
 * Base exception for API client errors
 */
export class APIClientError extends Error {
    /**
     * Constructor for APIClientError
     * @param message - The error message
     */
    constructor(message: string) {
        super(message);
        this.name = "APIClientError";
        Object.setPrototypeOf(this, APIClientError.prototype);
    }
}

/**
 * Exception raised when the server is not reachable
 */
export class ServerNotReachableError extends APIClientError {
    /**
     * Constructor for ServerNotReachableError
     * @param message - The error message
     */
    constructor(message: string) {
        super(message);
        this.name = "ServerNotReachableError";
        Object.setPrototypeOf(this, ServerNotReachableError.prototype);
    }
}

/**
 * Exception raised when the server returns an error
 */
export class ServerError extends APIClientError {
    /**
     * Constructor for ServerError
     * @param message - The error message
     */
    constructor(message: string) {
        super(message);
        this.name = "ServerError";
        Object.setPrototypeOf(this, ServerError.prototype);
    }
}
```

### Function for HTTP requests to AIshe

This function implements the HTTP specifics, timeout, error handling.

Focus more on working with AIshe and caching. Focus less on HTTP specifics.

```ts
import { REQUEST_TIMEOUT_MS } from "./config.js";
import { ServerError, APIClientError, ServerNotReachableError } from "./errors.js";

/**
 * Make a request to AIshe API server
 *
 * NOTE: this function is intentionally decoupled from the AIshe API client
 *       in ./client.ts for the purposes of the workshop. We want to hide the
 *       HTTP specifics & error handling away from attendees so they can focus
 *       on interacting with AIshe API server + caching the responses.
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
export async function aisheAPIRequest(
    method: string,
    endpoint: string,
    timeout?: number,
    body?: unknown,
): Promise<unknown> {
    timeout = timeout ?? REQUEST_TIMEOUT_MS;
    if (timeout <= 0) {
        throw new Error("Timeout must be greater than 0");
    }

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeout);

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
    } catch (error) {
        if (error instanceof Error && error.name === "AbortError") {
            throw new ServerNotReachableError(`${method} request timed out after ${timeout}ms`);
        } else if (error instanceof Error && error.name === "ServerError") {
            throw error;
        }
        throw new APIClientError(`${method} request failed! Unexpected error: ${error}`);
    } finally {
        clearTimeout(timeoutId);
    }
}
```

Usage:

```ts
import { aisheAPIRequest } from 'aishe-client';

const method = ...
const endpoint = ...
const timeout = ...
const body = ...
const response = await aisheAPIRequest(method, endpoint, timeout, body);
const typedResponse = response as ResponseType
```

### Function to generate cache keys

This function generates a cache key for you which look like:
`aishe:question:2990b8f25d9f7a585798544a7231ffcec5f0ef7507691f077cf70ba889af83ee`
`aishe:question:32b0ee199b6b95c4301c495665aad49af074b92d8a86f86d00b5c5535b360049`

```ts
/**
 * Generate a cache key for a question
 *
 * Uses SHA-256 hash to generate a unique key.
 *
 * @param question - Question to generate a cache key for.
 *
 * @returns Cache key.
 */
export function generateCacheKey(question: string): string {
    const hash = crypto.createHash("sha256").update(question).digest("hex");
    return `${REDIS_CACHE_KEY_PREFIX}:${hash}`;
}
```

## HELP: constructor() doesn't allow `await`

JavaScript's constructor is not an asynchronous method, which is why you can't
call and await other async methods in it.

You'll still need to use asynchronous methods from redis-js client library to
connect and interact with Redis. That's why you'll end up using a factory creation
design pattern. It goes like this:

- Make your constructor private and only use it to initalize properties
- Create `static async create()` method which creates a Redis client + a new instance
  of AIsheHTTPClient
- Connect to Redis & check it's health in your `create()` factory method
- Use your `create()` method to get an instance of AIsheHTTPClient instead of `new`

Example:

```ts
// client.ts
import { createClient, RedisClientType } from 'redis';

export class AIsheHTTPClient {
    private readonly redisClient;

    private constructor(redisClient: RedisClientType) {
        this.redisClient = redisClient;
    }

    static async create(...): Promise<AIsheHTTPClient> {
        const redisClient: RedisClientType = createClient({ ... });
        // connect to Redis via your client
        // PING Redis for a basic healthcheck
        return new AIsheHTTPClient(redisClient);
    }
}

// main.ts
import { AIsheHTTPClient } from './client.js';

const client = await AIsheHTTPClient.create(...);
```
