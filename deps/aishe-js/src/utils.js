import crypto from "node:crypto";
import { REQUEST_TIMEOUT_MS, REDIS_CACHE_KEY_PREFIX } from "./config.js";
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
export async function aisheAPIRequest(method, endpoint, timeout, body) {
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
    }
    catch (error) {
        if (error instanceof Error && error.name === "AbortError") {
            throw new ServerNotReachableError(`${method} request timed out after ${timeout}ms`);
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
 * Generate a cache key for a question
 *
 * Uses SHA-256 hash to generate a unique key.
 *
 * @param question - Question to generate a cache key for.
 *
 * @returns Cache key.
 */
export function generateCacheKey(question) {
    const hash = crypto.createHash("sha256").update(question).digest("hex");
    return `${REDIS_CACHE_KEY_PREFIX}:${hash}`;
}
//# sourceMappingURL=utils.js.map