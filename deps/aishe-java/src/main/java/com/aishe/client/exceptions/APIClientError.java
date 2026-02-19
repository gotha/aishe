package com.aishe.client.exceptions;

/**
 * Base exception for API client errors.
 */
public class APIClientError extends Exception {
    /**
     * Constructor for APIClientError.
     * 
     * @param message The error message
     */
    public APIClientError(String message) {
        super(message);
    }

    /**
     * Constructor for APIClientError with cause.
     * 
     * @param message The error message
     * @param cause The cause of the error
     */
    public APIClientError(String message, Throwable cause) {
        super(message, cause);
    }
}

