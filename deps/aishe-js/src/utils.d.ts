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
export declare function aisheAPIRequest(method: string, endpoint: string, timeout?: number, body?: unknown): Promise<unknown>;
/**
 * Generate a cache key for a question
 *
 * Uses SHA-256 hash to generate a unique key.
 *
 * @param question - Question to generate a cache key for.
 *
 * @returns Cache key.
 */
export declare function generateCacheKey(question: string): string;
//# sourceMappingURL=utils.d.ts.map