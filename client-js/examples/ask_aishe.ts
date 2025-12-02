import { RAGAPIClient } from "../src/client.js";

console.log("EXAMPLE: Asking AIshe a question... 🤖");

const client = new RAGAPIClient();
const question = "What is the Python programming language?";
const answer = await client.askQuestion(question);
console.log("Answer:", answer.answer);
console.log("Sources:", answer.sources);
console.log("Processing time:", answer.processing_time);
