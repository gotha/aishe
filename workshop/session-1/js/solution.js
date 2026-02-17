import { createInterface } from "node:readline/promises";
const { stdin: input, stdout: output } = process;

const AISHE_API_URL = process.env["AISHE_API_URL"];

const rl = createInterface({ input, output });

let startTime = 0;

async function main() {
  const answer = await rl.question("Please enter your question: ");

  startTime = performance.now();

  const question = answer.trim();
  if (!question) {
    console.log("Error: You must provide a question.");
    process.exit(1);
  }

  console.log(`Asking: ${question}`);
  console.log("Waiting for response...\n");

  const response = await fetch(`${AISHE_API_URL}/api/v1/ask`, {
    body: JSON.stringify({ question }),
    headers: {
      "Content-Type": "application/json",
    },
    method: "POST",
  });

  const data = await response.json();

  // Print answer
  console.log("=".repeat(70));
  console.log("ANSWER:");
  console.log("=".repeat(70));
  console.log(data.answer);

  // Print sources if available
  if (data.sources?.length) {
    console.log("\n" + "=".repeat(70));
    console.log("SOURCES:");
    console.log("=".repeat(70));
    for (const source of data.sources) {
      console.log(`[${source.number}] ${source.title}`);
      console.log(`    ${source.url}`);
    }
  }
  // Print processing time
  console.log("\n" + "=".repeat(70));
  console.log(`Processing time: ${data.processing_time.toFixed(2)} seconds`);
  console.log("=".repeat(70));
}

try {
  await main();
} catch (error) {
  console.error(error);
  process.exit(1);
} finally {
  rl.close();
  console.log(
    `Total time:  ${((performance.now() - startTime) / 1000).toFixed(2)} seconds`,
  );
  console.log("=".repeat(70));
}
