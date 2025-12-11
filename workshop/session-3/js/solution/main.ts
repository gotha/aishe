import readline from "node:readline/promises";
import { stdin as input, stdout as output } from "node:process";

import {
    AISHE_API_URL,
    REQUEST_TIMEOUT_MS,
    LANGCACHE_STRICT_SIMILARITY_THRESHOLD,
    LANGCACHE_CLOSE_SIMILARITY_THRESHOLD,
    LANGCACHE_LOOSE_SIMILARITY_THRESHOLD,
    type AnswerResponse,
    type HealthResponse,
} from "aishe-client";

import { AIsheHTTPClient } from "./client.js";

async function main(): Promise<void> {
    // TODO: create a new AIsheHTTPClient instance
    const client = await AIsheHTTPClient.create(
        AISHE_API_URL,
        REQUEST_TIMEOUT_MS,
        LANGCACHE_CLOSE_SIMILARITY_THRESHOLD,
    );

    console.log("Checking AIshe's health...");
    // TODO: check AIshe's health
    // Hint: you'll need to use the 'await' operator with async functions.
    const health: HealthResponse = await client.checkHealth();

    // TODO: print the health status
    // Hint: you need to print `status`, `ollama_accessible`, and `message` if it exists.
    const status = health.status;
    const ollamaAccessible = health.ollama_accessible;
    const message = health.message;

    console.log("======================================================================");
    console.log("AIshe server status:", status);
    console.log("Ollama accessible:", ollamaAccessible);
    if (message) {
        console.log("AIshe server message:", message);
    }
    console.log("======================================================================");
    console.log("\n");

    // Interactive question loop
    console.log("=== AIshe Question Answering (Session 3: LangCache Caching) ===");
    console.log("Type your question and press Enter to get an answer.");
    console.log("Enter CTRL+C, 'quit', 'exit', or 'q' to exit.");
    console.log();

    const rl = readline.createInterface({ input, output });
    // Handle CTRL+C at readline level
    rl.on("SIGINT", async () => {
        console.log("\nReceived CTRL+C. Exiting...");
        await client.close();
        rl.close();
    });

    while (true) {
        const question = (await rl.question("Your question: ")).trim();
        if (["quit", "exit", "q"].includes(question.toLowerCase())) {
            console.log("\nGoodbye!");
            break;
        }
        // Skip empty questions
        if (!question) {
            continue;
        }

        // TODO: ask AIshe a question, handle errors, measure execution time
        // Hint: use performance for measuring execution time
        // Hint: performance measures in milliseconds, so you need to convert it to seconds
        let answer: AnswerResponse;
        const startTime = performance.now();
        try {
            answer = await client.askQuestion(question);
        } catch (error) {
            console.error("Error:", error);
            continue;
        }
        const endTime = performance.now();

        // Asking: Does France have a capital?
        //
        // ======================================================================
        // ANSWER:
        // ======================================================================
        // Yes, France has two capitals. The capital of the country is Paris, which serves as the administrative center and seat of government.
        // However, there is another entity called "capital district" or "grand-duché" in French that holds special status, and it's the Grand Duchy of Luxembourg (although that's outside the scope)
        //
        // ======================================================================
        // SOURCES:
        // ======================================================================
        // [1] Capital punishment in France
        //    https://en.wikipedia.org/wiki/Capital_punishment_in_France
        // [2] Capital punishment by country
        //    https://en.wikipedia.org/wiki/Capital_punishment_by_country
        // [3] Capital districts and territories
        //    https://en.wikipedia.org/wiki/Capital_districts_and_territories

        // ======================================================================
        // Source: ASIHE API
        // Processing time: 2.345 seconds
        // Measured execution time: 2.531 seconds
        // ======================================================================

        const source = (await client.isCached(question)) ? "LangCache" : "AIshe API";
        const processingTime = answer.processing_time;
        const measuredTime = (endTime - startTime) / 1000;

        console.log(`Asking: ${question}`);
        console.log("\n");
        console.log("======================================================================");
        console.log("ANSWER:");
        console.log("======================================================================");
        console.log(answer.answer);
        console.log("\n");
        console.log("======================================================================");
        console.log("SOURCES:");
        console.log("======================================================================");
        for (const source of answer.sources) {
            console.log(`  [${source.number}] ${source.title}`);
            console.log(`      ${source.url}`);
        }
        console.log("\n");
        console.log("======================================================================");
        console.log("Source:", source);
        console.log("Processing time:", processingTime * 1000, "ms");
        console.log("Measured execution time:", measuredTime * 1000, "ms");
        console.log("======================================================================");
        console.log("\n");
    }

    await client.close();
    rl.close();
}

// Run only when executed directly
if (process.argv[1] && import.meta.url.endsWith(process.argv[1])) {
    main().catch((err) => {
        console.error("Fatal error:", err);
        process.exit(1);
    });
}
