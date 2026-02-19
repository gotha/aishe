package com.aishe.client;

import com.aishe.client.exceptions.APIClientError;
import com.aishe.client.exceptions.ServerError;
import com.aishe.client.exceptions.ServerNotReachableError;
import com.aishe.client.models.AnswerResponse;
import com.aishe.client.models.QuestionRequest;
import com.google.gson.Gson;
import com.google.gson.JsonSyntaxException;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.SocketTimeoutException;
import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.util.stream.Collectors;

/**
 * AISHE API client for interacting with the RAG server.
 */
public class RAGAPIClient {
    private final String baseUrl;
    private final int timeout;
    private final Gson gson;

    /**
     * Initialize AISHE API client.
     *
     * @param baseUrl Base URL of the AISHE API server (default: Config.AISHE_API_URL)
     * @param timeout Request timeout in milliseconds (default: Config.REQUEST_TIMEOUT_MS)
     */
    public RAGAPIClient(String baseUrl, int timeout) {
        this.baseUrl = baseUrl != null ? baseUrl : Config.AISHE_API_URL;
        this.timeout = timeout > 0 ? timeout : Config.REQUEST_TIMEOUT_MS;
        this.gson = new Gson();
    }

    /**
     * Initialize AISHE API client with default configuration.
     */
    public RAGAPIClient() {
        this(null, 0);
    }

    /**
     * Ask AISHE a question.
     *
     * Answer response is in the format:
     * {
     *   "answer": "The answer to the question",
     *   "sources": [
     *     {
     *       "number": 1,
     *       "title": "Source title",
     *       "url": "Source URL"
     *     }
     *   ],
     *   "processing_time": 0.123
     * }
     *
     * @param question Question to ask
     * @return Answer from AISHE server
     * @throws IllegalArgumentException If the question is empty
     * @throws ServerNotReachableError If the request timed out
     * @throws ServerError If the request failed
     * @throws APIClientError If the request or response is malformed
     */
    public AnswerResponse askQuestion(String question) throws APIClientError {
        if (question == null || question.trim().isEmpty()) {
            throw new IllegalArgumentException("Question cannot be empty");
        }

        String endpoint = askEndpoint();
        QuestionRequest request = new QuestionRequest(question.trim());
        String requestBody = gson.toJson(request);
        
        String responseBody = aisheRequest("POST", endpoint, requestBody);
        
        try {
            AnswerResponse response = gson.fromJson(responseBody, AnswerResponse.class);
            if (!isValidAnswerResponse(response)) {
                throw new APIClientError("Malformed answer response from server");
            }
            return response;
        } catch (JsonSyntaxException e) {
            throw new APIClientError("Failed to parse answer response", e);
        }
    }

    /**
     * Make a request to AISHE API server.
     *
     * @param method HTTP method to use (GET, POST, etc.)
     * @param endpoint Endpoint to send request to
     * @param body Body of the request (optional, null for GET requests)
     * @return Response body as string
     * @throws ServerNotReachableError If the request timed out
     * @throws ServerError If the request failed
     * @throws APIClientError If the request is malformed
     */
    private String aisheRequest(String method, String endpoint, String body) throws APIClientError {
        HttpURLConnection conn = null;
        try {
            URI uri = URI.create(endpoint);
            conn = (HttpURLConnection) uri.toURL().openConnection();
            conn.setRequestMethod(method);
            conn.setRequestProperty("Content-Type", "application/json");
            conn.setConnectTimeout(timeout);
            conn.setReadTimeout(timeout);

            // Send request body if provided
            if (body != null) {
                conn.setDoOutput(true);
                try (OutputStream os = conn.getOutputStream()) {
                    byte[] input = body.getBytes(StandardCharsets.UTF_8);
                    os.write(input, 0, input.length);
                }
            }

            // Check response code
            int responseCode = conn.getResponseCode();
            if (responseCode != 200) {
                throw new ServerError(method + " request failed! Status: " + responseCode);
            }

            // Read response
            try (BufferedReader br = new BufferedReader(
                    new InputStreamReader(conn.getInputStream(), StandardCharsets.UTF_8))) {
                return br.lines().collect(Collectors.joining("\n"));
            }

        } catch (SocketTimeoutException e) {
            throw new ServerNotReachableError(method + " request timed out after " + timeout + "ms", e);
        } catch (IOException e) {
            throw new APIClientError(method + " request failed! Unexpected error: " + e.getMessage(), e);
        } finally {
            if (conn != null) {
                conn.disconnect();
            }
        }
    }

    /**
     * Get the endpoint for asking a question to AISHE server.
     *
     * @return Endpoint for asking a question to AISHE server
     */
    private String askEndpoint() {
        return baseUrl + "/api/v1/ask";
    }

    /**
     * Validate answer response from AISHE server.
     *
     * @param response Answer response from AISHE server
     * @return True if the answer response is valid, false otherwise
     */
    private boolean isValidAnswerResponse(AnswerResponse response) {
        return response != null &&
               response.getAnswer() != null &&
               response.getSources() != null;
    }

    /**
     * Get the base URL of the AISHE API server.
     *
     * @return Base URL
     */
    public String getBaseUrl() {
        return baseUrl;
    }

    /**
     * Get the request timeout in milliseconds.
     *
     * @return Timeout in milliseconds
     */
    public int getTimeout() {
        return timeout;
    }
}

