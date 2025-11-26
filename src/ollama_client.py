"""Ollama client for LLM interactions."""

from typing import AsyncIterator, Dict, List, Optional
import ollama


class OllamaClient:
    """Client for interacting with Ollama LLM."""

    def __init__(self, host: str = "http://localhost:11434", model: str = "llama3.2:3b"):
        """Initialize the Ollama client.

        Args:
            host: Ollama server host URL
            model: Default model to use for generation
        """
        self.host = host
        self.model = model
        self.client = ollama.Client(host=host)

    def list_models(self) -> List[Dict]:
        """List available models.

        Returns:
            List of available models with their details
        """
        response = self.client.list()
        return response.get('models', [])

    def generate(self, prompt: str, model: Optional[str] = None, **kwargs) -> str:
        """Generate a response from the LLM.

        Args:
            prompt: The prompt to send to the LLM
            model: Model to use (defaults to self.model)
            **kwargs: Additional parameters for generation

        Returns:
            Generated response text
        """
        model = model or self.model
        response = self.client.generate(
            model=model,
            prompt=prompt,
            **kwargs
        )
        return response.get('response', '')

    def chat(self, messages: List[Dict[str, str]], model: Optional[str] = None, **kwargs) -> str:
        """Chat with the LLM using a conversation format.

        Args:
            messages: List of message dictionaries with 'role' and 'content'
            model: Model to use (defaults to self.model)
            **kwargs: Additional parameters for chat

        Returns:
            Generated response text
        """
        model = model or self.model
        response = self.client.chat(
            model=model,
            messages=messages,
            **kwargs
        )
        return response.get('message', {}).get('content', '')

    def stream_generate(self, prompt: str, model: Optional[str] = None, **kwargs) -> AsyncIterator[str]:
        """Stream a response from the LLM.

        Args:
            prompt: The prompt to send to the LLM
            model: Model to use (defaults to self.model)
            **kwargs: Additional parameters for generation

        Yields:
            Chunks of generated text
        """
        model = model or self.model
        stream = self.client.generate(
            model=model,
            prompt=prompt,
            stream=True,
            **kwargs
        )

        for chunk in stream:
            if 'response' in chunk:
                yield chunk['response']

    def stream_chat(self, messages: List[Dict[str, str]], model: Optional[str] = None, **kwargs) -> AsyncIterator[str]:
        """Stream a chat response from the LLM.

        Args:
            messages: List of message dictionaries with 'role' and 'content'
            model: Model to use (defaults to self.model)
            **kwargs: Additional parameters for chat

        Yields:
            Chunks of generated text
        """
        model = model or self.model
        stream = self.client.chat(
            model=model,
            messages=messages,
            stream=True,
            **kwargs
        )

        for chunk in stream:
            if 'message' in chunk and 'content' in chunk['message']:
                yield chunk['message']['content']

    def generate_with_context(self, prompt: str, context: str, model: Optional[str] = None, **kwargs) -> str:
        """Generate a response with additional context.

        Args:
            prompt: The user's question or prompt
            context: Additional context to provide to the LLM
            model: Model to use (defaults to self.model)
            **kwargs: Additional parameters for generation

        Returns:
            Generated response text
        """
        full_prompt = f"""Context:
{context}

Question: {prompt}

Answer based on the context provided above:"""

        return self.generate(full_prompt, model=model, **kwargs)

    def chat_with_context(self, question: str, context: str, model: Optional[str] = None, **kwargs) -> str:
        """Chat with the LLM using context (RAG pattern).

        Args:
            question: The user's question
            context: Retrieved context to answer the question
            model: Model to use (defaults to self.model)
            **kwargs: Additional parameters for chat

        Returns:
            Generated response text
        """
        messages = [
            {
                "role": "system",
                "content": "You are a helpful assistant. Answer questions based on the provided context. If the context doesn't contain enough information, say so."
            },
            {
                "role": "user",
                "content": f"Context:\n{context}\n\nQuestion: {question}"
            }
        ]

        return self.chat(messages, model=model, **kwargs)
