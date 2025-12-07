#!/usr/bin/env node
// bin/aishe-cli.ts

/**
 * Command-line interface for RAG question answering.
 * 
 * Example usage:
 * 
 * ```ts
 * import { AIsheCLI } from "aishe-client";
 * 
 * const apiUrl = process.env.AISHE_API_URL;
 * const cli = new AIsheCLI(apiUrl);
 * await cli.run();
 * ```
 * 
 * NOTE: 100% ported by chatGPT from the Python version. NOT tested extensively or refactored.
 */

import readline from "node:readline/promises";
import { stdin as input, stdout as output } from "node:process";

import { RAGAPIClient } from "./client.js";
import { ServerNotReachableError, ServerError, APIClientError } from "./errors.js";
import type { AnswerResponse, HealthResponse } from "./models.js";

/**
 * Command-line interface for RAG question answering.
 */
export class AIsheCLI {
  private apiClient: RAGAPIClient;
  private rl = readline.createInterface({ input, output });

  /**
   * Initialize the CLI.
   *
   * @param apiUrl Optional API URL. If undefined, uses env/default from client.
   */
  constructor(apiUrl?: string) {
    this.apiClient = new RAGAPIClient(apiUrl);
  }

  /**
   * Print welcome banner.
   */
  private printBanner(): void {
    console.log("=".repeat(70));
    console.log("Wikipedia RAG Question Answering System");
    console.log("=".repeat(70));
    console.log("Ask questions and get answers based on Wikipedia articles.");
    console.log("Type 'quit' or 'exit' to stop.");
    console.log("=".repeat(70));
    console.log();
  }

  /**
   * Print RAG result in a formatted way.
   */
  private printResult(result: AnswerResponse): void {
    console.log("\n" + "─".repeat(70));
    console.log("ANSWER:");
    console.log("─".repeat(70));
    console.log(result.answer);

    if (result.sources && result.sources.length > 0) {
      console.log("\n" + "─".repeat(70));
      console.log("SOURCES:");
      console.log("─".repeat(70));
      for (const source of result.sources) {
        console.log(`[${source.number}] ${source.title}`);
        console.log(`    ${source.url}`);
      }
    }

    console.log("─".repeat(70));
  }

  /**
   * Run the interactive CLI.
   */
  async run(): Promise<void> {
    this.printBanner();

    // Handle CTRL+C at readline level
    this.rl.on("SIGINT", () => {
      console.log("\nReceived CTRL+C. Exiting...");
      this.rl.close();
    });

    // Check server health on startup
    try {
      console.log("Checking server connection...");
      const health: HealthResponse = await this.apiClient.checkHealth();
      if (health.status === "healthy" || health.status === "ok") {
        console.log("✓ Connected to server");
      } else {
        console.log(`⚠ Server status: ${health.status}`);
        if (health.message) {
          console.log(`  ${health.message}`);
        }
      }
      console.log();
    } catch (error) {
      if (error instanceof ServerNotReachableError) {
        console.error(`\n❌ Error: ${error.message}`);
        console.error("\nPlease start the server first:");
        console.error("  nix run .#server");
        console.error("\nOr set AISHE_API_URL to point to a running server.");
        process.exit(1);
      } else {
        console.error(`\n⚠ Warning: Could not check server health: ${error}`);
        console.error("Continuing anyway...\n");
      }
    }

    // Main loop
    while (true) {
      try {
        const question = (await this.rl.question("\nYour question: ")).trim();

        // Exit commands
        if (["quit", "exit", "q"].includes(question.toLowerCase())) {
          console.log("\nGoodbye!");
          break;
        }

        // Skip empty questions
        if (!question) {
          continue;
        }

        console.log("\nSearching Wikipedia and generating answer...");
        const result = await this.apiClient.askQuestion(question);
        this.printResult(result);
      } catch (error: unknown) {
        if (error instanceof ServerNotReachableError) {
          console.error(`\n❌ Server Error: ${error.message}`);
          console.error("\nThe server may have stopped. Please restart it:");
          console.error("  nix run .#server");
        } else if (error instanceof ServerError) {
          console.error(`\n❌ Server Error: ${error.message}`);
        } else if (error instanceof APIClientError) {
          console.error(`\n❌ Error: ${error.message}`);
        } else if (error instanceof Error && error.name === "AbortError") {
          console.error("\n❌ Request aborted.");
        } else {
          console.error(`\nUnexpected error: ${error}`);
          console.error(error);
        }
      }
    }

    this.rl.close();
  }
}
