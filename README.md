# GO-HN
GO-HN is a personalized news engine for Hacker News.

It lets you add interests to receive relevant, personalized news while maintaining a rolling database of 1,000 stories.

It periodically fetches fresh news from the Hacker News API, automatically replacing older items to keep the database up to date.

# Features
Hybrid search combining BM25 keyword ranking with semantic embeddings for improved relevance.

Two categories of news:
    Top News – Popular stories ranked by Hacker News.
    New News – Latest stories as they appear.
    
Personalized recommendations based on user-defined interests.

Rolling database that keeps only the most recent 1,000 stories.
