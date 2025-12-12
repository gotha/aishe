interface BaseResponse {
    statusCode?: number;
}

/**
 * Request model for asking a question
 */
export interface QuestionRequest {
    /** The question to answer */
    question: string;
}

/**
 * Model for source reference
 */
export interface Source {
    /** Source reference number */
    number: number;

    /** Source article title */
    title: string;

    /** URL of the source article */
    url: string;
}

/**
 * Response model for an answer
 */
export interface AnswerResponse extends BaseResponse {
    /** Generated answer to the given question */
    answer: string;

    /** List of sources used to generate the answer */
    sources: Source[];

    /** Time taken to process the question in seconds */
    processing_time: number;
}

/**
 * Response model for health check
 */
export interface HealthResponse extends BaseResponse {
    /** Server status */
    status: string;

    /** Whether Ollama service is accessible */
    ollama_accessible: boolean;

    /** Additional status message */
    message?: string;
}

/**
 * Response model for errors
 */
export interface ErrorResponse {
    /** Error message */
    error: string;

    /** Detailed error information */
    detail?: string;
}
