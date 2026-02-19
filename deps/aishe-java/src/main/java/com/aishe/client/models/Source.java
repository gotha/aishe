package com.aishe.client.models;

/**
 * Model for source reference.
 */
public class Source {
    /** Source reference number */
    private int number;
    
    /** Source article title */
    private String title;
    
    /** URL of the source article */
    private String url;

    /**
     * Default constructor for JSON deserialization.
     */
    public Source() {
    }

    /**
     * Constructor with all fields.
     */
    public Source(int number, String title, String url) {
        this.number = number;
        this.title = title;
        this.url = url;
    }

    public int getNumber() {
        return number;
    }

    public void setNumber(int number) {
        this.number = number;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getUrl() {
        return url;
    }

    public void setUrl(String url) {
        this.url = url;
    }

    @Override
    public String toString() {
        return "Source{" +
                "number=" + number +
                ", title='" + title + '\'' +
                ", url='" + url + '\'' +
                '}';
    }
}

