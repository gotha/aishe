"""MCP Client for connecting to Wikipedia MCP server."""

import asyncio
import os
import sys
from pathlib import Path
from typing import Any, Dict, List

from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client


class WikipediaMCPClient:
    """Client for interacting with Wikipedia MCP server."""

    def __init__(self, command: str = None, args: List[str] = None):
        """Initialize the Wikipedia MCP client.

        Args:
            command: Command to run the Wikipedia MCP server (defaults to venv path)
            args: Additional arguments for the server command
        """
        # Default to the venv wikipedia-mcp if no command specified
        if command is None:
            venv_path = Path(sys.executable).parent / "wikipedia-mcp"
            if venv_path.exists():
                command = str(venv_path)
            else:
                command = "wikipedia-mcp"

        self.server_params = StdioServerParameters(
            command=command,
            args=args or [],
            env={**os.environ}  # Pass current environment to subprocess
        )
        self.session: ClientSession | None = None
    
    async def __aenter__(self):
        """Async context manager entry."""
        self._client_context = stdio_client(self.server_params)
        self._read, self._write = await self._client_context.__aenter__()
        
        self._session_context = ClientSession(self._read, self._write)
        self.session = await self._session_context.__aenter__()
        
        # Initialize the connection
        await self.session.initialize()
        
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        """Async context manager exit."""
        if self._session_context:
            await self._session_context.__aexit__(exc_type, exc_val, exc_tb)
        if self._client_context:
            await self._client_context.__aexit__(exc_type, exc_val, exc_tb)
    
    async def search_wikipedia(self, query: str, limit: int = 10) -> Dict[str, Any]:
        """Search Wikipedia for articles matching a query.
        
        Args:
            query: The search term
            limit: Maximum number of results to return
            
        Returns:
            Dictionary containing search results
        """
        if not self.session:
            raise RuntimeError("Client not initialized. Use async with context manager.")
        
        result = await self.session.call_tool(
            "search_wikipedia",
            arguments={"query": query, "limit": limit}
        )
        
        return self._parse_tool_result(result)
    
    async def get_article(self, title: str) -> Dict[str, Any]:
        """Get the full content of a Wikipedia article.
        
        Args:
            title: The title of the Wikipedia article
            
        Returns:
            Dictionary containing article content
        """
        if not self.session:
            raise RuntimeError("Client not initialized. Use async with context manager.")
        
        result = await self.session.call_tool(
            "get_article",
            arguments={"title": title}
        )
        
        return self._parse_tool_result(result)
    
    async def get_summary(self, title: str) -> str:
        """Get a concise summary of a Wikipedia article.
        
        Args:
            title: The title of the Wikipedia article
            
        Returns:
            Summary text
        """
        if not self.session:
            raise RuntimeError("Client not initialized. Use async with context manager.")
        
        result = await self.session.call_tool(
            "get_summary",
            arguments={"title": title}
        )
        
        parsed = self._parse_tool_result(result)
        return parsed.get("text", "")
    
    async def list_tools(self) -> List[str]:
        """List available tools from the Wikipedia MCP server.
        
        Returns:
            List of tool names
        """
        if not self.session:
            raise RuntimeError("Client not initialized. Use async with context manager.")
        
        tools = await self.session.list_tools()
        return [tool.name for tool in tools.tools]
    
    def _parse_tool_result(self, result) -> Dict[str, Any]:
        """Parse tool result into a dictionary.
        
        Args:
            result: CallToolResult from MCP
            
        Returns:
            Parsed result as dictionary
        """
        # Check for structured content first (new format)
        if hasattr(result, 'structuredContent') and result.structuredContent:
            return result.structuredContent
        
        # Fall back to parsing text content
        parsed = {}
        for content in result.content:
            if hasattr(content, 'text'):
                # Try to parse as JSON if possible
                import json
                try:
                    parsed = json.loads(content.text)
                except (json.JSONDecodeError, AttributeError):
                    parsed = {"text": content.text}
        
        return parsed

