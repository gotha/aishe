async function main(): Promise<void> {
    console.log("Session 1 solution");
}

// Run only when executed directly
if (process.argv[1] && import.meta.url.endsWith(process.argv[1])) {
    main().catch((err) => {
        console.error("Fatal error:", err);
        process.exit(1);
    });
}
