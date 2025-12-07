import { AIsheCLI } from "./cli.js";

/**
 * Entry point for the CLI.
 */
async function main(): Promise<void> {
    const apiUrl = process.env.AISHE_API_URL;
    const cli = new AIsheCLI(apiUrl);
    await cli.run();
}

// Run only when executed directly
if (process.argv[1] && import.meta.url.endsWith(process.argv[1])) {
    main().catch((err) => {
        console.error("Fatal error:", err);
        process.exit(1);
    });
}
