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
 * Display the current configuration
 */
export function displayConfig(): void {
    console.log("Configuration:");
    console.log(`  Server Host: ${SERVER_HOST}`);
    console.log(`  Server Port: ${SERVER_PORT}`);
    console.log(`  API URL: ${AISHE_API_URL}`);
    console.log(`  Request Timeout: ${REQUEST_TIMEOUT_MS}ms`);
}
