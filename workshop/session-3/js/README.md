# Session 3 - Criteria for Success

If you succeeded in time with session 2, copy your implementation of session 2
(`main.ts` & `client.ts`) inside _session3-langcache/starter/_.

Otherwise, check out the reference implementation inside _session2-redis/solution/_
and use it to complete session 3.

## Criteria

`main.ts`:

- [x] Create a basic AIshe HTTP client
- [x] Check AIshe's health
- [x] Print health response
- [x] Ask AIshe a question
- [x] Display AIshe's results
- [x] Time asking AIshe a question
- [ ] Check if cache key is in ~~Redis~~ LangCache before asking a question
- [ ] Display the source of your retrieved answer ~~(Redis cache HIT or AIshe API)~~ (LangCache or AIshe API)

`client.ts`

- [x] implement AIsheHTTPClient constructor
- [x] implement checkHealth() method
- [x] implement askQuestion() method
- [x] modify constructor() to create a Redis client
- [x] implement create() static method + add a basic healthcheck by pinging Redis
- [x] implement close()
- [x] implement isCached()
- [x] implement generateCacheKey()
- [x] modify askQuestion() to check the cache before fetching from AIshe API
    - [x] if answer is cahed, retrieve cached value
    - [x] otherwise, generate a cache key, fetch from AIshe API, cache & return your result
- [ ] modify constructor() / create() to create a new LangCache client
- [ ] clean up close() (no operation needed)
- [ ] modify isCached() to check LangCache instead of Redis
- [ ] modify askQuestion() to use LangCache instead of a basic Redis Cache
      NOTE: remember to use your API KEY, Cache ID and URL list

## Prerequisites

- [Node](https://nodejs.org/en) v25.2.1+
- [npm](https://www.npmjs.com/) v11.6.2+
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/) (for AIshe server, Redis, Ollama)
- Running AISHE stack on http://localhost:8000 or another port
- Familiarity with JavaScript / TypeScript
- Basic understanding of HTPP and REST APIs

## Setup

```bash
cd workshop/session-3/js/
npm install
export LANGCACHE_API_KEY="<YOUR_API_KEY>"
export LANGCACHE_CACHE_ID="<YOUR_CACHE_ID>"
export LANGCACHE_SERVER_URL="<LANGCACHE_SERVER_URL>"
```

## Run

Starter:

```bash
cd workshop/session-3/js/
npm run session
```

Reference solution:

```bash
cd workshop/session-3/js/
npm run solution
```

## LangCache Basics

Redis LangCache is a fully-managed semantic caching service. It's used to cache
large language model (LLM) responses based on meaning. This allows you to reuse
the same response for similar prompts, massively speeding up your app & cutting
LLM costs (+ saving GPU/TPU lives).

### Create a LangCache account

- Visit [https://redis.io/langcache/](https://redis.io/langcache/) and click "Try it for free".
- Create an account with GitHub, Google account, etc.
- IMPORTANT: Copy the API KEY and save it somewhere; You'll only see it once.
- Quick create a new LangCache service
- Find the Cache ID and copy it somewhere
- Copy one of the URLs somewhere

### How to use LangCache

Basic API examples: [https://redis.io/docs/latest/develop/ai/langcache/api-examples/](https://redis.io/docs/latest/develop/ai/langcache/api-examples/)
Theoretical intro to langcache: [https://redis.io/docs/latest/develop/ai/langcache/](https://redis.io/docs/latest/develop/ai/langcache/)

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

/**
 * LangCache strict similarity threshold
 */
export const LANGCACHE_STRICT_SIMILARITY_THRESHOLD: number = parseFloat(
    process.env.LANGCACHE_STRICT_SIMILARITY_THRESHOLD || "0.95",
);

/**
 * LangCache close similarity threshold
 */
export const LANGCACHE_CLOSE_SIMILARITY_THRESHOLD: number = parseFloat(
    process.env.LANGCACHE_CLOSE_SIMILARITY_THRESHOLD || "0.9",
);

/**
 * LangCache loose similarity threshold
 */
export const LANGCACHE_LOOSE_SIMILARITY_THRESHOLD: number = parseFloat(
    process.env.LANGCACHE_LOOSE_SIMILARITY_THRESHOLD || "0.8",
);

/**
 * LangCache API key
 */
export const LANGCACHE_API_KEY: string = process.env.LANGCACHE_API_KEY || "YOUR_API_KEY";

/**
 * LangCache cache ID
 */
export const LANGCACHE_CACHE_ID: string = process.env.LANGCACHE_CACHE_ID || "YOUR_CACHE_ID";

/**
 * LangCache server URL
 */
export const LANGCACHE_SERVER_URL: string = process.env.LANGCACHE_SERVER_URL || "YOUR_SERVER_URL";
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
