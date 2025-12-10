/**
 * Base exception for API client errors
 */
export declare class APIClientError extends Error {
    /**
     * Constructor for APIClientError
     * @param message - The error message
     */
    constructor(message: string);
}
/**
 * Exception raised when the server is not reachable
 */
export declare class ServerNotReachableError extends APIClientError {
    /**
     * Constructor for ServerNotReachableError
     * @param message - The error message
     */
    constructor(message: string);
}
/**
 * Exception raised when the server returns an error
 */
export declare class ServerError extends APIClientError {
    /**
     * Constructor for ServerError
     * @param message - The error message
     */
    constructor(message: string);
}
//# sourceMappingURL=errors.d.ts.map