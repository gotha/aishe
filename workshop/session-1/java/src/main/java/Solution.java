import com.aishe.client.RAGAPIClient;
import com.aishe.client.models.AnswerResponse;
import com.aishe.client.models.Source;
import com.aishe.client.exceptions.APIClientError;

/**
 * Simple command line application to ask questions to AISHE server.
 *
 * Usage: java Main <your question>
 * Example: java Main "What is the capital of France?"
 */
public class Solution {
    public static void main(String[] args) {
        // Check if question was provided
        if (args.length == 0) {
            System.out.println("Usage: java Main <your question>");
            System.out.println("Example: java Main \"What is the capital of France?\"");
            System.exit(1);
        }

        // Get question from command line arguments
        String question = String.join(" ", args);

        // Start timing
        long startTime = System.currentTimeMillis();

        System.out.println("Asking: " + question);
        System.out.println("Waiting for response...\n");

        // Create AISHE client
        // Uses AISHE_API_URL environment variable or defaults to http://localhost:8000
        RAGAPIClient client = new RAGAPIClient();

        try {
            // Ask question
            AnswerResponse response = client.askQuestion(question);

            // Calculate total execution time
            double executionTime = (System.currentTimeMillis() - startTime) / 1000.0;

            // Display results
            System.out.println("======================================================================");
            System.out.println("ANSWER:");
            System.out.println("======================================================================");
            System.out.println(response.getAnswer());
            System.out.println();

            System.out.println("======================================================================");
            System.out.println("SOURCES:");
            System.out.println("======================================================================");
            for (Source source : response.getSources()) {
                System.out.printf("[%d] %s%n", source.getNumber(), source.getTitle());
                System.out.printf("    %s%n", source.getUrl());
            }
            System.out.println();

            System.out.println("======================================================================");
            System.out.printf("Processing time: %.2f seconds%n", response.getProcessingTime());
            System.out.printf("Total execution time: %.2f seconds%n", executionTime);
            System.out.println("======================================================================");

        } catch (APIClientError e) {
            System.err.println("Error: " + e.getMessage());
            e.printStackTrace();
            System.exit(1);
        }
    }
}

