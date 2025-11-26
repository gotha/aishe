"""Command-line interface for the RAG system."""

import asyncio
import sys
from pathlib import Path

from rag_pipeline import RAGPipeline


class RAGCLI:
    """Command-line interface for RAG question answering."""

    def __init__(self):
        """Initialize the CLI."""
        self.pipeline = RAGPipeline(
            ollama_model="qwen2.5-coder:7b",
            max_context_length=4000,
            max_search_results=3
        )
    
    def print_banner(self):
        """Print welcome banner."""
        print("=" * 70)
        print("Wikipedia RAG Question Answering System")
        print("=" * 70)
        print("Ask questions and get answers based on Wikipedia articles.")
        print("Type 'quit' or 'exit' to stop.")
        print("=" * 70)
        print()
    
    def print_result(self, result):
        """Print RAG result in a formatted way.
        
        Args:
            result: RAGResult object
        """
        print("\n" + "─" * 70)
        print("ANSWER:")
        print("─" * 70)
        print(result.answer)
        
        if result.sources:
            print("\n" + "─" * 70)
            print("SOURCES:")
            print("─" * 70)
            for source in result.sources:
                print(f"[{source['number']}] {source['title']}")
                print(f"    {source['url']}")
        
        print("─" * 70)
    
    async def run(self):
        """Run the interactive CLI."""
        self.print_banner()
        
        while True:
            try:
                # Get user input
                question = input("\nYour question: ").strip()
                
                # Check for exit commands
                if question.lower() in ['quit', 'exit', 'q']:
                    print("\nGoodbye!")
                    break
                
                # Skip empty questions
                if not question:
                    continue
                
                # Process question
                print("\nSearching Wikipedia and generating answer...")
                result = await self.pipeline.answer_question(question)
                
                # Display result
                self.print_result(result)
                
            except KeyboardInterrupt:
                print("\n\nGoodbye!")
                break
            except Exception as e:
                print(f"\nError: {e}")
                import traceback
                traceback.print_exc()


def main():
    """Entry point for the CLI."""
    cli = RAGCLI()
    asyncio.run(cli.run())


if __name__ == "__main__":
    main()

