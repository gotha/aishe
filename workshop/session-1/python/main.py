#!/usr/bin/env python3
"""Simple command line script to ask questions to AISHE server."""

import sys


def main():
    """Ask a question to AISHE server and print the result."""
    # Check if question was provided
    if len(sys.argv) < 2:
        print("Usage: python main.py <your question>")
        print("Example: python main.py 'What is the capital of France?'")
        sys.exit(1)

    # Get question from command line arguments
    question = " ".join(sys.argv[1:])

    print(f"Asking: {question}")
    print("Waiting for response...\n")

    raise Exception("not implemented ...")


if __name__ == "__main__":
    main()
