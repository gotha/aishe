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
 * Request timeout in seconds
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
 * Redis cache key prefix
 */
export const REDIS_CACHE_KEY_PREFIX: string = process.env.REDIS_CACHE_KEY_PREFIX || "aishe:question";

/**
 * Display the current configuration
 */
export function displayConfig(): void {
    console.log("Configuration:");
    console.log(`  Server Host: ${SERVER_HOST}`);
    console.log(`  Server Port: ${SERVER_PORT}`);
    console.log(`  API URL: ${AISHE_API_URL}`);
    console.log(`  Request Timeout: ${REQUEST_TIMEOUT_MS}ms`);
    console.log(`  Redis Host: ${REDIS_HOST}`);
    console.log(`  Redis Port: ${REDIS_PORT}`);
    console.log(`  Redis Database: ${REDIS_DATABASE}`);
    console.log(`  Redis Username: ${REDIS_USERNAME}`);
    console.log(`  Redis Password: ${REDIS_PASSWORD}`);
    console.log(`  Redis URL: ${REDIS_URL}`);
    console.log(`  Redis Cache Key Prefix: ${REDIS_CACHE_KEY_PREFIX}`);
}
