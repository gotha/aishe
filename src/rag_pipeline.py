"""RAG Pipeline for Wikipedia-based question answering."""

import asyncio
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass

from mcp_client import WikipediaMCPClient
from ollama_client import OllamaClient


@dataclass
class RAGResult:
    """Result from RAG pipeline."""
    answer: str
    sources: List[Dict[str, str]]
    context_used: str
    query: str


class RAGPipeline:
    """RAG Pipeline combining Wikipedia retrieval and Ollama generation."""

    def __init__(
        self,
        ollama_model: str = "qwen2.5-coder:7b",
        max_context_length: int = 4000,
        max_search_results: int = 3
    ):
        """Initialize the RAG pipeline.
        
        Args:
            ollama_model: Model to use for generation
            max_context_length: Maximum context length in characters
            max_search_results: Maximum number of Wikipedia articles to retrieve
        """
        self.ollama_client = OllamaClient(model=ollama_model)
        self.max_context_length = max_context_length
        self.max_search_results = max_search_results
    
    def process_query(self, query: str) -> str:
        """Process and clean user query.
        
        Args:
            query: Raw user query
            
        Returns:
            Processed query suitable for Wikipedia search
        """
        # Remove common question words that don't help search
        query = query.strip()
        
        # Extract key terms for better Wikipedia search
        # Remove question marks and common filler words
        search_query = query.replace("?", "").strip()
        
        return search_query
    
    async def retrieve_articles(self, query: str) -> List[Dict]:
        """Retrieve relevant Wikipedia articles.
        
        Args:
            query: Search query
            
        Returns:
            List of article data dictionaries
        """
        processed_query = self.process_query(query)
        
        async with WikipediaMCPClient() as mcp_client:
            # Search for articles
            search_results = await mcp_client.search_wikipedia(
                processed_query,
                limit=self.max_search_results
            )
            
            articles = []
            
            # Get the results list from search
            results = search_results.get('results', [])
            if isinstance(results, list):
                for result in results[:self.max_search_results]:
                    # Get article title
                    title = result.get('title', '')
                    if not title:
                        continue
                    
                    try:
                        # Get article summary
                        summary = await mcp_client.get_summary(title)
                        
                        articles.append({
                            'title': title,
                            'summary': summary,
                            'url': f"https://en.wikipedia.org/wiki/{title.replace(' ', '_')}"
                        })
                    except Exception as e:
                        print(f"Error retrieving article '{title}': {e}")
                        continue
            
            return articles
    
    def prepare_context(self, articles: List[Dict]) -> Tuple[str, List[Dict[str, str]]]:
        """Prepare context from retrieved articles.
        
        Args:
            articles: List of article dictionaries
            
        Returns:
            Tuple of (formatted_context, sources)
        """
        context_parts = []
        sources = []
        current_length = 0
        
        for i, article in enumerate(articles, 1):
            title = article.get('title', 'Unknown')
            summary = article.get('summary', '')
            url = article.get('url', '')
            
            # Format article section
            article_text = f"[Source {i}] {title}\n{summary}\n"
            
            # Check if adding this would exceed max length
            if current_length + len(article_text) > self.max_context_length:
                # Truncate to fit
                remaining = self.max_context_length - current_length
                if remaining > 100:  # Only add if we have reasonable space
                    article_text = article_text[:remaining] + "...\n"
                else:
                    break
            
            context_parts.append(article_text)
            sources.append({
                'number': i,
                'title': title,
                'url': url
            })
            current_length += len(article_text)
        
        context = "\n".join(context_parts)
        return context, sources
    
    def generate_answer(self, query: str, context: str) -> str:
        """Generate answer using Ollama with context.
        
        Args:
            query: User's question
            context: Retrieved context from Wikipedia
            
        Returns:
            Generated answer
        """
        # Use the RAG-specific method from OllamaClient
        answer = self.ollama_client.chat_with_context(query, context)
        return answer
    
    async def answer_question(self, query: str) -> RAGResult:
        """Answer a question using the RAG pipeline.
        
        Args:
            query: User's question
            
        Returns:
            RAGResult with answer, sources, and metadata
        """
        # Step 1: Retrieve articles
        articles = await self.retrieve_articles(query)
        
        if not articles:
            return RAGResult(
                answer="I couldn't find any relevant information in Wikipedia to answer your question.",
                sources=[],
                context_used="",
                query=query
            )
        
        # Step 2: Prepare context
        context, sources = self.prepare_context(articles)
        
        # Step 3: Generate answer
        answer = self.generate_answer(query, context)
        
        return RAGResult(
            answer=answer,
            sources=sources,
            context_used=context,
            query=query
        )

