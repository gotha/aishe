package com.aishe.client.exceptions;

/**
 * Exception raised when the server returns an error.
 */
public class ServerError extends APIClientError {
    /**
     * Constructor for ServerError.
     * 
     * @param message The error message
     */
    public ServerError(String message) {
        super(message);
    }

    /**
     * Constructor for ServerError with cause.
     * 
     * @param message The error message
     * @param cause The cause of the error
     */
    public ServerError(String message, Throwable cause) {
        super(message, cause);
    }
}

