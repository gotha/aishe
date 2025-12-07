import readline from "node:readline/promises";
import { stdin as input, stdout as output } from "node:process";

import type { AnswerResponse, HealthResponse } from "aishe-client";

import { AIsheHTTPClient } from "./client.js";

async function main(): Promise<void> {
    // TODO: create a new AIsheHTTPClient instance
    // ..........................

    // TODO: check AIshe's health
    // Hint: you'll need to use the 'await' operator with async functions.
    console.log("Checking AIshe's health...");
    // ..........................

    // TODO: print the health status
    // Hint: you need to print `status`, `ollama_accessible`, and `message` if it exists.
    // ..........................

    // Interactive question loop
    console.log("=== AIshe Question Answering (Session 1: Basic Client) ===");
    console.log("Type your question and press Enter to get an answer.");
    console.log("Enter CTRL+C, 'quit', 'exit', or 'q' to exit.");
    console.log();

    const rl = readline.createInterface({ input, output });
    // Handle CTRL+C at readline level
    rl.on("SIGINT", () => {
        console.log("\nReceived CTRL+C. Exiting...");
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
        // ..........................

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

        // ..........................
    }
    rl.close();
}

// Run only when executed directly
if (process.argv[1] && import.meta.url.endsWith(process.argv[1])) {
    main().catch((err) => {
        console.error("Fatal error:", err);
        process.exit(1);
    });
}
