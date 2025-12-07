"""Configuration management for the AISHE application."""

import os
from typing import Optional


class Config:
    """Application configuration."""

    # Server configuration
    SERVER_HOST: str = os.getenv("AISHE_SERVER_HOST", "0.0.0.0")
    SERVER_PORT: int = int(os.getenv("AISHE_SERVER_PORT", "8000"))

    # API configuration
    API_URL: str = os.getenv(
        "AISHE_API_URL",
        f"http://localhost:{SERVER_PORT}"
    )

    # Ollama configuration
    OLLAMA_HOST: str = os.getenv("OLLAMA_HOST", "http://localhost:11434")
    OLLAMA_MODEL: str = os.getenv("AISHE_OLLAMA_MODEL", "llama3.2:3b")

    # RAG configuration
    MAX_CONTEXT_LENGTH: int = int(
        os.getenv("AISHE_MAX_CONTEXT_LENGTH", "4000"))
    MAX_SEARCH_RESULTS: int = int(os.getenv("AISHE_MAX_SEARCH_RESULTS", "3"))

    # Client configuration
    REQUEST_TIMEOUT: float = float(os.getenv("AISHE_REQUEST_TIMEOUT", "120.0"))

    @classmethod
    def get_server_url(cls) -> str:
        """Get the full server URL.

        Returns:
            Server URL in format http://host:port
        """
        return f"http://{cls.SERVER_HOST}:{cls.SERVER_PORT}"

    @classmethod
    def display_config(cls):
        """Display current configuration."""
        print("Configuration:")
        print(f"  Server Host: {cls.SERVER_HOST}")
        print(f"  Server Port: {cls.SERVER_PORT}")
        print(f"  API URL: {cls.API_URL}")
        print(f"  Ollama Host: {cls.OLLAMA_HOST}")
        print(f"  Ollama Model: {cls.OLLAMA_MODEL}")
        print(f"  Max Context Length: {cls.MAX_CONTEXT_LENGTH}")
        print(f"  Max Search Results: {cls.MAX_SEARCH_RESULTS}")
        print(f"  Request Timeout: {cls.REQUEST_TIMEOUT}s")


# Create a singleton instance
config = Config()
