#!/usr/bin/env python3
"""Simple command line script to ask questions to AISHE server."""

import sys
import requests


def main():
    """Ask a question to AISHE server and print the result."""
    # Check if question was provided
    if len(sys.argv) < 2:
        print("Usage: python main.py <your question>")
        print("Example: python main.py 'What is the capital of France?'")
        sys.exit(1)

    # Get question from command line arguments
    question = " ".join(sys.argv[1:])

    # AISHE server URL (running in Docker on port 8000)
    url = "http://localhost:8000/api/v1/ask"

    # Prepare request payload
    payload = {"question": question}

    try:
        print(f"Asking: {question}")
        print("Waiting for response...\n")

        # Send POST request to AISHE server
        response = requests.post(url, json=payload, timeout=120)
        response.raise_for_status()

        # Parse response
        data = response.json()

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
        print(f"Processing time: {data['processing_time']:.2f} seconds")
        print("=" * 70)

    except requests.exceptions.ConnectionError:
        print("Error: Could not connect to AISHE server at http://localhost:8000")
        print("Make sure the server is running in Docker.")
        sys.exit(1)
    except requests.exceptions.Timeout:
        print("Error: Request timed out. The question may be too complex.")
        sys.exit(1)
    except requests.exceptions.HTTPError as e:
        print(f"Error: Server returned an error: {e}")
        if response.text:
            print(f"Details: {response.text}")
        sys.exit(1)
    except Exception as e:
        print(f"Unexpected error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
