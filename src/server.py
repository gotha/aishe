"""FastAPI server for the RAG question answering system."""

import asyncio
import time
from contextlib import asynccontextmanager
from typing import Optional

from fastapi import FastAPI, HTTPException, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from api_models import (
    QuestionRequest,
    AnswerResponse,
    HealthResponse,
    ErrorResponse,
    Source
)
from rag_pipeline import RAGPipeline
from ollama_client import OllamaClient
from config import config


# Global pipeline instance
pipeline: Optional[RAGPipeline] = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifespan context manager for startup and shutdown events."""
    global pipeline

    # Startup
    print("Initializing RAG pipeline...")
    print(f"Using Ollama host: {config.OLLAMA_HOST}")
    print(f"Using Ollama model: {config.OLLAMA_MODEL}")
    pipeline = RAGPipeline(
        ollama_model=config.OLLAMA_MODEL,
        ollama_host=config.OLLAMA_HOST,
        max_context_length=config.MAX_CONTEXT_LENGTH,
        max_search_results=config.MAX_SEARCH_RESULTS
    )
    print("RAG pipeline initialized successfully")

    yield

    # Shutdown
    print("Shutting down RAG pipeline...")
    pipeline = None


# Create FastAPI app
app = FastAPI(
    title="AISHE - AI Search & Help Engine",
    description="RAG-based question answering system using Wikipedia and Ollama",
    version="1.0.0",
    lifespan=lifespan
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/", tags=["Root"])
async def root():
    """Root endpoint."""
    return {
        "message": "AISHE - AI Search & Help Engine",
        "version": "1.0.0",
        "docs": "/docs",
        "health": "/health"
    }


@app.get("/health", response_model=HealthResponse, tags=["Health"])
async def health_check():
    """Health check endpoint.

    Returns:
        Health status including Ollama accessibility
    """
    ollama_accessible = False
    message = None

    try:
        # Try to check if Ollama is accessible
        test_client = OllamaClient(host=config.OLLAMA_HOST, model=config.OLLAMA_MODEL)
        # Try to list models to verify connection
        test_client.list_models()
        ollama_accessible = True
    except Exception as e:
        message = f"Ollama not accessible: {str(e)}"

    return HealthResponse(
        status="healthy" if ollama_accessible else "degraded",
        ollama_accessible=ollama_accessible,
        message=message
    )


@app.post(
    "/api/v1/ask",
    response_model=AnswerResponse,
    responses={
        500: {"model": ErrorResponse},
        400: {"model": ErrorResponse}
    },
    tags=["Question Answering"]
)
async def ask_question(request: QuestionRequest):
    """Answer a question using the RAG pipeline.

    Args:
        request: Question request containing the question to answer

    Returns:
        Answer response with the generated answer and sources

    Raises:
        HTTPException: If the pipeline is not initialized or an error occurs
    """
    if pipeline is None:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="RAG pipeline not initialized"
        )

    try:
        # Record start time
        start_time = time.time()

        # Process the question
        result = await pipeline.answer_question(request.question)

        # Calculate processing time
        processing_time = time.time() - start_time

        # Convert sources to API model
        sources = [
            Source(
                number=source["number"],
                title=source["title"],
                url=source["url"]
            )
            for source in result.sources
        ]

        return AnswerResponse(
            answer=result.answer,
            sources=sources,
            processing_time=processing_time
        )

    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Error processing question: {str(e)}"
        )


if __name__ == "__main__":
    import uvicorn
    import argparse

    parser = argparse.ArgumentParser(description="AISHE API Server")
    parser.add_argument(
        "--host",
        default=config.SERVER_HOST,
        help=f"Host to bind to (default: {config.SERVER_HOST})"
    )
    parser.add_argument(
        "--port",
        type=int,
        default=config.SERVER_PORT,
        help=f"Port to bind to (default: {config.SERVER_PORT})"
    )
    args = parser.parse_args()

    print(f"Starting AISHE server on {args.host}:{args.port}")
    uvicorn.run(app, host=args.host, port=args.port)
