import { createInterface } from "node:readline/promises";
const { stdin: input, stdout: output } = process;

const AISHE_API_URL = process.env["AISHE_API_URL"];

const rl = createInterface({ input, output });

let startTime = 0;

try {
  const answer = await rl.question("Please enter your question: ");

  startTime = performance.now();

  const question = answer.trim();

  if (!question) {
    console.log("Error: You must provide a question.");
    process.exit(1);
  }

  console.log(`Asking: ${question}`);
  console.log("Waiting for response...\n");
  
  // *******************************
  // * WARNING: Solution missing
  // * TODO: Fetch and parse Aishe response via HTTP GET request
  // *******************************
  throw new Error("Not Implemented!");
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
