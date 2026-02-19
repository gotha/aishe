/**
 * Simple command line application to ask questions to AISHE server.
 *
 * Usage: java Main <your question>
 * Example: java Main "What is the capital of France?"
 */
public class Main {
    public static void main(String[] args) {
        // Check if question was provided
        if (args.length == 0) {
            System.out.println("Usage: java Main <your question>");
            System.out.println("Example: java Main \"What is the capital of France?\"");
            System.exit(1);
        }

        // Get question from command line arguments
        String question = String.join(" ", args);

        System.out.println("Asking: " + question);
        System.out.println("Waiting for response...\n");

        throw new UnsupportedOperationException("not implemented ...");
    }
}
