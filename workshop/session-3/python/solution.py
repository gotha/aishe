#!/usr/bin/env python3
"""Command line script to ask questions to AISHE server with semantic caching using langcache."""

import sys
import json
import os
import requests
from langcache import LangCache
from dotenv import load_dotenv


def get_from_cache(lang_cache: LangCache, question: str) -> dict | None:
    """Search for a cached response using semantic search.

    Args:
        lang_cache: LangCache client instance
        question: The question to search for

    Returns:
        Cached response data or None if not found
    """
    try:
        # Search for semantically similar questions in the cache
        # Use a lower similarity threshold (0.8) to allow for more semantic matches
        result = lang_cache.search(prompt=question, similarity_threshold=0.8)

        # Check if we got any results
        if result and hasattr(result, 'data') and result.data:
            # Get the first (most similar) entry
            entry = result.data[0]

            # Parse the cached response from the entry's response field
            if hasattr(entry, 'response') and entry.response:
                # The response is stored as a JSON string
                cached_data = json.loads(entry.response)
                return cached_data

    except Exception as e:
        print(f"Warning: Error reading from cache: {e}")

    return None


def save_to_cache(lang_cache: LangCache, question: str, response_data: dict) -> None:
    """Save response to semantic cache.

    Args:
        lang_cache: LangCache client instance
        question: The question
        response_data: The response data to cache
    """
    try:
        # Convert response data to JSON string for storage
        response_json = json.dumps(response_data)

        # Save to langcache with the question as prompt and response as JSON
        lang_cache.set(prompt=question, response=response_json)

    except Exception as e:
        print(f"Warning: Error saving to cache: {e}")


def main():
    """Ask a question to AISHE server with semantic caching using langcache."""
    # Load environment variables from .env file
    env_path = os.path.join(os.path.dirname(__file__), '.env')
    load_dotenv(env_path)

    # Check if question was provided
    if len(sys.argv) < 2:
        print("Usage: python main.py <your question>")
        print("Example: python main.py 'What is the capital of France?'")
        sys.exit(1)

    # Get question from command line arguments
    question = " ".join(sys.argv[1:])

    # Get credentials from environment variables
    api_key = os.getenv('API_KEY')
    cache_id = os.getenv('CACHE_ID')
    server_url = os.getenv('SERVER_URL')

    # Validate required credentials
    missing_fields = []
    if not api_key:
        missing_fields.append('API_KEY')
    if not cache_id:
        missing_fields.append('CACHE_ID')
    if not server_url or server_url == 'YOUR_REDIS_CLOUD_LANGCACHE_HOST_HERE':
        missing_fields.append('SERVER_URL')

    if missing_fields:
        print(f"Error: Missing or invalid credentials in .env file")
        print(f"Missing fields: {', '.join(missing_fields)}")
        print("\nPlease update .env file with your Redis Cloud LangCache credentials:")
        print("- SERVER_URL: Your Redis Cloud LangCache host (e.g., 'your-instance.redis.cloud')")
        print("- CACHE_ID: Your cache ID")
        print("- API_KEY: Your LangCache API key")
        sys.exit(1)

    # Initialize LangCache client
    # Ensure server_url has https:// prefix
    if not server_url.startswith('http'):
        server_url = f'https://{server_url}'

    try:
        with LangCache(
            server_url=server_url,
            cache_id=cache_id,
            api_key=api_key,
        ) as lang_cache:

            print(f"Asking: {question}")

            # Check cache first using semantic search
            cached_response = get_from_cache(lang_cache, question)

            if cached_response:
                print("✓ Found in semantic cache! (no API call needed)\n")
                data = cached_response
                from_cache = True
            else:
                print("✗ Not in cache, calling AISHE API...")
                print("Waiting for response...\n")

                # AISHE server URL (running in Docker on port 8000)
                url = "http://localhost:8000/api/v1/ask"
                payload = {"question": question}

                try:
                    # Send POST request to AISHE server
                    response = requests.post(url, json=payload, timeout=120)
                    response.raise_for_status()

                    # Parse response
                    data = response.json()

                    # Save to semantic cache for future use
                    save_to_cache(lang_cache, question, data)
                    print("✓ Response saved to semantic cache\n")
                    from_cache = False

                except requests.exceptions.ConnectionError:
                    print(
                        "Error: Could not connect to AISHE server at http://localhost:8000")
                    print("Make sure the server is running in Docker.")
                    sys.exit(1)
                except requests.exceptions.Timeout:
                    print("Error: Request timed out. The question may be too complex.")
                    sys.exit(1)
                except requests.exceptions.HTTPError as e:
                    print(f"Error: Server returned an error: {e}")
                    if 'response' in locals():
                        print(f"Details: {response.text}")
                    sys.exit(1)
                except Exception as e:
                    print(f"Unexpected error: {e}")
                    sys.exit(1)

            # Print answer
            print("=" * 70)
            print("ANSWER:")
            print("=" * 70)
            print(data["answer"])

            # Print sources if available
            if data.get("sources"):
                print("\n" + "=" * 70)
                print("SOURCES:")
                print("=" * 70)
                for source in data["sources"]:
                    print(f"[{source['number']}] {source['title']}")
                    print(f"    {source['url']}")

            # Print processing time
            print("\n" + "=" * 70)
            if from_cache:
                print("Source: Semantic Cache (LangCache)")
            else:
                print(
                    f"Processing time: {data['processing_time']:.2f} seconds")
            print("=" * 70)

    except Exception as e:
        print(f"Error initializing LangCache: {e}")
        print("Please check your credentials in .env file")
        sys.exit(1)


if __name__ == "__main__":
    main()
