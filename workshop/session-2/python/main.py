#!/usr/bin/env python3
"""Command line script to ask questions to AISHE server with Redis caching."""

import sys
import json
import hashlib
import requests
import redis


def get_cache_key(question: str) -> str:
    raise NotImplementedError


def get_from_cache(redis_client: redis.Redis, question: str) -> dict | None:
    raise NotImplementedError


def save_to_cache(redis_client: redis.Redis, question: str, response_data: dict) -> None:
    raise NotImplementedError


def main():
    """Ask a question to AISHE server with Redis caching."""
    # Check if question was provided
    if len(sys.argv) < 2:
        print("Usage: python main.py <your question>")
        print("Example: python main.py 'What is the capital of France?'")
        sys.exit(1)

    # Get question from command line arguments
    question = " ".join(sys.argv[1:])

    # Connect to Redis (running in Docker on port 6379)
    try:
        redis_client = redis.Redis(
            host='localhost',
            port=6379,
            db=0,
            decode_responses=False  # We'll handle JSON encoding/decoding ourselves
        )
        # Test connection
        redis_client.ping()
    except redis.ConnectionError:
        print("Error: Could not connect to Redis at localhost:6379")
        print("Make sure Redis is running in Docker.")
        sys.exit(1)

    print(f"Asking: {question}")

    # Check cache first
    cached_response = get_from_cache(redis_client, question)

    if cached_response:
        print("✓ Found in cache! (no API call needed)\n")
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

            # Save to cache for future use
            save_to_cache(redis_client, question, data)
            print("✓ Response saved to cache\n")
            from_cache = False

        except requests.exceptions.ConnectionError:
            print("Error: Could not connect to AISHE server at http://localhost:8000")
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
        print("Source: Redis Cache")
    else:
        print(f"Processing time: {data['processing_time']:.2f} seconds")
    print("=" * 70)


if __name__ == "__main__":
    main()
