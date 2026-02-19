import { createInterface } from "node:readline/promises";
import { LangCache } from "@redis-ai/langcache";

const { stdin: input, stdout: output } = process;

/**
 * Searches the language cache using the given question.
 *
 * @param {LangCache} langCache - Cache instance used for semantic search.
 * @param {string} question - Query prompt.
 * @returns {Promise<any>} Search result from the cache.
 */
const getFromCache = async (langCache: LangCache, question: string) => {
  const { data } = await langCache.search({
    prompt: question,
    similarityThreshold: 0.8,
  });

  return data;
};

/**
 * Saves a prompt–response pair to the language cache.
 *
 * @param {LangCache} langCache - Cache instance used for storage.
 * @param {string} question - Prompt to cache.
 * @param {any} response - Response associated with the prompt.
 * @returns {Promise<any>} Result of the cache set operation.
 */
const saveToCache = async (
  langCache: LangCache,
  question: string,
  response: any,
) => {
  return langCache.set({
    prompt: question,
    response: JSON.stringify(response),
  });
};

const rl = createInterface({ input, output });
let startTime = 0;

const langCache = new LangCache({
  serverURL: process.env["LANGCACHE_SERVER_URL"]!,
  cacheId: process.env["LANGCACHE_CACHE_ID"],
  apiKey: process.env["LANGCACHE_API_KEY"],
});

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

  let data;
  const cached = await getFromCache(langCache, question);

  if (cached.length) {
    data = JSON.parse(cached[0].response);
  } else {
    const response = await fetch(
      `${process.env["AISHE_API_URL"]}/api/v1/ask`,
      {
        body: JSON.stringify({ question }),
        headers: {
          "Content-Type": "application/json",
        },
        method: "POST",
      },
    );

    data = await response.json();
    await saveToCache(langCache, question, data);
  }

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
    `Total time:  ${(performance.now() - startTime).toFixed(2)} ms`,
  );
  console.log("=".repeat(70));
}
