/**
 * AIshe API server host
 */
export const SERVER_HOST = process.env.SERVER_HOST || "localhost";
/**
 * AIshe API server port
 */
export const SERVER_PORT = parseInt(process.env.SERVER_PORT || "8000");
/**
 * AIshe API URL
 */
export const AISHE_API_URL = process.env.AISHE_API_URL || `http://${SERVER_HOST}:${SERVER_PORT}`;
/**
 * Request timeout in seconds
 */
export const REQUEST_TIMEOUT_MS = parseInt(process.env.REQUEST_TIMEOUT || "120") * 1000;
/**
 * Redis host
 */
export const REDIS_HOST = process.env.REDIS_HOST || "localhost";
/**
 * Redis port
 */
export const REDIS_PORT = parseInt(process.env.REDIS_PORT || "6379");
/**
 * Redis database
 */
export const REDIS_DATABASE = parseInt(process.env.REDIS_DATABASE || "0");
/**
 * Redis username
 */
export const REDIS_USERNAME = process.env.REDIS_USERNAME || "default";
/**
 * Redis password
 */
export const REDIS_PASSWORD = process.env.REDIS_PASSWORD || "";
/**
 * Redis URL
 */
export const REDIS_URL = process.env.REDIS_URL ||
    `redis://${REDIS_USERNAME}:${REDIS_PASSWORD}@${REDIS_HOST}:${REDIS_PORT}/${REDIS_DATABASE}`;
/**
 * Redis cache key prefix
 */
export const REDIS_CACHE_KEY_PREFIX = process.env.REDIS_CACHE_KEY_PREFIX || "aishe:question";
/**
 * LangCache strict similarity threshold
 */
export const LANGCACHE_STRICT_SIMILARITY_THRESHOLD = parseFloat(process.env.LANGCACHE_STRICT_SIMILARITY_THRESHOLD || "0.95");
/**
 * LangCache close similarity threshold
 */
export const LANGCACHE_CLOSE_SIMILARITY_THRESHOLD = parseFloat(process.env.LANGCACHE_CLOSE_SIMILARITY_THRESHOLD || "0.9");
/**
 * LangCache loose similarity threshold
 */
export const LANGCACHE_LOOSE_SIMILARITY_THRESHOLD = parseFloat(process.env.LANGCACHE_LOOSE_SIMILARITY_THRESHOLD || "0.8");
/**
 * LangCache API key
 */
export const LANGCACHE_API_KEY = process.env.LANGCACHE_API_KEY;
/**
 * LangCache cache ID
 */
export const LANGCACHE_CACHE_ID = process.env.LANGCACHE_CACHE_ID;
/**
 * LangCache server URL
 */
export const LANGCACHE_SERVER_URL = process.env.LANGCACHE_SERVER_URL;
/**
 * Display the current configuration
 */
export function displayConfig() {
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
    console.log(`  LangCache Strict Similarity Threshold: ${LANGCACHE_STRICT_SIMILARITY_THRESHOLD}`);
    console.log(`  LangCache Close Similarity Threshold: ${LANGCACHE_CLOSE_SIMILARITY_THRESHOLD}`);
    console.log(`  LangCache Loose Similarity Threshold: ${LANGCACHE_LOOSE_SIMILARITY_THRESHOLD}`);
}
//# sourceMappingURL=config.js.map