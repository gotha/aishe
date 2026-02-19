import com.aishe.client.RAGAPIClient;
import com.aishe.client.models.AnswerResponse;
import com.aishe.client.models.Source;
import com.aishe.client.exceptions.APIClientError;
import io.github.cdimascio.dotenv.Dotenv;

import java.util.List;
import java.util.ArrayList;

/**
 * Command line application to ask questions to AISHE server with semantic caching using LangCache.
 *
 * Usage: java Main <your question>
 * Example: java Main "What is the capital of France?"
 */
public class Main {
    private static final double DEFAULT_SIMILARITY_THRESHOLD = 0.8;

    public static void main(String[] args) {
        long startTime = System.currentTimeMillis();

        // Check if question was provided
        if (args.length == 0) {
            System.out.println("Usage: java Main <your question>");
            System.out.println("Example: java Main \"What is the capital of France?\"");
            System.exit(1);
        }

        // Get question from command line arguments
        String question = String.join(" ", args);

        // Load environment variables from .env file
        Dotenv dotenv = Dotenv.configure()
                .directory(".")
                .ignoreIfMissing()
                .load();

        // Get credentials from environment variables
        String apiKey = getEnv(dotenv, "API_KEY");
        String cacheId = getEnv(dotenv, "CACHE_ID");
        String serverUrl = getEnv(dotenv, "SERVER_URL");

        // Get similarity threshold from environment variable (default: 0.8)
        double threshold = DEFAULT_SIMILARITY_THRESHOLD;
        String thresholdStr = getEnv(dotenv, "SIMILARITY_THRESHOLD");
        if (thresholdStr != null && !thresholdStr.isEmpty()) {
            try {
                threshold = Double.parseDouble(thresholdStr);
            } catch (NumberFormatException ignored) {}
        }

        // Validate required credentials
        List<String> missingFields = new ArrayList<>();
        if (apiKey == null || apiKey.isEmpty()) {
            missingFields.add("API_KEY");
        }
        if (cacheId == null || cacheId.isEmpty()) {
            missingFields.add("CACHE_ID");
        }
        if (serverUrl == null || serverUrl.isEmpty() || serverUrl.equals("YOUR_REDIS_CLOUD_LANGCACHE_HOST_HERE")) {
            missingFields.add("SERVER_URL");
        }

        if (!missingFields.isEmpty()) {
            System.out.println("Error: Missing or invalid credentials in .env file");
            System.out.println("Missing fields: " + String.join(", ", missingFields));
            System.out.println("\nPlease update .env file with your Redis Cloud LangCache credentials:");
            System.out.println("- SERVER_URL: Your Redis Cloud LangCache host");
            System.out.println("- CACHE_ID: Your cache ID");
            System.out.println("- API_KEY: Your LangCache API key");
            System.exit(1);
        }

        // Ensure server_url has https:// prefix
        if (!serverUrl.startsWith("http")) {
            serverUrl = "https://" + serverUrl;
        }

        // Initialize LangCache client
        LangCacheClient langCache = new LangCacheClient(serverUrl, cacheId, apiKey);

        System.out.println("Asking: " + question);

        // Check cache first using semantic search
        AnswerResponse response;
        boolean fromCache;
        Double similarity = null;

        try {
            CachedResponse cachedResponse = getFromCache(langCache, question, threshold);

            if (cachedResponse != null) {
                System.out.println("✓ Found in semantic cache! (no API call needed)");
                if (cachedResponse.similarity != null) {
                    System.out.printf("  Similarity score: %.4f%n", cachedResponse.similarity);
                }
                System.out.println();
                response = cachedResponse.response;
                similarity = cachedResponse.similarity;
                fromCache = true;
            } else {
                System.out.println("✗ Not in cache, calling AISHE API...");
                System.out.println("Waiting for response...\n");

                // Create AISHE client
                RAGAPIClient client = new RAGAPIClient();
                response = client.askQuestion(question);

                // Save to semantic cache for future use
                try {
                    saveToCache(langCache, question, response);
                    System.out.println("✓ Response saved to semantic cache\n");
                } catch (Exception e) {
                    System.out.println("Warning: Error saving to cache: " + e.getMessage());
                }
                fromCache = false;
            }
        } catch (APIClientError e) {
            System.err.println("Error: " + e.getMessage());
            System.exit(1);
            return;
        } catch (Exception e) {
            System.err.println("Error: " + e.getMessage());
            System.exit(1);
            return;
        }

        // Display results
        displayResponse(response, fromCache, similarity, startTime);
    }

    private static String getEnv(Dotenv dotenv, String key) {
        String value = dotenv.get(key);
        if (value == null) {
            value = System.getenv(key);
        }
        return value;
    }

    private static CachedResponse getFromCache(LangCacheClient client, String question, double threshold) {
        throw new UnsupportedOperationException();
    }

    private static void saveToCache(LangCacheClient client, String question, AnswerResponse response) throws Exception {
        throw new UnsupportedOperationException();
    }

    private static void displayResponse(AnswerResponse response, boolean fromCache, Double similarity, long startTime) {
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

        // Print processing time or cache info
        System.out.println();
        System.out.println("======================================================================");
        if (fromCache) {
            System.out.println("Source: Semantic Cache (LangCache)");
            if (similarity != null) {
                System.out.printf("Similarity score: %.4f%n", similarity);
            }
            System.out.printf("Original processing time: %.2f seconds%n", response.getProcessingTime());
        } else {
            System.out.printf("Processing time: %.2f seconds%n", response.getProcessingTime());
        }
        System.out.println("======================================================================");

        // Print total execution time
        double executionTime = (System.currentTimeMillis() - startTime) / 1000.0;
        System.out.println();
        System.out.println("======================================================================");
        System.out.printf("Execution time: %.2f seconds%n", executionTime);
        System.out.println("======================================================================");
    }
}
