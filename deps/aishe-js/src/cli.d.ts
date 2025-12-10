#!/usr/bin/env node
/**
 * Command-line interface for RAG question answering.
 */
export declare class AIsheCLI {
    private apiClient;
    private rl;
    /**
     * Initialize the CLI.
     *
     * @param apiUrl Optional API URL. If undefined, uses env/default from client.
     */
    constructor(apiUrl?: string);
    /**
     * Print welcome banner.
     */
    private printBanner;
    /**
     * Print RAG result in a formatted way.
     */
    private printResult;
    /**
     * Run the interactive CLI.
     */
    run(): Promise<void>;
}
//# sourceMappingURL=cli.d.ts.map