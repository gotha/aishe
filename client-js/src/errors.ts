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
