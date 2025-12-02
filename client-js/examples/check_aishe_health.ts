import { RAGAPIClient } from "../src/client.js";

console.log("EXAMPLE: Checking AIshe's health... 🤖");

const client = new RAGAPIClient();
const health = await client.checkHealth();
console.log("Health:", health);
