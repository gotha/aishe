import com.aishe.client.RAGAPIClient;
import com.aishe.client.Config;
import com.aishe.client.Utils;
import com.aishe.client.models.AnswerResponse;
import com.aishe.client.models.Source;
import com.aishe.client.exceptions.APIClientError;
import com.google.gson.Gson;
import redis.clients.jedis.Jedis;
import redis.clients.jedis.exceptions.JedisConnectionException;

/**
 * Command line application to ask questions to AISHE server with Redis caching.
 *
 * Usage: java Main <your question>
 * Example: java Main "What is the capital of France?"
 */
public class Solution {
    private static final int CACHE_TTL_SECONDS = 86400; // 24 hours
    private static final Gson gson = new Gson();

    public static void main(String[] args) {
        // Check if question was provided
        if (args.length == 0) {
            System.out.println("Usage: java Main <your question>");
            System.out.println("Example: java Main \"What is the capital of France?\"");
            System.exit(1);
        }

        // Get question from command line arguments
        String question = String.join(" ", args);

        // Connect to Redis
        Jedis redis = null;
        try {
            redis = new Jedis(Config.REDIS_HOST, Config.REDIS_PORT);
            redis.ping(); // Test connection
        } catch (JedisConnectionException e) {
            System.err.println("Error: Could not connect to Redis at " + Config.REDIS_HOST + ":" + Config.REDIS_PORT);
            System.err.println("Make sure Redis is running in Docker.");
            System.exit(1);
        } catch (Exception e) {
            System.err.println("Error connecting to Redis: " + e.getMessage());
            System.exit(1);
        }

        System.out.println("Asking: " + question);

        // Create AISHE client
        RAGAPIClient client = new RAGAPIClient();

        // Check cache first
        String cachedResponse = getFromCache(redis, question);
        AnswerResponse response;
        boolean fromCache;

        if (cachedResponse != null) {
            System.out.println("✓ Found in cache! (no API call needed)\n");
            response = gson.fromJson(cachedResponse, AnswerResponse.class);
            fromCache = true;
        } else {
            System.out.println("✗ Not in cache, calling AISHE API...");
            System.out.println("Waiting for response...\n");

            try {
                // Make API call using client
                response = client.askQuestion(question);

                // Save to cache for future use
                saveToCache(redis, question, gson.toJson(response));
                System.out.println("✓ Response saved to cache\n");
                fromCache = false;

            } catch (APIClientError e) {
                System.err.println("Error: " + e.getMessage());
                redis.close();
                System.exit(1);
                return;
            }
        }

        // Display results
        displayResponse(response, fromCache);

        // Close Redis connection
        redis.close();
    }

    /**
     * Generate a cache key from the question using the client library utility.
     */
    private static String getCacheKey(String question) {
        return Utils.generateCacheKey(question);
    }

    /**
     * Get cached response for a question.
     */
    private static String getFromCache(Jedis redis, String question) {
        String cacheKey = getCacheKey(question);
        try {
            return redis.get(cacheKey);
        } catch (Exception e) {
            System.err.println("Warning: Error reading from cache: " + e.getMessage());
            return null;
        }
    }

    /**
     * Save response to cache with 24-hour expiration.
     */
    private static void saveToCache(Jedis redis, String question, String responseData) {
        String cacheKey = getCacheKey(question);
        try {
            redis.setex(cacheKey, CACHE_TTL_SECONDS, responseData);
        } catch (Exception e) {
            System.err.println("Warning: Error saving to cache: " + e.getMessage());
        }
    }

    /**
     * Display the response to the user.
     */
    private static void displayResponse(AnswerResponse response, boolean fromCache) {
        // Print answer
        System.out.println("======================================================================");
        System.out.println("ANSWER:");
        System.out.println("======================================================================");
        System.out.println(response.getAnswer());

        // Print sources if available
        if (response.getSources() != null && !response.getSources().isEmpty()) {
            System.out.println();
            System.out.println("======================================================================");
            System.out.println("SOURCES:");
            System.out.println("======================================================================");
            for (Source source : response.getSources()) {
                System.out.printf("[%d] %s%n", source.getNumber(), source.getTitle());
                System.out.printf("    %s%n", source.getUrl());
            }
        }

        // Print processing time or cache indicator
        System.out.println();
        System.out.println("======================================================================");
        if (fromCache) {
            System.out.println("Source: Redis Cache");
        } else {
            System.out.printf("Processing time: %.2f seconds%n", response.getProcessingTime());
        }
        System.out.println("======================================================================");
    }
}

