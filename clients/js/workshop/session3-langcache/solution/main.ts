import readline from "node:readline/promises";
import { stdin as input, stdout as output } from "node:process";

import {
    AISHE_API_URL,
    REQUEST_TIMEOUT_MS,
    LANGCACHE_LOOSE_SIMILARITY_THRESHOLD,
    LANGCACHE_CLOSE_SIMILARITY_THRESHOLD,
    LANGCACHE_STRICT_SIMILARITY_THRESHOLD,
    type AnswerResponse,
    type HealthResponse,
} from "aishe-client";

import { AIsheHTTPClient } from "./client.js";

async function main(): Promise<void> {
    // TODO: create a new AIsheHTTPClient instance
    const client = await AIsheHTTPClient.create(
        AISHE_API_URL,
        REQUEST_TIMEOUT_MS,
        LANGCACHE_STRICT_SIMILARITY_THRESHOLD,
    );

    // TODO: check AIshe's health
    // Hint: you'll need to use the 'await' operator with async functions.
    console.log("Checking AIshe's health...");
    const health: HealthResponse = await client.checkHealth();

    // TODO: print the health status
    // Hint: you need to print `status`, `ollama_accessible`, and `message` if it exists.
    console.log("AIshe server status:", health.status);
    console.log("Ollama accessible:", health.ollama_accessible);
    if (health.message) {
        console.log("AIshe server message:", health.message);
    }

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

        // TODO: check if the question is cached in LangCache
        //       - if found (cache HIT), mark source as LangCache
        //       - otherwise, mark source as AIshe API
        const retrievalSource: string = (await client.isCached(question)) ? "LangCache" : "AIshe API";

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
        const measuredTime = endTime - startTime;

        // TODO: dispaly results
        // Expected output format:
        //
        // Answer: <answer>
        // Source: AIshe API
        //
        // Processing time: <processing_time>
        // Measured execution time: <measured_time>
        //
        // Wikipedia sources:
        //   [1] <title>
        //       <url>
        //   [2] <title>
        //       <url>
        //   [3] <title>
        //       <url>

        console.log("Answer:", answer.answer);
        console.log("Source:", retrievalSource);
        console.log("Processing time:", answer.processing_time * 1000, "ms");
        console.log("Measured execution time:", measuredTime, "ms");
        for (const source of answer.sources) {
            console.log(`  [${source.number}] ${source.title}`);
            console.log(`      ${source.url}`);
        }
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
