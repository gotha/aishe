package com.aishe.client.models;

/**
 * Request model for asking a question.
 */
public class QuestionRequest {
    /** The question to answer */
    private String question;

    /**
     * Default constructor for JSON serialization.
     */
    public QuestionRequest() {
    }

    /**
     * Constructor with question.
     */
    public QuestionRequest(String question) {
        this.question = question;
    }

    public String getQuestion() {
        return question;
    }

    public void setQuestion(String question) {
        this.question = question;
    }

    @Override
    public String toString() {
        return "QuestionRequest{" +
                "question='" + question + '\'' +
                '}';
    }
}

