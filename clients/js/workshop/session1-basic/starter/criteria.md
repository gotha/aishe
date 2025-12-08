# Session 1 - Criteria for Success

## Criteria

`main.ts`:

- [ ] Create a basic AIshe HTTP client
- [ ] Check AIshe's health
- [ ] Print health response
- [ ] Ask AIshe a question
- [ ] Display AIshe's results
- [ ] Time asking AIshe a question

`client.ts`

- [ ] implement AIsheHTTPClient constructor
- [ ] implement checkHealth() method
- [ ] implement askQuestion() method

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
