package com.aishe.client.exceptions;

/**
 * Exception raised when the server is not reachable.
 */
public class ServerNotReachableError extends APIClientError {
    /**
     * Constructor for ServerNotReachableError.
     * 
     * @param message The error message
     */
    public ServerNotReachableError(String message) {
        super(message);
    }

    /**
     * Constructor for ServerNotReachableError with cause.
     * 
     * @param message The error message
     * @param cause The cause of the error
     */
    public ServerNotReachableError(String message, Throwable cause) {
        super(message, cause);
    }
}

