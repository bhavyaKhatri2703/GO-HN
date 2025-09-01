

def create_text_analyzer(conn):

    cursor = conn.cursor()
    cursor.execute("""
        SELECT create_text_analyzer('text_analyzer1', $$
        pre_tokenizer = "unicode_segmentation"
        [[character_filters]]
        to_lowercase = {}
        [[character_filters]]
        unicode_normalization = "nfkd"
        [[token_filters]]
        skip_non_alphanumeric = {}
        [[token_filters]]
        stopwords = "nltk_english"
        [[token_filters]]
        stemmer = "english_porter2"
        $$);
        """
    )

    cursor.execute("""
        SELECT create_custom_model_tokenizer_and_trigger(
            tokenizer_name => 'tokenizer1',
            model_name => 'model1',
            text_analyzer_name => 'text_analyzer1',
            table_name => 'newStories',
            source_column => 'full_text',
            target_column => 'bm25-embedding'
        );
        """
    )

    cursor.execute("""
        SELECT create_custom_model_tokenizer_and_trigger(
            tokenizer_name => 'tokenizer1',
            model_name => 'model1',
            text_analyzer_name => 'text_analyzer1',
            table_name => 'topStories',
            source_column => 'full_text',
            target_column => 'bm25-embedding'
        );
        """
    )
    cursor.execute("""
        CREATE INDEX new_embedding_bm25 ON newStories USING bm25 (bm25-embedding bm25_ops);
        CREATE INDEX top_embedding_bm25 ON topStories USING bm25 (bm25-embedding bm25_ops);
        """
    )
