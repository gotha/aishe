package com.aishe.client.models;

import java.util.List;

/**
 * Response model for an answer.
 */
public class AnswerResponse {
    /** Generated answer to the given question */
    private String answer;
    
    /** List of sources used to generate the answer */
    private List<Source> sources;
    
    /** Time taken to process the question in seconds */
    private double processing_time;

    /**
     * Default constructor for JSON deserialization.
     */
    public AnswerResponse() {
    }

    /**
     * Constructor with all fields.
     */
    public AnswerResponse(String answer, List<Source> sources, double processingTime) {
        this.answer = answer;
        this.sources = sources;
        this.processing_time = processingTime;
    }

    public String getAnswer() {
        return answer;
    }

    public void setAnswer(String answer) {
        this.answer = answer;
    }

    public List<Source> getSources() {
        return sources;
    }

    public void setSources(List<Source> sources) {
        this.sources = sources;
    }

    public double getProcessingTime() {
        return processing_time;
    }

    public void setProcessingTime(double processingTime) {
        this.processing_time = processingTime;
    }

    @Override
    public String toString() {
        return "AnswerResponse{" +
                "answer='" + answer + '\'' +
                ", sources=" + sources +
                ", processing_time=" + processing_time +
                '}';
    }
}

